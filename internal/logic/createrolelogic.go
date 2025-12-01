// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"test/model/admin_role"

	"test/internal/svc"
	"test/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.CreateRoleResponse, err error) {
	_, err = l.svcCtx.AdminRoleModel.Insert(l.ctx, &admin_role.AdminRole{
		RoleName: req.RoleName,
		Desc:     req.Desc,
		Status:   req.Status,
		IsSuper:  req.IsSuper,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateRoleResponse{
		Result: true,
	}, nil
}
