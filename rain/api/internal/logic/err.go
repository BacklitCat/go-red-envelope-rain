package logic

import "errors"

var (
	ErrForbidden     = errors.New("用户被该活动封禁")
	ErrWrongSequence = errors.New("用户在初始化活动前非法访问接口")
	ErrWait          = errors.New("用户访问接口过快")
	ErrNoRemaining   = errors.New("用户无剩余活动次数")
	ErrUnWin         = errors.New("用户因系统无余额/其他异常未中奖")
)
