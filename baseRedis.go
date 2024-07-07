package gmodel

import "github.com/redis/go-redis/v9"

type IRedis interface {
	RedisName() string                               //[实体类必须实现] 即对应的redis.key,由实体子类实现 - TODO 改进:加可变参数type,同一对象可存放不同的redis.key里
	RedisGet(child IRedis) (*redis.StringCmd, error) //[子基类默认实现]
	RedisSet(child IRedis) error                     //[子基类默认实现]
	RedisDel(child IRedis) error                     //[子基类默认实现]
}

// redis基类,实现最基本的共有方法,暂时无方法
type BaseRedis struct { /*这里不能写字段名*/
	//Base
}

//func (this *BaseRedis) Locks(keys ...string) []*lock.Locker {
//	return rds.ClientIns().LocksV2(keys...)
//}
//
//// 解锁 by cc 2019-6-28 11:43:32
//func (this *BaseRedis) Unlocks(locks []*lock.Locker) {
//	NewRedisDao().Unlocks(locks)
//}
