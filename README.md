# gozero 入门指南1

## 一、可以用 goctl 自动生成的步骤 ✅

### 1. 创建新项目（仅适用于从零开始）⭐

```bash
# 在项目根目录执行（例如 D:\Docker\code）
goctl api new test
```

**使用场景：**  在**空目录**或**新项目**开始时使用，会创建完整的项目结构。

### 2.编写 API 定义文件

编辑 `api/admin.api`，定义接口：

```go
syntax = "v1"

info (
	title:  "test"
	desc:   "test api"
	author: "your name"
	email:  "your email"
)

type GetUserRequest {
	Id int64 `path:"userId"`
}

type GetUserResponse {
	UserName string `json:"user_name"`
	Account     string `json:"account"`
	Status   int    `json:"status"`
}

type CreateUserRequest {
	UserName string `json:"user_name"`
	Account     string `json:"account"`
	Status   int    `json:"status"`
}

type CreateUserResponse {
	Result bool `json:"result"`
}

service admin-api {
	@handler GetUser
	get /users/id/:userId (GetUserRequest) returns (GetUserResponse)

	@handler CreateUser
	post /users/create (CreateUserRequest) returns (CreateUserResponse)
}
```

 **⚠️ 注意：**  如果项目已存在，**不要使用**此命令，它会创建新的子目录。

### 3. 生成 Model（如果还没有）

```bash
# 从 SQL 文件生成 model
# 添加 -cache 参数可以让生成的 admin_users model 具备自动缓存能力，优化数据访问性能
goctl model mysql ddl -src model/admin_users/admin_users.sql -dir model/admin_users -cache

# 或者从数据库直接生成
goctl model mysql datasource -url="root:root@tcp(127.0.0.1:3306)/hh_overseas" -table="admin_users" -dir="model/admin_users" -cache
```

**生成的文件：**

- ​`model/admin_users/adminusersmodel_gen.go` - 自动生成的模型代码
- ​`model/admin_users/adminusersmodel.go` - 可自定义的模型代码
- ​`model/admin_users/vars.go` - 变量定义

---

### 4. 从 API 文件生成代码（⚠️ 注意：会覆盖现有文件）

```bash
# 生成 handler、logic、types、routes 等
goctl api go -api api/admin.api -dir . --style gozero
```

**自动生成的文件包括：**

- ✅ `internal/types/types.go` - 类型定义（GetRoleRequest, GetRoleResponse, CreateRoleRequest, CreateRoleResponse）
- ✅ `internal/handler/adminrolehandler.go` - Handler 文件（GetRoleHandler, CreateRoleHandler）
- ✅ `internal/logic/adminrolelogic.go` - Logic 文件（GetRoleLogic, CreateRoleLogic）
- ✅ `internal/handler/routes.go` - 路由注册（会合并到现有路由中）

---

## 二、需要手动配置的步骤 ⚙️

### 1. 配置数据库连接（config.go）

在 `internal/config/config.go` 中添加：

```go
package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DataSource string          // 数据库连接字符串
	Cache      cache.CacheConf // 缓存配置（可选）
}
```

### 2. 初始化 Model（servicecontext.go）

在 `internal/svc/servicecontext.go` 中：

```go
package svc

import (
	"test/internal/config"
	"test/model/admin_users"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	AdminUsersModel admin_users.AdminUsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:          c,
		AdminUsersModel: admin_users.NewAdminUsersModel(sqlConn, c.Cache),
	}
}
```

### 3. 配置文件（etc/myproject.yaml）

添加数据库配置：

```yaml
# 数据库配置
DataSource: root:root@tcp(127.0.0.1:3306)/hh_overseas?parseTime=true

# Redis 缓存配置（可选）
Cache:
  - Host: 127.0.0.1:6379
    Pass:
    Type: node
```

### 4. 实现业务逻辑（logic 文件）

在 `internal/logic/getuserlogic.go` 中实现具体的数据库操作：

```go
package logic

import (
	"context"

	"test/internal/svc"
	"test/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	user, err := l.svcCtx.AdminUsersModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}
	return &types.GetUserResponse{
		UserName: user.UserName,
		Account:  user.Account,
		Status:   int(user.Status),
	}, nil
}
```

在 `internal/logic/createuserlogic.go` 中实现具体的数据库操作：

```go
package logic

import (
	"context"

	"test/internal/svc"
	"test/internal/types"
	"test/model/admin_users"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) (resp *types.CreateUserResponse, err error) {
	_, err = l.svcCtx.AdminUsersModel.Insert(l.ctx, &admin_users.AdminUsers{
		UserName: req.UserName,
		Account:  req.Account,
		Status:   int64(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateUserResponse{
		Result: true,
	}, nil
}
```

### 5. 修复模块依赖

如果 model 有独立的 go.mod，需要删除：

