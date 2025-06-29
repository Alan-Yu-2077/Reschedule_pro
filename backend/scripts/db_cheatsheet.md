# 数据库操作速查表

## 基本CRUD操作

### 创建 (Create)
```go
// 创建单个记录
user := &models.User{Username: "test", Password: "123"}
database.DB.Create(user)

// 批量创建
users := []models.User{
    {Username: "user1", Password: "123"},
    {Username: "user2", Password: "123"},
}
database.DB.Create(&users)

// 创建关联记录
schedule := &models.ClassSchedule{
    ClassName: "班级1",
    Courses: []models.Course{
        {Name: "数学", Day: 0, Slot: 0, WeekFrom: 1, WeekTo: 16},
    },
}
database.DB.Create(schedule)
```

### 查询 (Read)
```go
// 查询所有记录
var users []models.User
database.DB.Find(&users)

// 根据条件查询
var user models.User
database.DB.Where("username = ?", "test").First(&user)

// 查询单个记录
database.DB.First(&user, 1) // 根据ID查询

// 查询关联数据
var schedule models.ClassSchedule
database.DB.Preload("Courses").First(&schedule, 1)

// 分页查询
var users []models.User
database.DB.Limit(10).Offset(20).Find(&users)
```

### 更新 (Update)
```go
// 更新单个记录
user.Password = "newpassword"
database.DB.Save(&user)

// 更新特定字段
database.DB.Model(&user).Update("password", "newpassword")

// 批量更新
database.DB.Model(&models.User{}).Where("username LIKE ?", "test%").Update("password", "newpass")

// 更新多个字段
updates := map[string]interface{}{
    "password": "newpass",
    "username": "newname",
}
database.DB.Model(&user).Updates(updates)
```

### 删除 (Delete)
```go
// 软删除（默认）
database.DB.Delete(&user)

// 永久删除
database.DB.Unscoped().Delete(&user)

// 批量删除
database.DB.Where("username LIKE ?", "test%").Delete(&models.User{})
```

## 高级查询

### 条件查询
```go
// 基本条件
database.DB.Where("age > ?", 18).Find(&users)

// 多个条件
database.DB.Where("age > ? AND username LIKE ?", 18, "test%").Find(&users)

// IN查询
database.DB.Where("id IN ?", []int{1, 2, 3}).Find(&users)

// 模糊查询
database.DB.Where("username LIKE ?", "%test%").Find(&users)
```

### 排序
```go
// 升序
database.DB.Order("created_at").Find(&users)

// 降序
database.DB.Order("created_at desc").Find(&users)

// 多字段排序
database.DB.Order("age desc, username asc").Find(&users)
```

### 聚合查询
```go
// 计数
var count int64
database.DB.Model(&models.User{}).Count(&count)

// 求和
var total int64
database.DB.Model(&models.Course{}).Select("SUM(week_to)").Scan(&total)

// 平均值
var avg float64
database.DB.Model(&models.Course{}).Select("AVG(slot)").Scan(&avg)

// 最大值/最小值
var maxWeek int
database.DB.Model(&models.Course{}).Select("MAX(week_to)").Scan(&maxWeek)
```

### 关联查询
```go
// 预加载关联数据
database.DB.Preload("Courses").Find(&schedules)

// 带条件的预加载
database.DB.Preload("Courses", "week_from <= ?", 8).Find(&schedules)

// 关联查询
database.DB.Joins("JOIN courses ON courses.schedule_id = class_schedules.id").Find(&schedules)
```

### 原生SQL
```go
// 原生SQL查询
var result []map[string]interface{}
database.DB.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&result)

// 执行原生SQL
database.DB.Exec("UPDATE users SET password = ? WHERE id = ?", "newpass", 1)
```

## 事务操作
```go
// 开始事务
tx := database.DB.Begin()

// 在事务中执行操作
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&schedule).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
tx.Commit()
```

## 常用技巧

### 错误处理
```go
if err := database.DB.Create(&user).Error; err != nil {
    log.Printf("创建用户失败: %v", err)
    return err
}
```

### 检查记录是否存在
```go
var count int64
database.DB.Model(&models.User{}).Where("username = ?", "test").Count(&count)
if count > 0 {
    fmt.Println("用户已存在")
}
```

### 获取最后插入的ID
```go
user := &models.User{Username: "test", Password: "123"}
database.DB.Create(user)
fmt.Printf("新用户ID: %d\n", user.ID)
```

### 批量操作
```go
// 批量插入
var users []models.User
for i := 0; i < 100; i++ {
    users = append(users, models.User{
        Username: fmt.Sprintf("user%d", i),
        Password: "123",
    })
}
database.DB.Create(&users)
```

## 数据库连接
```go
// 初始化数据库
database.InitDB()

// 获取数据库实例
db := database.DB

// 关闭数据库连接（通常在程序结束时）
sqlDB, err := database.DB.DB()
if err == nil {
    sqlDB.Close()
}
``` 