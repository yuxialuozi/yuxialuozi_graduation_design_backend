# 租户信息管理系统 - 后端 API

基于 Go + Gin + GORM + PostgreSQL 的租户信息管理系统后端 API 服务。

## 项目概述

本系统是一个完整的租户信息管理平台后端服务，提供租户管理、合同管理、房间管理、费用管理、维修工单管理和报表统计等功能。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: PostgreSQL
- **依赖注入**: Wire v0.5.0
- **认证**: JWT (golang-jwt/jwt v5.2.0)
- **配置管理**: Viper v1.18.2
- **日志**: Zap v1.26.0
- **API 文档**: Swagger (swaggo/swag v1.16.2)

## 功能特性

### 用户认证
- JWT 用户登录认证
- 获取当前用户信息
- 用户登出

### 租户管理
- 租户 CRUD 操作
- 支持关键字搜索
- 支持状态筛选

### 合同管理
- 合同 CRUD 操作
- 支持关键字搜索
- 支持状态筛选
- 支持日期范围查询

### 房间管理
- 房间 CRUD 操作
- 支持关键字搜索
- 支持楼栋筛选
- 支持状态筛选
- 租户分配与释放

### 费用管理
- 费用记录 CRUD 操作
- 支持多条件筛选（租户、房间、费用类型、状态、账期）
- 缴费确认

### 维修工单管理
- 工单 CRUD 操作
- 支持多条件筛选
- 指派维修人员
- 完成工单

### 报表统计
- 收入统计（按月、按类型）
- 出租率统计
- 费用构成分析
- 维修统计数据
- 租户缴费排行榜
- 仪表盘汇总数据

## 项目结构

```
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── config/                  # 配置管理
│   │   └── config.go
│   ├── database/                # 数据库连接
│   │   └── postgres.go
│   ├── middleware/              # 中间件
│   │   ├── auth.go              # JWT 认证
│   │   ├── cors.go              # CORS 跨域
│   │   ├── logger.go            # 请求日志
│   │   └── recovery.go          # 错误恢复
│   ├── model/                   # GORM 模型
│   │   ├── user.go
│   │   ├── tenant.go
│   │   ├── contract.go
│   │   ├── room.go
│   │   ├── fee.go
│   │   └── maintenance.go
│   ├── repository/              # 数据访问层
│   │   ├── user_repo.go
│   │   ├── tenant_repo.go
│   │   ├── contract_repo.go
│   │   ├── room_repo.go
│   │   ├── fee_repo.go
│   │   └── maintenance_repo.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── tenant_service.go
│   │   ├── contract_service.go
│   │   ├── room_service.go
│   │   ├── fee_service.go
│   │   ├── maintenance_service.go
│   │   └── report_service.go
│   ├── handler/                 # HTTP 处理器
│   │   ├── auth_handler.go
│   │   ├── tenant_handler.go
│   │   ├── contract_handler.go
│   │   ├── room_handler.go
│   │   ├── fee_handler.go
│   │   ├── maintenance_handler.go
│   │   └── report_handler.go
│   ├── router/                  # 路由配置
│   │   └── router.go
│   ├── dto/                     # 数据传输对象
│   │   ├── request.go           # 请求 DTO
│   │   └── response.go          # 响应 DTO
│   └── wire/                    # Wire 依赖注入
│       ├── wire.go
│       └── wire_gen.go
├── pkg/
│   ├── response/                # 统一响应
│   │   └── response.go
│   └── utils/                   # 工具函数
│       └── jwt.go
├── docs/                        # Swagger 文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── config.yaml                  # 配置文件
├── go.mod
├── go.sum
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 12+

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/yuxialuozi/yuxialuozi_graduation_design_backend.git
   cd yuxialuozi_graduation_design_backend
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置数据库**

   创建 PostgreSQL 数据库：
   ```sql
   CREATE DATABASE tenant_management;
   ```

   修改 `config.yaml` 中的数据库配置：
   ```yaml
   database:
     host: localhost
     port: 5432
     user: postgres
     password: your_password
     dbname: tenant_management
     sslmode: disable
   ```

4. **运行项目**
   ```bash
   go run cmd/server/main.go
   ```

   项目启动后会自动创建数据库表。

### 默认管理员账号

- 用户名: `admin`
- 密码: `admin123`

## API 文档

启动项目后，访问 Swagger UI 查看完整的 API 文档：

```
http://localhost:8080/swagger/index.html
```

### 主要 API 端点

#### 认证 `/api/auth`

| 方法 | 路径    | 说明         |
|------|---------|--------------|
| POST | /login  | 用户登录     |
| GET  | /me     | 获取当前用户 |
| POST | /logout | 退出登录     |

#### 租户管理 `/api/tenants`

| 方法   | 路径 | 说明   | 查询参数                        |
|--------|------|--------|---------------------------------|
| GET    | /    | 租户列表 | page, pageSize, keyword, status |
| GET    | /:id | 租户详情 | -                               |
| POST   | /    | 创建租户 | -                               |
| PUT    | /:id | 更新租户 | -                               |
| DELETE | /:id | 删除租户 | -                               |

#### 合同管理 `/api/contracts`

| 方法   | 路径 | 说明   | 查询参数                                                    |
|--------|------|--------|-------------------------------------------------------------|
| GET    | /    | 合同列表 | page, pageSize, keyword, status, startDateFrom, startDateTo |
| GET    | /:id | 合同详情 | -                                                           |
| POST   | /    | 创建合同 | -                                                           |
| PUT    | /:id | 更新合同 | -                                                           |
| DELETE | /:id | 删除合同 | -                                                           |

#### 房间管理 `/api/rooms`

| 方法   | 路径        | 说明     | 查询参数                                  |
|--------|-------------|----------|-------------------------------------------|
| GET    | /           | 房间列表 | page, pageSize, keyword, building, status |
| GET    | /:id        | 房间详情 | -                                         |
| POST   | /           | 创建房间 | -                                         |
| PUT    | /:id        | 更新房间 | -                                         |
| DELETE | /:id        | 删除房间 | -                                         |
| POST   | /:id/assign | 分配租户 | {tenantId}                                |

#### 费用管理 `/api/fees`

| 方法   | 路径     | 说明     | 查询参数                                                  |
|--------|----------|----------|-----------------------------------------------------------|
| GET    | /        | 费用列表 | page, pageSize, tenantId, roomNo, feeType, status, period |
| GET    | /:id     | 费用详情 | -                                                         |
| POST   | /        | 创建费用 | -                                                         |
| PUT    | /:id     | 更新费用 | -                                                         |
| DELETE | /:id     | 删除费用 | -                                                         |
| POST   | /:id/pay | 确认缴费 | {paidDate?}                                               |

#### 维修管理 `/api/maintenance`

| 方法   | 路径          | 说明         | 查询参数                                        |
|--------|---------------|--------------|-------------------------------------------------|
| GET    | /             | 工单列表     | page, pageSize, keyword, type, status, priority |
| GET    | /:id          | 工单详情     | -                                               |
| POST   | /             | 创建工单     | -                                               |
| PUT    | /:id          | 更新工单     | -                                               |
| DELETE | /:id          | 删除工单     | -                                               |
| POST   | /:id/assign   | 指派维修人员 | {assignee}                                      |
| POST   | /:id/complete | 完成工单     | {completedAt?}                                  |

#### 报表统计 `/api/reports`

| 方法 | 路径               | 说明       | 查询参数            |
|------|--------------------|------------|---------------------|
| GET  | /income            | 收入统计   | start, end, groupBy |
| GET  | /occupancy         | 出租率统计 | start, end          |
| GET  | /fees/composition  | 费用构成   | start, end          |
| GET  | /maintenance/stats | 维修统计   | start, end          |
| GET  | /tenants/ranking   | 租户排行   | limit, start, end   |
| GET  | /dashboard         | 仪表盘数据 | -                   |

## 开发命令

### 安装依赖
```bash
go mod tidy
```

### 运行开发服务器
```bash
go run cmd/server/main.go
```

### 构建
```bash
go build -o server cmd/server/main.go
```

### 生成 Swagger 文档
```bash
swag init -g cmd/server/main.go -o docs
```

### 运行测试
```bash
go test ./...
```

## 配置说明

### config.yaml

```yaml
server:
  port: 8080            # 服务端口
  mode: debug           # 运行模式: debug, release

