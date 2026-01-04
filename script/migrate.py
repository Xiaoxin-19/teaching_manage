#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
数据库迁移脚本
从 MySQL 数据库读取数据，迁移到 SQLite 数据库
注意：MySQL 和 SQLite 的表结构不完全对应，需要根据业务逻辑处理：
1. MySQL records 表没有 subject_id，需要通过学生-教师关系查找
2. MySQL 没有 student_subjects 表，需要根据 records 和 orders 计算生成
3. 订单数据需要关联到具体的学生课程
"""

import mysql.connector
import sqlite3
import os
import sys

from datetime import datetime, time
from typing import List, Dict, Tuple, Optional
from collections import defaultdict
import time as time_module


class DatabaseMigrator:
    """数据库迁移器：从 MySQL 迁移到 SQLite"""
    
    def __init__(self, mysql_config: Dict, sqlite_path: str, clear_data: bool = True):
        """初始化迁移器
        
        Args:
            mysql_config: MySQL 连接配置
            sqlite_path: SQLite 数据库文件路径
            clear_data: 是否清空现有数据（默认 True）
        """
        self.mysql_config = mysql_config
        self.sqlite_path = sqlite_path
        self.clear_data = clear_data
        self.mysql_conn = None
        self.sqlite_conn = None
        # 缓存映射关系
        self.student_teacher_subject_map = {}  # (student_id, teacher_id) -> subject_id
        self.student_subject_ids = {}  # (student_id, teacher_id, subject_id) -> student_subject_id
        self.student_deleted_at = {}  # student_id -> deleted_at
        self.student_current_teacher = {}  # student_id -> current_teacher_id（用于查询历史教师的记录）
    
    def connect(self) -> bool:
        """连接到 MySQL 和 SQLite 数据库
        
        Returns:
            True 如果连接成功，False 否则
        """
        try:
            # 连接 MySQL
            print("正在连接 MySQL 数据库...")
            self.mysql_conn = mysql.connector.connect(**self.mysql_config)
            print(f"[OK] MySQL 连接成功: {self.mysql_config['host']}")
            
            # 创建或连接 SQLite
            print(f"正在创建/连接 SQLite 数据库: {self.sqlite_path}")
            self.sqlite_conn = sqlite3.connect(self.sqlite_path)
            self.sqlite_conn.row_factory = sqlite3.Row
            print("[OK] SQLite 连接成功")
            
            return True
        except Exception as e:
            print(f"[ERROR] 数据库连接失败: {e}")
            return False
    
    def close(self):
        """关闭数据库连接"""
        if self.mysql_conn:
            self.mysql_conn.close()
        if self.sqlite_conn:
            self.sqlite_conn.close()
    
    def _map_gender(self, gender: str) -> str:
        """映射性别值：男→male，女→female
        
        Args:
            gender: MySQL中的性别值
            
        Returns:
            SQLite中的性别值
        """
        if not gender:
            return ''
        
        gender_str = str(gender).strip()
        gender_map = {
            '男': 'male',
            '女': 'female',
            'male': 'male',
            'female': 'female'
        }
        
        return gender_map.get(gender_str, gender_str)
    
    def _format_datetime_for_gorm(self, dt) -> Optional[str]:
        """将 datetime 转换为 GORM 兼容的格式
        
        Args:
            dt: datetime 对象或字符串
            
        Returns:
            格式化后的时间字符串，格式: YYYY-MM-DD HH:MM:SS.ffffff+08:00
        """
        if dt is None:
            return None
        
        if isinstance(dt, str):
            # 如果已经是字符串，尝试解析
            try:
                dt = datetime.strptime(dt, '%Y-%m-%d %H:%M:%S')
            except ValueError:
                return dt
        
        if isinstance(dt, datetime):
            # 转换为 GORM 格式: YYYY-MM-DD HH:MM:SS.ffffff+08:00
            # 假设使用东八区时区
            return dt.strftime('%Y-%m-%d %H:%M:%S.%f') + '+08:00'
        
        return None
    
    def _generate_order_number(self, create_time, sequence: int = 0) -> str:
        """根据创建时间生成类雪花算法的订单号
        
        Args:
            create_time: 订单创建时间
            sequence: 序列号（同一毫秒内的订单序号）
            
        Returns:
            订单号字符串
        """
        if create_time is None:
            create_time = datetime.now()
        
        # 确保 create_time 是 datetime 对象
        if isinstance(create_time, str):
            try:
                create_time = datetime.strptime(create_time, '%Y-%m-%d %H:%M:%S')
            except ValueError:
                create_time = datetime.now()
        
        # 使用时间戳（毫秒）作为基础
        # 基准时间：2020-01-01 00:00:00
        epoch = datetime(2020, 1, 1)
        timestamp_ms = int((create_time - epoch).total_seconds() * 1000)
        
        # 生成雪花ID样式的订单号
        # 格式：时间戳(41位) + 序列号(12位) = 53位数字
        # 为了方便阅读，转换为字符串
        order_id = (timestamp_ms << 12) | (sequence & 0xFFF)
        
        return str(order_id)
    
    def clear_existing_data(self):
        """清空 SQLite 数据库中的现有数据"""
        print("\n清空现有数据...")
        cursor = self.sqlite_conn.cursor()
        
        try:
            # 按照依赖顺序删除数据（先删除依赖表，后删除主表）
            tables = ['records', 'recharge_orders', 'student_subjects', 'students', 'subjects', 'teachers']
            
            for table in tables:
                cursor.execute(f"DELETE FROM {table}")
                print(f"  清空 {table} 表")
            
            self.sqlite_conn.commit()
            print("[OK] 现有数据清空完成")
            
        except Exception as e:
            print(f"[ERROR] 清空数据失败: {e}")
            raise
    
    def verify_schema(self):
        """验证 SQLite 数据库表结构是否存在"""
        print("\n验证数据库表结构...")
        cursor = self.sqlite_conn.cursor()
        
        required_tables = ['teachers', 'subjects', 'students', 'student_subjects', 'records', 'recharge_orders']
        
        try:
            cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
            existing_tables = [row[0] for row in cursor.fetchall()]
            
            missing_tables = []
            for table in required_tables:
                if table in existing_tables:
                    print(f"  [OK] {table} 表存在")
                else:
                    missing_tables.append(table)
                    print(f"  [ERROR] {table} 表不存在")
            
            if missing_tables:
                raise Exception(f"缺少必需的表: {', '.join(missing_tables)}。请确保 SQLite 数据库已正确初始化。")
            
            print("[OK] 数据库表结构验证通过")
            
        except Exception as e:
            print(f"[ERROR] 验证表结构失败: {e}")
            raise
    
    def build_student_teacher_subject_relationships(self):
        """构建学生-教师-科目的关系映射
        从 MySQL students 表的 teacher_id 字段获取当前学生-教师关系
        
        设计意图：
        - 只维护学生的当前教师关系（简洁）
        - 当 records 中有历史教师时，用当前教师的科目ID代替
        - 这样所有记录都能找到科目映射，不会有"找不到科目映射"的错误
        """
        print("\n构建学生-教师-科目关系（从 students 表的 teacher_id 字段）...")
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        
        try:
            # 从 students 表获取当前的学生-教师关系
            mysql_cursor.execute("""
                SELECT id as student_id, teacher_id
                FROM students
                WHERE teacher_id IS NOT NULL
            """)
            combinations = mysql_cursor.fetchall()
            
            # 保存学生的当前教师映射（用于后续查询）
            self.student_current_teacher = {}
            for combo in combinations:
                self.student_current_teacher[combo['student_id']] = combo['teacher_id']
            
            # 获取 SQLite 中的科目列表
            sqlite_cursor = self.sqlite_conn.cursor()
            sqlite_cursor.execute("SELECT id, name FROM subjects")
            subjects = {row[1]: row[0] for row in sqlite_cursor.fetchall()}
            
            # 默认科目 ID（如果没有匹配的科目，使用第一个科目），如果没有科目则创建默认科目
            if subjects:
                default_subject_id = list(subjects.values())[0]
            else:
                # 创建默认科目
                print("  未找到科目，创建默认科目...")
                sqlite_cursor.execute('''
                    INSERT INTO subjects (name, subject_number, created_at, updated_at)
                    VALUES (?, ?, ?, ?)
                ''', ('默认科目', 'SUBJ0001', self._format_datetime_for_gorm(datetime.now()), self._format_datetime_for_gorm(datetime.now())))
                self.sqlite_conn.commit()
                default_subject_id = sqlite_cursor.lastrowid
                subjects['默认科目'] = default_subject_id
                print(f"  [OK] 已创建默认科目 (ID: {default_subject_id})")
            
            # 构建映射：每个(student_id, current_teacher_id)对应一个科目
            for combo in combinations:
                student_id = combo['student_id']
                teacher_id = combo['teacher_id']
                
                # 使用默认科目
                subject_id = default_subject_id
                
                self.student_teacher_subject_map[(student_id, teacher_id)] = subject_id
            
            print(f"[OK] 从 students 表构建了 {len(self.student_teacher_subject_map)} 个学生-教师-科目关系")
            
        except Exception as e:
            print(f"[ERROR] 构建关系失败: {e}")
            raise
        finally:
            mysql_cursor.close()
    
    def migrate_teachers(self) -> int:
        """迁移 teachers 表"""
        print("\n正在迁移 teachers 表...")
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            mysql_cursor.execute("SELECT * FROM teachers")
            teachers = mysql_cursor.fetchall()
            
            count = 0
            for teacher in teachers:
                sqlite_cursor.execute('''
                    INSERT INTO teachers 
                    (id, teacher_number, name, gender, phone, remark, created_at, updated_at, deleted_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
                ''', (
                    teacher['id'],
                    f"T{teacher['id']:08d}",  # 格式: T + 8位数字 (T00000001)
                    teacher['name'],
                    self._map_gender(teacher.get('gender', '')),
                    teacher.get('phone_code'),
                    '',
                    self._format_datetime_for_gorm(teacher.get('create_time')),
                    self._format_datetime_for_gorm(teacher.get('update_time')),
                    self._format_datetime_for_gorm(teacher.get('delete_time'))
                ))
                count += 1
            
            self.sqlite_conn.commit()
            print(f"[OK] 迁移 {count} 条教师记录")
            return count
            
        except Exception as e:
            print(f"[ERROR] 迁移 teachers 表失败: {e}")
            raise
        finally:
            mysql_cursor.close()
    
    def migrate_subjects(self) -> int:
        """检查并创建默认科目（MySQL 中没有 subjects 表）"""
        print("\n检查科目表...")
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            # 检查是否已有科目
            sqlite_cursor.execute("SELECT COUNT(*) FROM subjects")
            count = sqlite_cursor.fetchone()[0]
            
            if count == 0:
                # 创建默认科目
                print("  未找到科目，创建默认科目...")
                sqlite_cursor.execute('''
                    INSERT INTO subjects (name, subject_number, created_at, updated_at)
                    VALUES (?, ?, ?, ?)
                ''', ('默认科目', 'SUBJ0001', self._format_datetime_for_gorm(datetime.now()), self._format_datetime_for_gorm(datetime.now())))
                self.sqlite_conn.commit()
                print(f"  [OK] 已创建默认科目")
                return 1
            else:
                print(f"[OK] 已有 {count} 个科目，跳过迁移")
                return 0
            
        except Exception as e:
            print(f"[ERROR] 检查科目表失败: {e}")
            raise
    
    def migrate_students(self) -> int:
        """迁移 students 表"""
        print("\n正在迁移 students 表...")
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            mysql_cursor.execute("SELECT * FROM students")
            students = mysql_cursor.fetchall()
            
            count = 0
            current_year = datetime.now().year
            for student in students:
                student_id = student['id']
                deleted_at = student.get('delete_time')
                
                # 缓存学生的删除状态
                self.student_deleted_at[student_id] = deleted_at
                
                sqlite_cursor.execute('''
                    INSERT INTO students 
                    (id, name, student_number, gender, phone, status, remark, created_at, updated_at, deleted_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                ''', (
                    student_id,
                    student['name'],
                    f"S{current_year}{student_id:05d}",  # 格式: S + 年份 + 5位数字 (S202400005)
                    self._map_gender(student.get('gender', '')),
                    student.get('phone_code'),
                    1,  # status, 默认为 1(正常)
                    '',
                    self._format_datetime_for_gorm(student.get('create_time')),
                    self._format_datetime_for_gorm(student.get('update_time')),
                    self._format_datetime_for_gorm(deleted_at)
                ))
                count += 1
            
            self.sqlite_conn.commit()
            print(f"[OK] 迁移 {count} 条学生记录")
            return count
            
        except Exception as e:
            print(f"[ERROR] 迁移 students 表失败: {e}")
            raise
        finally:
            mysql_cursor.close()
    
    def migrate_student_subjects(self) -> int:
        """迁移/生成 student_subjects 表
        注意：MySQL 中没有 student_subjects 表，需要根据 records 生成，课时余额从 students 表的 hours 字段读取
        """
        print("\n生成 student_subjects 表数据...")
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            # 获取每个学生-教师-科目组合的课时统计
            student_course_data = {}  # (student_id, teacher_id, subject_id) -> data
            
            # 1. 获取所有学生的剩余课时（直接从 students 表读取 hours 字段）
            mysql_cursor.execute("""
                SELECT id, hours
                FROM students
            """)
            student_hours = {row['id']: row.get('hours', 0) for row in mysql_cursor.fetchall()}
            
            # 2. 生成 student_subjects 数据
            for (student_id, teacher_id), subject_id in self.student_teacher_subject_map.items():
                key = (student_id, teacher_id, subject_id)
                
                # 直接使用学生表中的剩余课时（可以为负数）
                balance = student_hours.get(student_id, 0)
                
                student_course_data[key] = {
                    'student_id': student_id,
                    'teacher_id': teacher_id,
                    'subject_id': subject_id,
                    'balance': balance,  # 允许负数
                    'avg_price': 0.0,  # 需要根据实际业务计算
                    'status': 1,  # 1=正常
                    'remark': '从 MySQL 迁移生成'
                }
            
            # 3. 插入到 SQLite
            count = 0
            for data in student_course_data.values():
                student_id = data['student_id']
                # 如果学生已删除，则该课程记录也标记为已删除
                deleted_at = self.student_deleted_at.get(student_id)
                
                sqlite_cursor.execute('''
                    INSERT INTO student_subjects 
                    (student_id, teacher_id, subject_id, balance, total_buy, avg_price, status, remark, created_at, updated_at, deleted_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                ''', (
                    student_id,
                    data['teacher_id'],
                    data['subject_id'],
                    data['balance'],
                    data['balance'],  # total_buy 初始值等于 balance
                    data['avg_price'],
                    data['status'],
                    data['remark'],
                    self._format_datetime_for_gorm(datetime.now()),
                    self._format_datetime_for_gorm(datetime.now()),
                    self._format_datetime_for_gorm(deleted_at)
                ))
                
                # 保存 student_subject_id 映射
                student_subject_id = sqlite_cursor.lastrowid
                key = (data['student_id'], data['teacher_id'], data['subject_id'])
                self.student_subject_ids[key] = student_subject_id
                
                count += 1
            
            self.sqlite_conn.commit()
            print(f"[OK] 生成 {count} 条学生课程记录")
            return count
            
        except Exception as e:
            print(f"[ERROR] 生成 student_subjects 表失败: {e}")
            raise
        finally:
            mysql_cursor.close()
    
    def migrate_records(self) -> int:
        """迁移 records 表
        注意：MySQL records 表没有 subject_id，需要从关系映射中获取
        MySQL 使用 teaching_time 字符串，SQLite 使用 start_time/end_time
        """
        print("\n正在迁移 records 表...")
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            mysql_cursor.execute('''
                SELECT id, student_id, teacher_id, teaching_date, teaching_time,
                       active, comment, create_time, update_time, delete_time
                FROM records
            ''')
            records = mysql_cursor.fetchall()
            
            count = 0
            skipped = 0
            for record in records:
                student_id = record['student_id']
                teacher_id = record['teacher_id']
                
                # 从映射中获取 subject_id
                subject_id = self.student_teacher_subject_map.get((student_id, teacher_id))
                if not subject_id:
                    print(f"  警告: 找不到学生 {student_id} 和教师 {teacher_id} 的科目映射，跳过")
                    skipped += 1
                    continue
                
                # 解析 teaching_time (格式：'18:15-19:00')
                teaching_date = record.get('teaching_date')
                teaching_time = record.get('teaching_time', '')
                start_time, end_time = self._parse_teaching_time(teaching_date, teaching_time)
                
                # 如果学生已删除，则该记录也标记为已删除
                record_deleted_at = record.get('delete_time')
                if not record_deleted_at:
                    record_deleted_at = self.student_deleted_at.get(student_id)
                
                # 转换时间为字符串格式 (HH:MM)
                start_time_str = start_time.strftime('%H:%M') if start_time else None
                end_time_str = end_time.strftime('%H:%M') if end_time else None
                
                # 生成 teaching_date_ms (Unix 毫秒，基于 teaching_date 的午夜时刻)
                teaching_date_ms = None
                if teaching_date:
                    base_date = teaching_date if isinstance(teaching_date, datetime) else datetime.strptime(str(teaching_date), '%Y-%m-%d')
                    teaching_date_ms = int(base_date.timestamp() * 1000)
                
                try:
                    sqlite_cursor.execute('''
                        INSERT OR IGNORE INTO records 
                        (student_id, teacher_id, subject_id, price_snapshot, teaching_date,
                         teaching_date_ms, start_time, end_time, active, remark, created_at, updated_at, deleted_at)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    ''', (
                        student_id,
                        teacher_id,
                        subject_id,
                        0.0,  # price_snapshot，需要根据实际业务计算
                        teaching_date,
                        teaching_date_ms,
                        start_time_str,
                        end_time_str,
                        record.get('active', 0),
                        record.get('comment', ''),
                        self._format_datetime_for_gorm(record.get('create_time')),
                        self._format_datetime_for_gorm(record.get('update_time')),
                        self._format_datetime_for_gorm(record_deleted_at)
                    ))
                    count += 1
                except Exception as record_error:
                    skipped += 1
                    stats['errors'].append(f"记录 {record['id']} 迁移失败: {str(record_error)}")
                    continue
            
            self.sqlite_conn.commit()
            print(f"[OK] 迁移 {count} 条教学记录" + (f"（跳过 {skipped} 条）" if skipped > 0 else ""))
            return count
            
        except Exception as e:
            print(f"[ERROR] 迁移 records 表失败: {e}")
            print(f"[DEBUG] 上下文信息:")
            print(f"  - 学生ID: {student_id}")
            print(f"  - 教师ID: {teacher_id if 'teacher_id' in locals() else 'N/A'}")
            print(f"  - 当前记录ID: {record['id'] if 'record' in locals() else 'N/A'}")
            raise
        finally:
            mysql_cursor.close()
    
    def _parse_teaching_time(self, teaching_date, teaching_time: str) -> Tuple[Optional[datetime], Optional[datetime]]:
        """解析 teaching_time 字符串为 start_time 和 end_time
        
        支持的格式:
        - 10:00-11:00 (标准格式)
        - 10：00-11：00 (全角冒号)
        - 10:00～11:00 (波浪线分隔)
        - 14::0-15::0 (双冒号)
        - 18:30:19:15 (多冒号 -> 智能提取为 18:30-19:15)
        - 20：00-20"45 (全角冒号和引号混用 -> 20:00-20:45)
        - 17;30-18:15 (分号代替冒号)
        
        Args:
            teaching_date: 上课日期
            teaching_time: 时间字符串
            
        Returns:
            (start_time, end_time) 元组
        """
        if not teaching_time:
            return None, None
        
        try:
            teaching_time = teaching_time.strip()
            
            # ===== Step 1: 标准化所有特殊字符 =====
            # 替换冒号变体
            teaching_time = teaching_time.replace('：', ':')  # 全角冒号
            teaching_time = teaching_time.replace('；', ':')  # 中文分号
            teaching_time = teaching_time.replace(';', ':')   # 英文分号
            
            # 替换连接符
            teaching_time = teaching_time.replace('－', '-')  # 全角连字符
            teaching_time = teaching_time.replace('—', '-')   # 破折号
            teaching_time = teaching_time.replace('~', '-')   # 波浪线
            teaching_time = teaching_time.replace('～', '-')  # 全角波浪线
            teaching_time = teaching_time.replace('_', '-')   # 下划线
            
            # 替换引号为冒号（引号常在时分位置被误用）
            teaching_time = teaching_time.replace('"', ':')   # 中文左双引号
            teaching_time = teaching_time.replace('"', ':')   # 中文右双引号
            teaching_time = teaching_time.replace("'", ':')   # 中文左单引号
            teaching_time = teaching_time.replace("'", ':')   # 中文右单引号
            teaching_time = teaching_time.replace('"', ':')   # 英文双引号
            teaching_time = teaching_time.replace("'", ':')   # 英文单引号
            
            # ===== Step 2: 处理常见格式错误 =====
            teaching_time = teaching_time.replace('::', ':')   # 双冒号
            teaching_time = teaching_time.replace(' ', '')     # 空格
            
            # ===== Step 3: 处理多冒号的情况 (18:30:19:15) =====
            colon_count = teaching_time.count(':')
            if colon_count > 2:
                parts = teaching_time.split('-')
                
                if len(parts) == 1:
                    # 没有分隔符 - 整个字符串可能是错误格式
                    # 例如 18:30:19:15 -> 提取为 18:30-19:15
                    colon_positions = [i for i, c in enumerate(teaching_time) if c == ':']
                    
                    if len(colon_positions) >= 2:
                        # 提取第一个时间 (HH:MM)
                        first_colon = colon_positions[0]
                        second_colon = colon_positions[1]
                        start_str = teaching_time[max(0, first_colon-2):min(len(teaching_time), second_colon+3)]
                        
                        # 从剩余部分提取第二个时间
                        if len(colon_positions) > 2:
                            # 有更多冒号，使用后面的冒号
                            third_colon = colon_positions[2]
                            fourth_colon = colon_positions[3] if len(colon_positions) > 3 else len(teaching_time)-1
                            end_str = teaching_time[max(0, third_colon-2):min(len(teaching_time), fourth_colon+3)]
                        else:
                            end_str = start_str  # 没有更多时间，重复
                        
                        # 清理提取的字符串，只保留数字和冒号
                        start_str = ''.join(c for c in start_str if c.isdigit() or c == ':')
                        end_str = ''.join(c for c in end_str if c.isdigit() or c == ':')
                        
                        # 确保是 HH:MM 格式
                        start_parts = start_str.split(':')
                        end_parts = end_str.split(':')
                        
                        start_str = f"{start_parts[0]}:{start_parts[1]}" if len(start_parts) >= 2 else start_str
                        end_str = f"{end_parts[0]}:{end_parts[1]}" if len(end_parts) >= 2 else end_str
                        
                        teaching_time = f"{start_str}-{end_str}"
                else:
                    # 有分隔符，处理每个部分
                    cleaned_parts = []
                    for part in parts:
                        part_colons = part.count(':')
                        if part_colons > 1:
                            # 多个冒号，只保留前两个
                            colon_pos = [i for i, c in enumerate(part) if c == ':']
                            if len(colon_pos) >= 2:
                                first = colon_pos[0]
                                second = colon_pos[1]
                                time_str = part[max(0, first-2):min(len(part), second+3)]
                                time_str = ''.join(c for c in time_str if c.isdigit() or c == ':')
                                parts_list = time_str.split(':')
                                time_str = f"{parts_list[0]}:{parts_list[1]}" if len(parts_list) >= 2 else time_str
                                cleaned_parts.append(time_str)
                            else:
                                cleaned_parts.append(part)
                        else:
                            cleaned_parts.append(part)
                    
                    if len(cleaned_parts) >= 2:
                        teaching_time = f"{cleaned_parts[0]}-{cleaned_parts[-1]}"
            
            # ===== Step 4: 分割时间段 =====
            if '-' not in teaching_time:
                return None, None
            
            parts = teaching_time.split('-')
            if len(parts) < 2:
                return None, None
            
            start_str = parts[0].strip()
            end_str = parts[-1].strip()
            
            # ===== Step 5: 解析时间 =====
            def parse_time_component(time_str):
                """解析单个时间部分 (HH:MM 或 HHMM)"""
                if ':' in time_str:
                    time_parts = time_str.split(':')
                    time_parts = [p for p in time_parts if p and p.isdigit()]
                    if not time_parts:
                        return None, None
                    hour = int(time_parts[0])
                    minute = int(time_parts[1]) if len(time_parts) > 1 else 0
                else:
                    digits = ''.join(c for c in time_str if c.isdigit())
                    if not digits:
                        return None, None
                    hour = int(digits[0:2]) if len(digits) >= 2 else int(digits[0])
                    minute = int(digits[2:4]) if len(digits) >= 4 else 0
                
                return hour, minute
            
            start_hour, start_min = parse_time_component(start_str)
            end_hour, end_min = parse_time_component(end_str)
            
            if start_hour is None or end_hour is None:
                return None, None
            
            # ===== Step 6: 验证时间范围和格式 =====
            # 检查小时和分钟是否为整数且在有效范围内
            if not isinstance(start_hour, int) or not isinstance(start_min, int) or \
               not isinstance(end_hour, int) or not isinstance(end_min, int):
                return None, None
            
            if not (0 <= start_hour <= 23 and 0 <= start_min <= 59 and 
                    0 <= end_hour <= 23 and 0 <= end_min <= 59):
                return None, None
            
            # ===== Step 7: 构建时间对象 =====
            if teaching_date:
                base_date = teaching_date if isinstance(teaching_date, datetime) else datetime.combine(teaching_date, time.min)
                start_time = base_date.replace(hour=start_hour, minute=start_min, second=0)
                end_time = base_date.replace(hour=end_hour, minute=end_min, second=0)
                
                # ===== Step 8: 再次验证格式为 HH:MM =====
                # 确保时间对象可以正确格式化为 HH:MM
                start_formatted = start_time.strftime('%H:%M')
                end_formatted = end_time.strftime('%H:%M')
                
                # 验证格式是否匹配 HH:MM（两位数字-冒号-两位数字）
                if not (len(start_formatted) == 5 and start_formatted[2] == ':' and 
                        len(end_formatted) == 5 and end_formatted[2] == ':'):
                    return None, None
                
                return start_time, end_time
            
        except (ValueError, AttributeError, IndexError, TypeError) as e:
            pass
        
        return None, None
    
    def migrate_recharge_orders(self) -> int:
        """迁移 recharge_orders 表（从 orders 表迁移）
        注意：需要关联到具体的 student_subject_id
        """
        print("\n正在迁移 recharge_orders 表...")
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            mysql_cursor.execute('''
                SELECT id, student_id, hours, comment, create_time, active, type
                FROM orders
            ''')
            orders = mysql_cursor.fetchall()
            
            count = 0
            skipped = 0
            for order in orders:
                student_id = order['student_id']
                # 根据 Snowflake 算法生成订单号（格式: ORD + Snowflake ID）
                order_id = order['id']
                snowflake_id = int(order.get('create_time').timestamp() * 1000) + order_id
                order_number = f"ORD{snowflake_id}"
                
                # 查找该学生的主要课程（第一个课程）
                student_subject_id = None
                for key, ss_id in self.student_subject_ids.items():
                    if key[0] == student_id:  # key = (student_id, teacher_id, subject_id)
                        student_subject_id = ss_id
                        break
                
                if not student_subject_id:
                    print(f"  警告: 找不到学生 {student_id} 的课程记录，跳过订单 {order_number}")
                    skipped += 1
                    continue
                
                # 如果学生已删除，则该订单也标记为已删除
                deleted_at = self.student_deleted_at.get(student_id)
                
                # 根据 hours 的正负判断订单类型
                hours = order['hours']
                
                sqlite_cursor.execute('''
                    INSERT INTO recharge_orders 
                    (order_number, student_course_id, hours, amount, remark, created_at, updated_at, deleted_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                ''', (
                    order_number,
                    student_subject_id,
                    hours,
                    0.0,  # amount 需要根据实际业务计算
                    order.get('comment'),
                    self._format_datetime_for_gorm(order.get('create_time')),
                    self._format_datetime_for_gorm(order.get('create_time')),
                    self._format_datetime_for_gorm(deleted_at)
                ))
                count += 1
            
            self.sqlite_conn.commit()
            print(f"[OK] 迁移 {count} 条充值订单" + (f"（跳过 {skipped} 条）" if skipped > 0 else ""))
            return count
            
        except Exception as e:
            print(f"[ERROR] 迁移 recharge_orders 表失败: {e}")
            raise
        finally:
            mysql_cursor.close()
    
    def migrate_student_data(self, student_id: int, student_name: str) -> Dict:
        """迁移单个学生的所有关联数据
        
        Args:
            student_id: 学生ID
            student_name: 学生姓名（用于日志）
            
        Returns:
            包含迁移统计信息的字典
        """
        stats = {
            'student_subjects': 0,
            'records': 0,
            'recharge_orders': 0,
            'errors': []
        }
        
        mysql_cursor = self.mysql_conn.cursor(dictionary=True)
        sqlite_cursor = self.sqlite_conn.cursor()
        
        try:
            print(f"\n  → 迁移学生: {student_name} (ID: {student_id})")
            
            # 1. 迁移该学生的 student_subjects
            mysql_cursor.execute("SELECT id, hours FROM students WHERE id = %s", (student_id,))
            student_data = mysql_cursor.fetchone()
            if not student_data:
                stats['errors'].append("找不到学生数据")
                return stats
            
            balance = student_data.get('hours', 0)
            deleted_at = self.student_deleted_at.get(student_id)
            
            # 查找该学生的所有教师关系
            student_subjects_created = []
            for (sid, tid), subject_id in self.student_teacher_subject_map.items():
                if sid != student_id:
                    continue
                
                try:
                    sqlite_cursor.execute('''
                        INSERT INTO student_subjects 
                        (student_id, teacher_id, subject_id, balance, total_buy, avg_price, status, remark, created_at, updated_at, deleted_at)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    ''', (
                        student_id, tid, subject_id, balance, balance, 0.0,
                        1,  # 1=正常
                        '从 MySQL 迁移生成',
                        self._format_datetime_for_gorm(datetime.now()),
                        self._format_datetime_for_gorm(datetime.now()),
                        self._format_datetime_for_gorm(deleted_at)
                    ))
                    
                    student_subject_id = sqlite_cursor.lastrowid
                    self.student_subject_ids[(student_id, tid, subject_id)] = student_subject_id
                    student_subjects_created.append((tid, student_subject_id))
                    stats['student_subjects'] += 1
                except Exception as ss_error:
                    stats['errors'].append(f"创建课程关系 (teacher_id={tid}) 失败: {str(ss_error)}")
                    print(f"    [DEBUG] 课程创建失败 - 学生: {student_id}, 教师: {tid}, 错误: {ss_error}")
                    continue
            
            print(f"    [OK] 创建 {stats['student_subjects']} 个课程关系")
            
            # 如果没有课程关系，只记录警告但不失败
            if stats['student_subjects'] == 0:
                stats['errors'].append(f"学生 {student_name} 没有课程关系")
            
            # 2. 迁移该学生的 records
            mysql_cursor.execute('''
                SELECT id, student_id, teacher_id, teaching_date, teaching_time,
                       active, comment, create_time, update_time, delete_time
                FROM records
                WHERE student_id = %s
            ''', (student_id,))
            records = mysql_cursor.fetchall()
            
            for record in records:
                teacher_id = record['teacher_id']
                subject_id = self.student_teacher_subject_map.get((student_id, teacher_id))
                
                # 如果找不到该(student_id, teacher_id)的映射，说明是历史教师
                # 直接使用当前教师的科目ID
                if not subject_id:
                    current_teacher_id = self.student_current_teacher.get(student_id)
                    if current_teacher_id:
                        subject_id = self.student_teacher_subject_map.get((student_id, current_teacher_id))
                
                if not subject_id:
                    stats['errors'].append(f"记录 {record['id']} 找不到科目映射")
                    continue
                
                teaching_date = record.get('teaching_date')
                teaching_time = record.get('teaching_time', '')
                start_time, end_time = self._parse_teaching_time(teaching_date, teaching_time)
                
                # 如果时间无法解析，跳过该记录
                if start_time is None or end_time is None:
                    stats['errors'].append(f"记录 {record['id']} 的时间无法解析: '{teaching_time}'")
                    continue
                
                record_deleted_at = record.get('delete_time') or deleted_at
                
                # 转换时间为字符串格式 (HH:MM)
                start_time_str = start_time.strftime('%H:%M') if start_time else None
                end_time_str = end_time.strftime('%H:%M') if end_time else None
                
                # 转换 teaching_date 为字符串格式 (YYYY-MM-DD)
                teaching_date_str = teaching_date.strftime('%Y-%m-%d') if isinstance(teaching_date, datetime) else teaching_date
                
                # 生成 teaching_date_ms (Unix 毫秒，基于 teaching_date 的午夜时刻)
                teaching_date_ms = None
                if teaching_date:
                    base_date = teaching_date if isinstance(teaching_date, datetime) else datetime.strptime(str(teaching_date), '%Y-%m-%d')
                    teaching_date_ms = int(base_date.timestamp() * 1000)
                
                try:
                    sqlite_cursor.execute('''
                        INSERT OR IGNORE INTO records 
                        (student_id, teacher_id, subject_id, price_snapshot, teaching_date,
                         teaching_date_ms, start_time, end_time, active, remark, created_at, updated_at, deleted_at)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    ''', (
                        student_id, teacher_id, subject_id, 0.0, teaching_date_str,
                        teaching_date_ms, start_time_str, end_time_str, record.get('active', 0),
                        record.get('comment', ''),
                        self._format_datetime_for_gorm(record.get('create_time')),
                        self._format_datetime_for_gorm(record.get('update_time')),
                        self._format_datetime_for_gorm(record_deleted_at)
                    ))
                    if sqlite_cursor.rowcount > 0:
                        stats['records'] += 1
                except Exception as record_error:
                    stats['errors'].append(f"记录 {record['id']} 插入失败: {str(record_error)}")
                    continue
            
            print(f"    [OK] 迁移 {stats['records']} 条教学记录")
            
            # 3. 迁移该学生的 recharge_orders
            mysql_cursor.execute('''
                SELECT id, student_id, hours, comment, create_time, active, type
                FROM orders
                WHERE student_id = %s
            ''', (student_id,))
            orders = mysql_cursor.fetchall()
            
            for order in orders:
                try:
                    # 根据 Snowflake 算法生成订单号（格式: ORD + Snowflake ID）
                    order_id = order['id']
                    snowflake_id = int(order.get('create_time').timestamp() * 1000) + order_id
                    order_number = f"ORD{snowflake_id}"
                    
                    # 如果学生有多个课程关系，优先匹配订单关联的教师
                    # 如果只有一个课程关系，使用该课程
                    student_subject_id = None
                    
                    if len(student_subjects_created) == 1:
                        # 只有一个课程关系，直接使用
                        student_subject_id = student_subjects_created[0][1]
                    elif len(student_subjects_created) > 1:
                        # 多个课程关系的情况
                        # 优先尝试匹配订单中记录的教师
                        # 如果无法确定，使用余额最高的课程
                        for tid, ss_id in student_subjects_created:
                            student_subject_id = ss_id
                            break  # 使用第一个作为默认
                    
                    if not student_subject_id:
                        stats['errors'].append(f"订单 {order_number} 找不到课程关联")
                        continue
                    
                    # 根据 hours 的正负判断订单类型（SQLite中不存储type，只存hours）
                    hours = order['hours']
                    
                    sqlite_cursor.execute('''
                        INSERT INTO recharge_orders 
                        (order_number, student_course_id, hours, amount, remark, created_at, updated_at, deleted_at)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    ''', (
                        order_number, student_subject_id, hours, 0.0,
                        order.get('comment'),
                        self._format_datetime_for_gorm(order.get('create_time')),
                        self._format_datetime_for_gorm(order.get('create_time')),
                        self._format_datetime_for_gorm(deleted_at)
                    ))
                    stats['recharge_orders'] += 1
                except Exception as order_error:
                    stats['errors'].append(f"订单 {order.get('id')} 处理失败: {str(order_error)}")
                    continue
            
            print(f"    [OK] 迁移 {stats['recharge_orders']} 条充值订单")
            
            # 提交该学生的所有数据
            self.sqlite_conn.commit()
            
            if stats['errors']:
                print(f"    [WARN] 警告: {len(stats['errors'])} 个问题")
                for error in stats['errors']:
                    print(f"      - {error}")
            
        except Exception as e:
            stats['errors'].append(str(e))
            print(f"    [ERROR] 迁移失败: {e}")
            print(f"    [DEBUG] 上下文信息:")
            print(f"      - 学生: {student_name} (ID: {student_id})")
            print(f"      - 余额: {balance}")
            print(f"      - 已创建课程关系: {len(student_subjects_created)}")
            print(f"      - 已迁移教学记录: {stats['records']}")
            print(f"      - 已迁移充值订单: {stats['recharge_orders']}")
            self.sqlite_conn.rollback()
        finally:
            mysql_cursor.close()
        
        return stats
    
    def run(self) -> bool:
        """执行迁移"""
        try:
            if not self.connect():
                return False
            
            print("\n" + "=" * 80)
            print("开始数据库迁移...")
            print("=" * 80)
            
            # 验证 SQLite 数据库表结构
            self.verify_schema()
            
            # 清空现有数据（如果需要）
            if self.clear_data:
                self.clear_existing_data()
            
            # 第一阶段：迁移基础表
            print("\n" + "─" * 80)
            print("第一阶段：迁移基础表")
            print("─" * 80)
            
            teachers_count = self.migrate_teachers()
            subjects_count = self.migrate_subjects()
            students_count = self.migrate_students()
            
            # 构建学生-教师-科目关系
            self.build_student_teacher_subject_relationships()
            
            # 第二阶段：逐个迁移学生数据
            print("\n" + "─" * 80)
            print("第二阶段：逐个迁移学生数据")
            print("─" * 80)
            
            # 获取所有学生列表
            mysql_cursor = self.mysql_conn.cursor(dictionary=True)
            mysql_cursor.execute("SELECT id, name FROM students ORDER BY id")
            students = mysql_cursor.fetchall()
            mysql_cursor.close()
            
            total_stats = {
                'students': len(students),
                'student_subjects': 0,
                'records': 0,
                'recharge_orders': 0,
                'failed_students': 0
            }
            
            for idx, student in enumerate(students, 1):
                print(f"\n[{idx}/{len(students)}]", end=" ")
                stats = self.migrate_student_data(student['id'], student['name'])
                
                total_stats['student_subjects'] += stats['student_subjects']
                total_stats['records'] += stats['records']
                total_stats['recharge_orders'] += stats['recharge_orders']
                
                if stats['errors']:
                    total_stats['failed_students'] += 1
            
            # 打印总结
            print("\n" + "=" * 80)
            print("[OK] 数据库迁移完成！")
            print("=" * 80)
            print(f"基础表迁移:")
            print(f"  - 教师: {teachers_count} 条")
            print(f"  - 科目: {subjects_count} 条")
            print(f"  - 学生: {students_count} 条")
            print(f"\n学生数据迁移:")
            print(f"  - 成功学生: {total_stats['students'] - total_stats['failed_students']}/{total_stats['students']}")
            print(f"  - 课程关系: {total_stats['student_subjects']} 条")
            print(f"  - 教学记录: {total_stats['records']} 条")
            print(f"  - 充值订单: {total_stats['recharge_orders']} 条")
            
            if total_stats['failed_students'] > 0:
                print(f"\n[WARN] 警告: {total_stats['failed_students']} 个学生迁移时出现问题")
            
            print(f"\nSQLite 数据库位置: {self.sqlite_path}")
            print("=" * 80)
            
            return True
            
        except Exception as e:
            print(f"\n[ERROR] 迁移失败: {e}")
            return False
        finally:
            self.close()


def main():
    """主函数"""
    import sys
    
    # MySQL 连接配置
    mysql_config = {
        'host': '49.235.136.31',  # 修改为实际的 MySQL 服务器地址
        'user': 'teaching_manage',       # 修改为实际的用户名
        'password': 'ZheShiHwx190800',       # 修改为实际的密码
        'database': 'teaching_manage'
    }
    
    # SQLite 数据库路径（使用现有数据库）
    script_dir = os.path.dirname(os.path.abspath(__file__))
    sqlite_path = os.path.join(os.path.dirname(script_dir), 'data', 'teaching_manage.db')
    
    # 检查 SQLite 数据库是否存在
    if not os.path.exists(sqlite_path):
        print(f"\n[ERROR] 错误: SQLite 数据库文件不存在: {sqlite_path}")
        print("请先运行应用程序以初始化数据库，或手动创建数据库文件。")
        sys.exit(1)
    
    print("\n" + "=" * 80)
    print("教学管理系统 - 数据库迁移工具")
    print("=" * 80)
    print(f"MySQL 服务器: {mysql_config['host']}")
    print(f"MySQL 数据库: {mysql_config['database']}")
    print(f"SQLite 数据库: {sqlite_path}")
    print(f"清空现有数据: 是")
    print("=" * 80 + "\n")
    
    # 创建迁移器并执行迁移
    # clear_data=True 表示迁移前清空现有数据
    # 如果想保留现有数据，改为 clear_data=False
    migrator = DatabaseMigrator(mysql_config, sqlite_path, clear_data=True)
    
    if migrator.run():
        print("\n[OK] 迁移成功！")
        print("\n[INFO] 注意事项:")
        print("  1. student_subjects 表是根据 records 和 orders 自动生成的")
        print("  2. records 表的 subject_id 是根据学生-教师关系推断的")
        print("  3. recharge_orders 已关联到具体的学生课程")
        print("  4. 如有数据不准确，请在应用中手动调整")
        print("\n[INFO] 后续步骤:")
        print("  1. 检查 SQLite 数据库中的数据是否正确")
        print("  2. 运行应用程序测试功能")
        print("  3. 验证课时余额和订单关联")
    else:
        print("\n[ERROR] 迁移失败！")
        sys.exit(1)


if __name__ == '__main__':
    main()

