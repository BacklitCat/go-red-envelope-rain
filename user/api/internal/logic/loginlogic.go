package logic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go-red-envelope-rain/user/model"
	"strings"
	"time"

	"go-red-envelope-rain/user/api/internal/svc"
	"go-red-envelope-rain/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) getJwtToken(secretKey, userAccount string, iat, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["account"] = userAccount
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *LoginLogic) Login(req *types.Request) (resp *types.Response, err error) {
	// 校验参数
	if len(strings.TrimSpace(req.Account)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, errors.New("有参数为空")
	}

	// 访问数据库
	u, err := l.svcCtx.UserModel.FindOneByUserAccount(l.ctx, req.Account)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("用户不存在")
	default:
		return nil, err
	}

	// 校验密码（这个仅作演示用无加密）
	if u.UserPassword != req.Password {
		return nil, errors.New("账号或者密码错误")
	}

	// jwt
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessSecret := l.svcCtx.Config.Auth.AccessSecret
	jwtToken, err := l.getJwtToken(accessSecret, u.UserAccount, now, accessExpire)
	if err != nil {
		return nil, err
	}

	// response
	return &types.Response{
		UserAccount:  u.UserAccount,
		UserName:     u.UserName,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}
