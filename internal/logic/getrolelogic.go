// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"test/internal/svc"
	"test/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(req *types.GetRoleRequest) (resp *types.GetRoleResponse, err error) {
	// todo: add your logic here and delete this line
	role, err := l.svcCtx.AdminRoleModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	return &types.GetRoleResponse{
		RoleName: role.RoleName,
		Desc:     role.Desc,
		IsSuper:  role.IsSuper,
		Status:   role.Status,
	}, nil
}
