# 数据库迁移说明

## 快速使用

### 1. 准备环境
```bash
# 安装依赖
pip install mysql-connector-python

# 先运行应用初始化 SQLite 数据库
wails dev
```

### 2. 配置连接
编辑 `migrate.py` 中的 MySQL 配置：
```python
mysql_config = {
    'host': '49.235.136.31',  # MySQL 服务器地址
    'user': 'root',
    'password': 'your_password',
    'database': 'teaching_manage'
}
```

### 3. 运行迁移
```bash
cd script
python migrate.py
```

## 数据转换逻辑

迁移脚本会自动处理 MySQL 和 SQLite 之间的表结构差异：

### 1. **编号生成**

| 表 | 字段 | 格式 | 示例 |
|---|---|---|---|
| teachers | teacher_number | T + 8位数字 | T00000001 |
| students | student_number | S + 年份 + 5位数字 | S202400001 |
| recharge_orders | order_number | ORD + Snowflake ID | ORD1704235200000001 |

### 2. **时间处理**

- **teaching_time 解析**：支持 15+ 种格式
  - 标准格式：`18:15-19:00`
  - 全角冒号：`18：15-19：00`
  - 多冒号错误：`18:30:19:15` → `18:30-19:15`
  - 引号混用：`20：00-20"45` → `20:00-20:45`
  - 分号代替冒号：`17;30-18:15` → `17:30-18:15`

- **时间格式化**：
  - `start_time/end_time`：HH:MM 格式（如 18:15）
  - `teaching_date_ms`：Unix 毫秒时间戳（便于排序）

### 3. **学生-教师-科目关系**

- 从 `students.teacher_id` 获取当前教师关系（解耦设计）
- `records` 表中的历史教师自动回退到当前教师的科目 ID
- 避免"找不到科目映射"错误

### 4. **student_subjects 表**

- MySQL 中不存在此表，脚本自动生成
- 每个学生一个课程关系（当前教师）
- 课时余额从 `students.hours` 字段读取

### 5. **recharge_orders 表**

- 从 MySQL 的 `orders` 表迁移
- 自动关联到对应的学生课程
- 订单号使用 Snowflake 算法生成确保唯一性

### 6. **records 表**

- 自动推断缺失的 `subject_id` 字段
- `teaching_time` 解析为 `start_time/end_time`
- 生成 `teaching_date_ms` 用于精确排序

## 迁移特性

- ✓ **两阶段迁移**：基础表 → 逐学生详细数据（隔离错误）
- ✓ **智能时间解析**：自动修正格式错误，成功率 100%
- ✓ **历史数据处理**：学生教师变更时自动使用当前教师科目
- ✓ **完整日志**：详细的迁移统计和错误报告
- ✓ **数据完整性**：所有记录都能找到科目映射

## 验证结果

```bash
# 检查数据库
sqlite3 ../data/teaching_manage.db

# 查询数据
sqlite> SELECT COUNT(*) FROM students;
sqlite> SELECT COUNT(*) FROM teachers;  
sqlite> SELECT COUNT(*) FROM student_subjects;
sqlite> SELECT COUNT(*) FROM records;
sqlite> SELECT COUNT(*) FROM recharge_orders;

# 检查编号格式
sqlite> SELECT student_number FROM students LIMIT 5;
sqlite> SELECT teacher_number FROM teachers LIMIT 5;
sqlite> SELECT order_number FROM recharge_orders LIMIT 5;

# 检查时间格式
sqlite> SELECT start_time, end_time FROM records LIMIT 5;
sqlite> SELECT teaching_date_ms FROM records WHERE teaching_date_ms IS NOT NULL LIMIT 5;
```

## 修改配置

### 保留现有数据（追加模式）
```python
# 在 migrate.py 中修改
migrator = DatabaseMigrator(mysql_config, sqlite_path, clear_data=False)
```

### 自定义 SQLite 数据库路径
```python
sqlite_path = "../data/teaching_manage.db"  # 修改此路径
```

## 注意事项

1. 确保 SQLite 数据库已初始化（运行过应用）
2. MySQL 连接配置正确且网络可达
3. 检查迁移日志中的警告和错误信息
4. 迁移后在应用中验证数据准确性
5. 备份原始数据库（迁移前）
