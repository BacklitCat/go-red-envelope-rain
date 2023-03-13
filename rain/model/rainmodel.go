package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RainModel = (*customRainModel)(nil)

type (
	// RainModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRainModel.
	RainModel interface {
		rainModel
		TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		TransactUpdate(ctx context.Context, session sqlx.Session, data *Rain) error
	}

	customRainModel struct {
		*defaultRainModel
	}
)

// NewRainModel returns a model for the database table.
func NewRainModel(conn sqlx.SqlConn, c cache.CacheConf) RainModel {
	return &customRainModel{
		defaultRainModel: newRainModel(conn, c),
	}
}

func (m *customRainModel) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.CachedConn.TransactCtx(ctx, fn)
}

func (m *customRainModel) TransactUpdate(ctx context.Context, session sqlx.Session, data *Rain) error {
	rainUserAccountKey := fmt.Sprintf("%s%v", cacheRainUserAccountPrefix, data.UserAccount)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, _ sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `user_account` = ?", m.table, rainRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.Status, data.Remaining, data.Balance, data.UserAccount)
	}, rainUserAccountKey)
	return err
}
