package logic

import (
	"github.com/go-redis/redis"
	"time"
)

type ListMutex struct {
	db       *redis.Client
	LockPath string
	WaitTime time.Duration
}

func NewListMutex(db *redis.Client, listName string, waitTime time.Duration) (*ListMutex, error) {
	// 检查连接
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}
	if waitTime < 0 {
		waitTime = time.Duration(0) // NO TTL LIMIT
	}

	listExistPath, listPath := "ListMutex:EXIST:"+listName, "ListMutex:LIST:"+listName
	// 检查List是否存在，若不存在，则抢占创建
	created, err := db.SetNX(listExistPath, "true", 0).Result()
	if err != nil {
		return nil, err
	}
	if created {
		_, err = db.RPush(listPath, "lock").Result()
		if err != nil {
			return nil, err
		}
	}

	return &ListMutex{
		db:       db,
		LockPath: listPath,
		WaitTime: waitTime,
	}, nil
}

func (m *ListMutex) Lock() {
	_, err := m.db.BLPop(m.WaitTime, m.LockPath).Result()
	if err != nil {
		panic(err)
	}
}

func (m *ListMutex) Unlock() {
	_, err := m.db.RPush(m.LockPath, "lock").Result()
	if err != nil {
		panic(err)
	}
}
