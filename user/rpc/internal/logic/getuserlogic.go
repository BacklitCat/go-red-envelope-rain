package logic

import (
	"context"

	"go-red-envelope-rain/user/rpc/internal/svc"
	"go-red-envelope-rain/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.AccountReq) (*user.UserInfoReply, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserAccount(l.ctx, in.Account)
	if err != nil {
		return nil, err
	}

	return &user.UserInfoReply{
		Id:      u.Id,
		Account: u.UserAccount,
		Name:    u.UserName,
	}, nil
}
