// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"test/internal/config"
	"test/model/admin_role"  // 导入 model
	"test/model/admin_users" // 导入 model

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	AdminUsersModel admin_users.AdminUsersModel // 添加 AdminUsersModel
	AdminRoleModel  admin_role.AdminRoleModel   // 添加 AdminRoleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:          c,
		AdminUsersModel: admin_users.NewAdminUsersModel(sqlConn, c.Cache), // 初始化 AdminUsersModel
		AdminRoleModel:  admin_role.NewAdminRoleModel(sqlConn, c.Cache),   // 初始化 AdminRoleModel
	}
}
