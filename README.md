# 教学管理系统 (Teaching Manage)

一个基于 **Wails** 框架开发的教学管理桌面应用，结合 **Go** 后端和 **Vue 3 + TypeScript** 前端，提供学生管理、课程管理、教师管理、消课记录、充值订单等功能。

## ✨ 特性

- 🖥️ **跨平台桌面应用** - 基于 Wails v2 构建，支持 Windows、macOS、Linux
- 🎨 **现代化 UI** - 使用 Vue 3 + Vuetify 3 构建优雅的用户界面
- 🔄 **命令分发架构** - 后端采用 Command Dispatcher 模式统一处理前后端通信
- 💾 **数据持久化** - 使用 GORM + SQLite 实现本地数据存储
- ☁️ **云端备份** - 支持 WebDAV 协议的云端数据备份与恢复
- 📊 **数据可视化** - 集成 ECharts 实现数据统计图表展示
- 🔐 **数据安全** - 支持自动备份功能

## 程序截图
<img width="2560" height="1440" alt="仪表盘" src="https://github.com/user-attachments/assets/7417cfe6-4a18-4d73-99c2-40d7d4ca111a" />
<img width="2560" height="1440" alt="数据备份" src="https://github.com/user-attachments/assets/7be94878-aa65-4196-8e44-27ccaae02bc7" />

## 🏗️ 架构设计

### 后端架构 (Go)

```
backend/
├── service/          # 业务逻辑层
│   ├── student_manager.go
│   ├── teacher_manager.go
│   ├── course_manage.go
│   ├── record_manager.go
│   ├── order_manager.go
│   ├── backup_manager.go
│   └── ...
├── repository/       # 数据访问层
├── dao/              # 数据访问对象
├── entity/           # 数据库实体
├── model/            # 业务模型
└── pkg/              # 公共包
    ├── dispatcher/   # 命令分发器
    ├── logger/       # 日志系统
    ├── crypto/       # 加密工具
    └── ...
```

**核心设计模式**：
- 采用 **Command Dispatcher** 模式，前端通过 `Dispatch` 方法统一调用后端服务
- 服务通过 `RegisterRoute` 方法注册路由，如 `student_manager/get_student_list`
- 支持依赖注入，所有服务在 `main.go` 中统一初始化和装配

### 前端架构 (Vue 3)

```
frontend/src/
├── views/            # 页面组件
│   ├── StudentManage/
│   │   ├── StudentManage.vue
│   │   └── StudentManage.logic.ts  # 业务逻辑分离
│   ├── CourseManage/
│   ├── Settings/
│   └── ...
├── api/              # API 封装
├── components/       # 可复用组件
├── composables/      # 组合式函数
├── router/           # 路由配置
└── types/            # TypeScript 类型定义
```

**设计原则**：
- **逻辑分离**：复杂组件的逻辑提取到 `*.logic.ts` 文件中
- **统一 API**：所有后端调用通过 `api/` 目录封装
- **类型安全**：完整的 TypeScript 类型定义

## 🚀 快速开始

### 环境要求

- **Go** >= 1.24
- **Node.js** >= 16.x
- **npm** 或 **pnpm**
- **Wails CLI** v2.11.0+

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 克隆项目

```bash
git clone https://github.com/Xiaoxin-19/teaching_manage.git
cd teaching_manage
```

### 安装依赖

Wails 会自动安装前端依赖，无需手动安装。

### 开发模式

```bash
wails dev
```

运行后会启动：
- Vite 开发服务器（热重载）
- Go 后端服务
- 桌面应用窗口

### 构建生产版本

```bash
wails build
```

或使用 VS Code 任务：
- `build` - 生产版本
- `build debug` - 带调试信息的版本
- `build dev` - 开发版本

构建产物位于 `build/bin/` 目录。

## 📦 技术栈

### 后端
- **[Wails v2](https://wails.io/)** - Go 桌面应用框架
- **[GORM](https://gorm.io/)** - ORM 框架
- **[SQLite](https://www.sqlite.org/)** - 嵌入式数据库
- **[Zap](https://github.com/uber-go/zap)** - 高性能日志库
- **[Validator](https://github.com/go-playground/validator)** - 数据验证
- **[GoWebDAV](https://github.com/studio-b12/gowebdav)** - WebDAV 客户端
- **[Excelize](https://github.com/qax-os/excelize)** - Excel 文档处理

### 前端
- **[Vue 3](https://vuejs.org/)** - 渐进式 JavaScript 框架
- **[TypeScript](https://www.typescriptlang.org/)** - 类型安全的 JavaScript
- **[Vuetify 3](https://vuetifyjs.com/)** - Material Design 组件库
- **[Vue Router](https://router.vuejs.org/)** - 官方路由管理器
- **[ECharts](https://echarts.apache.org/)** - 数据可视化库
- **[Vite](https://vitejs.dev/)** - 新一代前端构建工具

## 🔧 开发指南

### 添加新功能（后端）

1. 在 `backend/service/request` 和 `backend/service/response` 定义请求/响应结构
2. 在对应的 Service 中实现业务方法
3. 通过 `RegisterRoute` 注册路由：

```go
func (sm *StudentManager) RegisterRoute(d *dispatcher.Dispatcher) {
    dispatcher.RegisterTyped(d, "student_manager/my_method", sm.MyMethod)
}
```

4. 在 `main.go` 中确保服务已注册到 Dispatcher

### 添加新功能（前端）

1. 在 `src/api/` 创建 API 封装：

```typescript
export async function GetStudentList(req: GetStudentListRequest): Promise<StudentListResponse> {
  const result = await Dispatch('student_manager/get_student_list', JSON.stringify(req))
  return JSON.parse(result)
}
```

2. 创建页面组件和逻辑文件：
   - `MyView.vue` - 页面模板
   - `MyView.logic.ts` - 业务逻辑

3. 在 `router/` 中配置路由

### 项目配置

- `wails.json` - Wails 项目配置
- `frontend/vite.config.ts` - Vite 构建配置
- `backend/wirex/wirex.go` - 依赖注入配置

## 📝 功能模块

- **学生管理** - 学生信息的增删改查、状态管理
- **教师管理** - 教师信息维护
- **课程管理** - 课程安排、学生选课
- **消课记录** - 课时消耗记录与统计
- **充值订单** - 学生充值记录管理
- **数据备份** - 本地备份和云端 WebDAV 备份
- **数据统计** - 可视化数据报表
- **系统设置** - 应用配置和参数管理

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 👨‍💻 作者

**HeWenxin**
- Email: he_wenxin@foxmail.com
- GitHub: [@Xiaoxin-19](https://github.com/Xiaoxin-19)

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的 Go 桌面应用框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Vuetify](https://vuetifyjs.com/) - Material Design 组件库