database:
  host: localhost       # 数据库主机
  port: 5432            # 数据库端口
  user: postgres        # 数据库用户名
  password: your_password  # 数据库密码
  dbname: tenant_management  # 数据库名称
  sslmode: disable      # SSL 模式

jwt:
  secret: your-jwt-secret-key-please-change-in-production  # JWT 密钥
  expire: 24h            # Token 过期时间

log:
  level: debug           # 日志级别
  format: json           # 日志格式
```

## 中间件

- **CORS**: 允许跨域访问（生产环境建议配置具体域名）
- **JWT Auth**: 基于 Token 的用户认证
- **Logger**: 请求日志记录
- **Recovery**: Panic 恢复，防止服务崩溃

## 数据模型

### User 用户表
- 字段: ID, Username, Password, Nickname, Avatar, Role, Permissions
- 默认角色: admin, user

### Tenant 租户表
- 字段: ID, Name, ContactPerson, Phone, Email, Status
- 状态: active, inactive

### Contract 合同表
- 字段: ID, TenantID, ContractNo, StartDate, EndDate, Amount, Status
- 状态: draft, active, expired, terminated

### Room 房间表
- 字段: ID, RoomNo, Building, Floor, Area, MonthlyRent, Status, TenantID
- 状态: vacant, occupied, maintenance

### Fee 费用表
- 字段: ID, TenantID, RoomNo, FeeType, Amount, Period, DueDate, PaidDate, Status
- 费用类型: rent, water, electricity, property, other

### Maintenance 维修工单表
- 字段: ID, TicketNo, TenantID, RoomNo, Type, Description, Priority, Status, Assignee
- 类型: electrical, plumbing, appliance, furniture, other
- 状态: pending, processing, completed, cancelled
- 优先级: low, medium, high, urgent

## 统一响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "请求参数错误",
  "data": null
}
```

### 分页响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "pageSize": 10
  }
}
```

## 代码规范

1. 使用分层架构: Handler → Service → Repository → Model
2. 统一错误处理和响应格式
3. 使用 DTO 进行数据传输，避免直接暴露模型
4. Repository 层只处理数据库操作，Service 层处理业务逻辑
5. 使用 Wire 管理依赖注入
6. 敏感配置使用环境变量或配置文件

## 待实现功能

- [ ] 用户管理（增删改查）
- [ ] 文件上传功能（房间图片、租户证件等）
- [ ] 邮件/短信通知功能
- [ ] 定时任务（自动计算月度费用、过期提醒）
- [ ] 数据导入导出
- [ ] 操作日志记录
- [ ] 权限细粒度控制

## License

Apache 2.0

## 联系方式

- 作者: yuxialuozi
- 项目地址: https://github.com/yuxialuozi/yuxialuozi_graduation_design_backend