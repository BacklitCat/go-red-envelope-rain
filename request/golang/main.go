package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
)

var (
	host            = "192.168.2.176"
	loginUrl        = "http://" + host + ":7050/user/login"
	rainUrl         = "http://" + host + ":7054/rain"
	checkStatusUrl  = rainUrl + "/checkStatus"
	openEnvelopeUrl = rainUrl + "/openEnvelope"
	begin           = flag.Int("begin", 1000001, "user id begin")
	num             = flag.Int("num", 2000, "user num")
)

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(*num)
	for i := *begin; i < *begin+*num; i++ {
		u := NewUser(strconv.Itoa(i))
		go func() {
			u.work()
			wg.Done()
		}()
	}
	wg.Wait()

	// mysql
	db, err := sql.Open("mysql", "red_envelope_rain:123456@tcp(127.0.0.1:3306)/red_envelope_rain?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	row, err := db.Query(fmt.Sprintf("select sum(balance) from rain where user_account between %d and %d", *begin, *begin+*num))
	defer row.Close()
	var userBalance int
	row.Next()
	err = row.Scan(&userBalance)
	if err != nil {
		fmt.Println(err)
		return
	}

	// redis
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	sysRemainingStr, _ := redisClient.Get("rain:balance").Result()
	sysRemaining, err := strconv.Atoi(sysRemainingStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fmt.Sprintf("=========================\n"+
		"模拟完成，共有%d位用户参与模拟，共获得：%d，系统余额：%d，系统开支：%d",
		*num, userBalance, sysRemaining, 10000000000-sysRemaining))
}
