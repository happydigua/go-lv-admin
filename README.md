# Go Lv Admin

一个基于 Golang (Gin) + Vue3 (Naive UI) 的全栈后台管理系统。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D.svg)
![Gorm](https://img.shields.io/badge/Gorm-v2-red.svg)

## ✨ 特性

- **用户管理**：用户增删改查、角色分配、密码重置。
- **角色管理**：基于 Casbin 的 RBAC 权限控制，支持菜单权限分配。
- **菜单管理**：动态路由菜单，支持无限层级。
- **系统设置**：
    - **基础设置**：可视化配置系统名称、Logo、版权信息（数据库存储）。
    - **存储配置**：支持 Local、Aliyun OSS、Tencent COS、Cloudflare R2 等多种存储驱动（配置文件）。
- **文件管理**：统一的文件上传和管理界面，支持多种存储后端。
- **操作日志**：全系统操作审计，支持请求详情查看。
- **代码生成器**：一键生成前后端 CRUD 代码，**自动写入文件并注册菜单**：
    - 自动生成 Model、Service、API 后端代码
    - 自动追加路由到 `router.go`
    - 自动创建数据库菜单记录
    - 自动生成 Vue 页面和前端 API
    - 支持选择父菜单和设置图标
- **仪表盘**：ECharts 可视化展示系统运行状态。

## 🛠️ 技术栈

### 后端 (Backend)
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [Gorm](https://gorm.io/) (MySQL)
- **权限控制**: [Casbin](https://casbin.org/)
- **配置管理**: Viper
- **日志**: Zap
- **文档**: Swagger

### 前端 (Frontend)
- **核心框架**: Vue 3 + Vite + TypeScript
- **UI 组件库**: [Naive UI](https://www.naiveui.com/)
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP**: Axios
- **图表**: Apache ECharts

## 🚀 快速开始

### 环境要求
- Golang 1.21+
- Node.js 16+
- MySQL 8.0+

### 1. 数据库初始化
**方式一：自动初始化**
创建数据库 `go_lv_vue_admin`，后端启动时会自动迁移表结构并初始化默认数据（用户、角色、菜单、设置）。

**方式二：手动导入 SQL（推荐）**
如果需要快速部署或迁移，可以直接导入项目提供的 SQL 文件：
```bash
mysql -u root -p -h 127.0.0.1 < deploy/sql/go_lv_vue_admin.sql
```

### 2. 后端启动
```bash
cd backend

# 安装依赖
go mod tidy

# 修改配置 (config/config.yaml)
# 确保数据库配置正确

# 启动服务
go run cmd/server/main.go
```
后端服务将运行在 `http://127.0.0.1:8888`

### 3. 前端启动
```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```
前端服务将运行在 `http://localhost:5173`

## 📁 目录结构

```
go-lv-admin/
├── backend/            # 后端项目
│   ├── cmd/            # 入口文件
│   ├── config/         # 配置文件
│   ├── internal/       # 业务逻辑
│   │   ├── api/        # 控制器
│   │   ├── core/       # 核心组件
│   │   ├── middleware/ # 中间件
│   │   ├── model/      # 数据模型
│   │   ├── router/     # 路由配置
│   │   ├── service/    # 业务服务
│   │   └── storage/    # 存储驱动
│   └── ...
└── frontend/           # 前端项目
    ├── src/
    │   ├── api/        # 接口定义
    │   ├── components/ # 公共组件
    │   ├── layouts/    # 布局组件
    │   ├── store/      # 状态管理
    │   ├── views/      # 页面视图
    │   └── ...
    └── ...
```

## 📄 License

[MIT](./LICENSE)
