# 课程表管理系统 (Schedule Management System)

一个基于Go Gin后端和uni-app Vue3前端的现代化课程表管理系统，支持多用户权限管理、课程调换和跨周课程移动功能。

## 🚀 功能特性

### 用户管理
- **多用户类型支持**：admin（管理员）、user（普通用户）、viewer（只读用户）
- **用户注册**：支持10位数字UserID作为主键
- **权限控制**：不同用户类型拥有不同的操作权限
- **管理员功能**：用户增删改查、密码修改、用户信息编辑

### 课程表管理
- **班级管理**：支持多个班级的课程表管理
- **周数切换**：支持1-20周的课程表查看和编辑
- **课程调换**：支持同一周内课程位置调换
- **跨周移动**：支持将课程从一周移动到另一周
- **课程删除**：支持删除指定时间槽的课程
- **实时更新**：操作后自动刷新课程表和活动日志

### 界面特性
- **响应式设计**：适配不同屏幕尺寸
- **直观交互**：点击式课程选择和移动
- **视觉反馈**：选中课程高亮显示，操作状态清晰
- **英文化界面**：所有界面元素使用英文显示

## 🛠 技术栈

### 后端 (Backend)
- **语言**：Go 1.21+
- **框架**：Gin Web Framework
- **数据库**：SQLite (GORM)
- **中间件**：CORS支持
- **架构**：MVC模式

### 前端 (Frontend)
- **框架**：uni-app Vue3
- **语言**：TypeScript
- **构建工具**：Vite
- **UI组件**：原生HTML/CSS
- **状态管理**：Vue3 Composition API

## 📁 项目结构

```
Reschedule_program/
├── backend/                 # 后端代码
│   ├── database/           # 数据库配置
│   ├── middleware/         # 中间件
│   ├── models/            # 数据模型
│   ├── routes/            # API路由
│   ├── services/          # 业务逻辑
│   ├── scripts/           # 数据库脚本
│   └── main.go           # 主程序入口
├── frontend/              # 前端代码
│   └── res_pro/          # uni-app项目
│       ├── src/
│       │   ├── pages/    # 页面组件
│       │   │   ├── admin/    # 管理员页面
│       │   │   ├── login/    # 登录注册页面
│       │   │   ├── Main_page/ # 主功能页面
│       │   │   └── viewer/   # 只读用户页面
│       │   └── static/   # 静态资源
│       └── package.json
└── README.md
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- npm 或 yarn

### 安装步骤

#### 1. 克隆项目
```bash
git clone <repository-url>
cd Reschedule_program
```

#### 2. 启动后端服务
```bash
cd backend
go mod tidy
go run main.go
```
后端服务将在 `http://localhost:8080` 启动

#### 3. 启动前端服务
```bash
cd frontend/res_pro
npm install
npm run dev
```
前端服务将在 `http://localhost:5173` 启动

### 默认账户
- **管理员账户**：
  - UserID: Admin
  - 密码: 88888888
- **新注册用户**：默认为viewer类型

## 📖 使用指南

### 用户登录
1. 访问前端页面
2. 选择登录或注册
3. 输入UserID和密码
4. 根据用户类型跳转到相应页面

### 管理员功能
1. 使用管理员账户登录
2. 进入管理员仪表板
3. 管理用户：增删改查、修改密码、编辑用户信息
4. 查看系统活动日志

### 课程表操作
1. 选择班级和周数
2. 查看当前课程表
3. 课程调换操作：
   - 点击要移动的课程
   - 切换到目标周
   - 点击"Move to Current Week"
   - 点击目标位置完成移动
4. 删除课程：选中课程后点击"Delete Course"

### 只读用户
- viewer类型用户只能查看课程表
- 无法进行编辑操作
- 界面为只读模式

## 🔧 API文档

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 用户管理 (管理员)
- `GET /api/admin/users` - 获取所有用户
- `POST /api/admin/users` - 创建用户
- `PUT /api/admin/users/:id` - 更新用户信息
- `DELETE /api/admin/users/:id` - 删除用户

### 课程表管理
- `POST /api/schedule/save` - 保存课程表
- `GET /api/schedule/class/:className/week/:weekNumber` - 获取指定班级和周数的课程表
- `GET /api/schedule/classes` - 获取所有班级
- `DELETE /api/schedule/delete` - 删除课程
- `POST /api/schedule/move` - 移动课程（支持跨周）

## 🗄 数据库设计

### 主要表结构
- **users**: 用户信息表
- **classes**: 班级信息表
- **courses**: 课程信息表
- **weekly_schedules**: 周课程表

### 数据关系
- 一个班级可以有多个周课程表
- 一个课程可以在多个时间槽出现
- 用户通过UserID关联

## 🔒 安全特性

- **CORS配置**：支持跨域请求
- **输入验证**：所有API输入都经过验证
- **权限控制**：基于用户类型的操作权限
- **数据完整性**：数据库约束和事务处理

## 🐛 故障排除

### 常见问题

1. **端口占用**
   ```bash
   # 查找占用端口的进程
   lsof -i :8080
   # 杀死进程
   kill -9 <PID>
   ```

2. **数据库连接问题**
   ```bash
   # 重新初始化数据库
   cd backend
   go run scripts/init_data.go
   ```

3. **前端构建问题**
   ```bash
   # 清理缓存
   npm run clean
   npm install
   ```

### 日志查看
- 后端日志：控制台输出
- 前端日志：浏览器开发者工具

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📝 开发计划

### 已完成功能
- ✅ 用户注册登录系统
- ✅ 多用户权限管理
- ✅ 课程表基础CRUD操作
- ✅ 课程调换功能
- ✅ 跨周课程移动
- ✅ 管理员用户管理界面
- ✅ 只读用户界面

### 计划功能
- 🔄 课程冲突检测
- 🔄 批量课程操作
- 🔄 课程表导入导出
- 🔄 移动端适配优化
- 🔄 实时协作功能

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发送邮件至：[your-email@example.com]

---

**注意**：本项目仍在积极开发中，功能可能会有所变化。建议定期查看更新日志。 