```bash
# 删除 model 下的 go.mod（如果有）
rm model/admin_users/go.mod

# 更新依赖
go mod tidy
```

---

## 三、验证步骤 ✅

1. **编译项目**

    ```bash
    go build -o test.exe .
    ```
2. **运行服务**

    ```bash
    ./test.exe -f etc/test-api.yaml
    # 或
    go run myproject.go -f etc/test-api.yaml
    ```
3. **测试接口**

    ```bash
    # 获取用户
    curl http://localhost:8888/users/id/2

    # 创建角色
    curl -X POST http://localhost:8888/users/create \
      -H "Content-Type: application/json" \
      -d '{"user_name":"张三","desc":"测试用户","status":1}'
    ```

---

## 四、推荐的工作流程 🔄

### 方案 A：从零开始创建新项目（推荐）

**第一步：创建项目结构**

```bash
# 在项目根目录执行，会创建 test 文件夹
goctl api new test
cd test
```

这会自动生成：

- 完整的项目目录结构
- ​`test.go` - 主入口文件
- ​`etc/test-api.yaml` - 配置文件
- `internal/` - 内部代码目录
- ​`api/admin.api` - API 定义文件模板

**第二步：编写 API 定义**  
编辑 `api/admin.api` 文件，定义接口

**第三步：生成代码**

```bash
# 从 .api 文件生成 handler、logic、types、routes
goctl api go -api api/admin.api -dir . --style gozero
```

**第四步：生成 Model**

```bash
# 从 SQL 文件生成 model
goctl model mysql ddl -src model/admin_users/admin_users.sql -dir model/admin_users -cache
```

**第五步：配置数据库**

- 修改 `internal/config/config.go` 添加 DataSource
- 修改 `internal/svc/servicecontext.go` 初始化 Model
- 修改 `etc/test-api.yaml` 添加数据库连接

**第六步：实现业务逻辑**
在 `internal/logic/` 中实现具体的业务逻辑

---

### 方案 B：在已有项目中添加新 API（当前情况）

**注意：**  `goctl api new` **不适用**于已有项目，它会在新目录创建完整项目。

**正确的步骤：**

1. ✅ 已有 `.api`​ 文件（`api/admin.api`）
2. ✅ 已有 model（已生成）
3. ⚠️ **使用 **​**​`goctl api go`​** 生成代码（会覆盖现有文件）
    ```bash
    # 方式1：生成到临时目录，然后手动合并
    goctl api go -api api/admin.api -dir ./temp_admin --style gozero
    # 然后手动复制需要的文件

    # 方式2：直接生成（会覆盖 routes.go、types.go 等）
    goctl api go -api api/admin.api -dir . --style gozero
    ```

    - 或者直接手动创建 handler、logic 文件（已做）✅
4. ✅ 手动配置数据库连接和模型初始化（已做）
5. ✅ 实现业务逻辑（已做）

---

## 五、注意事项 ⚠️

1. **goctl api go 会覆盖文件**：如果项目已有代码，使用前先备份或使用临时目录
2. **routes.go 是自动生成的**：如果手动修改了 routes.go，下次生成会被覆盖
3. **types.go 是自动生成的**：类型定义应该在 .api 文件中维护
4. **model 的 go.mod**：如果 model 目录有独立的 go.mod，需要删除以使用主模块

---

## 总结

### goctl 命令对比

|命令|用途|使用场景|位置|
| ----| ----------------| -----------------------| ----------------|
|`goctl api new <name>`|创建全新项目|**项目初始化第一步**（空目录）|最开始|
|`goctl api go`|从 .api 生成代码|已有项目，添加/更新 API|编写 .api 文件后|
|`goctl model mysql`|生成 Model|需要数据库模型时|任何时候|

**可以用 goctl 自动生成：**

- ✅ **项目结构**（`goctl api new` - 仅新项目）
- ✅ **Model 代码**（从 SQL 或数据库）
- ✅ **Handler、Logic、Types、Routes**（从 .api 文件）

**需要手动配置：**

- ⚙️ 数据库配置（config.go）
- ⚙️ 模型初始化（servicecontext.go）
- ⚙️ 配置文件（yaml）
- ⚙️ 业务逻辑实现（logic 文件）

### 快速参考

**新项目流程：**

```bash
goctl api new admin_role          # 1. 创建项目（第一步）
# 编辑 api/admin_role.api          # 2. 编写 API 定义
goctl api go -api api/... -dir .  # 3. 生成代码
goctl model mysql ...              # 4. 生成 Model
# 配置数据库和业务逻辑              # 5. 手动配置
```

**已有项目流程：**

```bash
# 已有 .api 文件
goctl api go -api api/... -dir .  # 生成代码（注意覆盖）
goctl model mysql ...              # 生成 Model（如果需要）
# 配置数据库和业务逻辑              # 手动配置
```
