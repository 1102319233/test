// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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
