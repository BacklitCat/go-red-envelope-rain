package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-red-envelope-rain/rain/api/internal/svc"
	"go-red-envelope-rain/rain/api/internal/types"
	"go-red-envelope-rain/rain/model"
	"go-red-envelope-rain/user/rpc/types/user"
)

type CheckStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckStatusLogic {
	return &CheckStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckStatusLogic) CheckStatus() (resp *types.CheckResponse, err error) {
	jwtAccount := l.ctx.Value("account").(string)

	userInfo, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.AccountReq{
		Account: jwtAccount,
	})

	if err != nil {
		return nil, err
	}

	db := l.svcCtx.RainModel
	one, err := db.FindOne(l.ctx, userInfo.Account)
	switch err {
	case nil:
		// 已经参与过，首先判断是否在黑名单
		if !one.Status {
			return nil, ErrForbidden
		}
		return &types.CheckResponse{
			Account:   userInfo.Account,
			Remaining: int(one.Remaining),
			Balance:   int(one.Balance),
		}, nil
	case model.ErrNotFound:
		// 初次参加活动，初始化
		_, errInsert := db.Insert(l.ctx, &model.Rain{
			UserAccount: userInfo.Account,
			Status:      true,
			Remaining:   10,
			Balance:     0,
		})
		if errInsert != nil {
			return nil, errInsert
		}
		return &types.CheckResponse{
			Account:   userInfo.Account,
			Remaining: 10,
			Balance:   0,
		}, nil
	default:
		return nil, err
	}
}
