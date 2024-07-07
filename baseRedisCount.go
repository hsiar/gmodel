package gmodel

import (
	"context"
	"fmt"
	"github.com/hsiar/gmodel/rds"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"reflect"
	"time"
)

// 计数器,独立基类,供实体类继承使用
type BaseRedisCount struct {
}

// 参数child必须为指针类型
func (this *BaseRedisCount) redisCountName(child IBase, key string) string {
	return fmt.Sprintf("UCount.%s.%d.%s", reflect.ValueOf(child).Elem().Type().Name(), child.GetId(), key)
}

func (this *BaseRedisCount) getCmd(child IBase, key string) *redis.StringCmd {
	return rds.ClientIns().Get(context.Background(), this.redisCountName(child, key))
}

func (this *BaseRedisCount) GetRedisCountInt64(child IBase, key string) int64 {
	var (
		count int64
		err   error
	)
	if count, err = this.getCmd(child, key).Int64(); err != nil {
		return 0
	}
	return count
}

func (this *BaseRedisCount) GetRedisCountInt(child IBase, key string) int {
	var (
		count int
		err   error
	)
	if count, err = this.getCmd(child, key).Int(); err != nil {
		return 0
	}
	return count
}

// 默认一直存在
func (this *BaseRedisCount) SetRedisCount(child IBase, key string, num interface{}, expireTime ...int64) error {
	var (
		expire int64
		err    error
		ok     bool
	)
	if _, ok = num.(int); !ok {
		if _, ok = num.(int64); !ok {
			return errors.New("参数num必须为int或int64类型")
		}
	}
	if len(expireTime) > 0 {
		expire = expireTime[0]
	}
	_, err = rds.ClientIns().Set(context.Background(), this.redisCountName(child, key), num, time.Duration(expire)*time.Second).Result()
	return err
}

func (this *BaseRedisCount) IncrBy(child IBase, key string, num int) error {
	result := rds.ClientIns().IncrBy(context.Background(), this.redisCountName(child, key), int64(num))
	return result.Err()
}

func (this *BaseRedisCount) DecrBy(child IBase, key string, num int) error {
	result := rds.ClientIns().DecrBy(context.Background(), this.redisCountName(child, key), int64(num))
	return result.Err()
}
