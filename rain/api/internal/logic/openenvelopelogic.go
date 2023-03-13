package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-red-envelope-rain/rain/api/internal/svc"
	"go-red-envelope-rain/rain/api/internal/types"
	"go-red-envelope-rain/rain/model"
	"math/rand"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	maxMoney           int     = 1000 // 上限10元
	winningProbability float32 = 0.9  // 中奖概率90%
	balanceKey         string  = "rain:balance"
	balanceLock        string  = balanceKey
)

type OpenEnvelopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenEnvelopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenEnvelopeLogic {
	return &OpenEnvelopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenEnvelopeLogic) GetBalance() (int64, error) {
	balanceStr, _ := l.svcCtx.RedisClient.Get(balanceKey).Result()
	balanceInt, err := strconv.Atoi(balanceStr)
	if err != nil {
		return -1, err
	}
	return int64(balanceInt), nil
}

func (l *OpenEnvelopeLogic) SetBalance(balance int64) error {
	return l.svcCtx.RedisClient.Set(balanceKey, balance, 0).Err()
}

func (l *OpenEnvelopeLogic) OpenEnvelope() (resp *types.OpenResponse, err error) {
	account := l.ctx.Value("account").(string)
	db := l.svcCtx.RainModel
	one, err := db.FindOne(l.ctx, account)
	switch err {
	case nil:
		if !one.Status { // 检查封禁
			return nil, ErrForbidden
		}
		if one.Remaining <= 0 { // 检查剩余次数
			return nil, ErrNoRemaining
		}
		if time.Now().Sub(one.UpdateTime) < time.Second*3 { // 访问接口过快
			return nil, ErrWait
		}
	case model.ErrNotFound: // 找不到用户说明访问接口顺序错误
		return nil, ErrWrongSequence
	default:
		return nil, err
	}

	// 访问redis减金额库存
	randomMoney, err := l.updateRedis()
	if err != nil {
		return nil, err
	}

	// 开事务更新Mysql数据库
	err = l.updateDB(one, randomMoney)
	if err != nil {
		return nil, err
	}

	// 没有其他错误，返回
	return &types.OpenResponse{
		Account:   account,
		Amount:    int(randomMoney),
		Remaining: int(one.Remaining),
		Balance:   int(one.Balance),
	}, nil
}

func (l *OpenEnvelopeLogic) updateRedis() (int64, error) {
	// 生成随机金额
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMoney := int64(r.Intn(int(float32(maxMoney)/winningProbability)) - maxMoney)
	if randomMoney < 0 {
		return 0, nil // 未中奖（要占用一次机会，不能返回ErrUnWin）
	}

	balance, err := l.GetBalance()
	if err != nil || balance < randomMoney { // 错误设置金额或者系统余额不足
		return -1, ErrUnWin
	}

	mutex, err := NewListMutex(l.svcCtx.RedisClient, balanceLock, 0)
	if err != nil {
		return -1, ErrUnWin
	}

	mutex.Lock()
	defer mutex.Unlock()
	balance, err = l.GetBalance()
	if err != nil || balance < randomMoney { // 扣减前再次确认余额
		return -1, ErrUnWin
	}
	err = l.SetBalance(balance - randomMoney)
	if err != nil {
		return -1, ErrUnWin
	}

	return randomMoney, nil
}

func (l *OpenEnvelopeLogic) updateDB(u *model.Rain, num int64) error {
	// 开事务
	if err := l.svcCtx.RainModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		u.Remaining -= 1
		u.Balance += num
		err := l.svcCtx.RainModel.TransactUpdate(l.ctx, session, u)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
