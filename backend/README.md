# 课程调度系统 - 后端

## 数据库设置

本项目使用SQLite数据库来持久化存储数据。

### 安装依赖

```bash
go mod tidy
```

### 运行应用

```bash
go run main.go
```

### 初始化示例数据

```bash
go run scripts/init_data.go
```

这将创建：
- 示例用户：admin/admin123
- 示例课程表：Grade 23 - Class 1
- 示例活动日志

## 数据库结构

### 表结构

1. **users** - 用户表
   - id (主键)
   - username (用户名，唯一)
   - password (密码)
   - created_at, updated_at, deleted_at

2. **class_schedules** - 班级课程表
   - id (主键)
   - class_name (班级名称，唯一)
   - created_at, updated_at, deleted_at

3. **courses** - 课程表
   - id (主键)
   - name (课程名称)
   - day (星期几，0-6)
   - slot (时间段，0-4)
   - week_from (开始周)
   - week_to (结束周)
   - schedule_id (外键，关联班级课程表)
   - created_at, updated_at, deleted_at

4. **activity_logs** - 活动日志表
   - id (主键)
   - message (日志消息)
   - created_at, updated_at, deleted_at

## API端点

### 用户认证
- `POST /register` - 用户注册
- `POST /login` - 用户登录

### 课程调度
- `POST /schedule` - 创建课程表
- `GET /schedule/:class` - 获取指定班级的课程表
- `GET /schedules` - 获取所有课程表
- `GET /logs` - 获取活动日志

## 数据文件

数据库文件将自动创建为 `reschedule.db`，位于项目根目录。 