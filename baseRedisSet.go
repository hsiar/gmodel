package gmodel

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/hsiar/gmodel/rds"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// 专用用于实体类初始化redis模型用的
type BaseRedisSetModel struct {
	Name   string
	Value  any
	Expire int64

	DbWhere []any
	DbOrder string
}

type BaseRedisSet struct {
	Base
	//已实现 RedisSet RedisGet RedisDel RedisSetValue RedisSetExpire[默认无过期] RedisName
	//需实现 SetRedisModel GetRedisModels

	//私有成员
	redisName   string
	redisValue  interface{}
	redisExpire int64
	redisType   int8 //当前对象的类型,比如:shop.id.1和shop.user_id.1
}

// 默认为空,实体类也可以实现[参考:PrizeRedPacketSeckill类],也可以不实现
func (this *BaseRedisSet) RedisName() string {
	return ""
}

func (this *BaseRedisSet) SetRedisName(name string) {
	this.redisName = name
}
func (this *BaseRedisSet) GetRedisName() string {
	return this.redisName
}

func (this *BaseRedisSet) SetRedisValue(value interface{}) {
	this.redisValue = value
}
func (this *BaseRedisSet) GetRedisValue() interface{} {
	return this.redisValue
}

func (this *BaseRedisSet) SetRedisExpire(duration int64) {
	this.redisExpire = duration
}
func (this *BaseRedisSet) GetRedisExpire() int64 {
	return this.redisExpire
}

func (this *BaseRedisSet) SetRedisType(typ int8) IRedisSet {
	this.redisType = typ
	return this
}
func (this *BaseRedisSet) GetRedisType() int8 {
	return this.redisType
}

func (this *BaseRedisSet) SetRedisModel(typ int8) IRedisSet {
	panic("请在实体类里实现SetRedisModel方法")
	return nil
}

// 默认只有一种redis模型,实体类可以不实现
func (this *BaseRedisSet) GetRedisModels() []int8 {
	//panic("请在实体类里实现GetRedisModels方法")
	return []int8{1} //默认只有一种redis模型
}

// 通过BaseRedisSetV2Model初始化私有成员变量 by cc 2020-5-7 04:01:32
func (this *BaseRedisSet) Init(m *BaseRedisSetModel) {
	this.SetRedisName(m.Name)
	this.SetRedisValue(m.Value)
	this.SetRedisExpire(m.Expire)
}

func (this *BaseRedisSet) check(child IRedis) (IRedisSet, error) {
	var (
		ok     bool
		child2 IRedisSet
	)
	if child2, ok = child.(IRedisSet); !ok {
		return nil, errors.New("当前对象不是IRedisSetV2类型")
	}
	return child2, nil
}

func (this *BaseRedisSet) RedisExist(child IRedisSet) bool {
	child2, err := this.check(child)
	if err != nil {
		return false
	}
	child2.SetRedisModel(child2.GetRedisType())
	//name := child.GetRedisName()
	//logs.Debug(name)
	return rds.ClientIns().Exists(context.Background(), child.GetRedisName()).Val() == 1
}

// 以下实现IRedis接口
func (this *BaseRedisSet) RedisSet(child IRedis) error {
	var (
		child2 IRedisSet
		err    error
	)
	if child2, err = this.check(child); err != nil {
		return err
	}
	for _, v := range child2.GetRedisModels() { //同时保存所有redis模型
		child2.SetRedisType(v)
		child2.SetRedisModel(child2.GetRedisType())
		_, err = rds.ClientIns().Set(context.Background(), child2.GetRedisName(), child2.GetRedisValue(), time.Duration(child2.GetRedisExpire())*time.Second).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// 该类返回的err都为nil
// 默认实现为获取对象
func (this *BaseRedisSet) RedisGet(child IRedis) (*redis.StringCmd, error) {
	var (
		child2 IRedisSet
		err    error
	)
	if child2, err = this.check(child); err != nil {
		return nil, err
	}
	child2.SetRedisModel(child2.GetRedisType())
	stringCmd := rds.ClientIns().Get(context.Background(), child2.GetRedisName())
	if jsonstr, err := stringCmd.Result(); err != nil {
		return nil, err
	} else if govalidator.IsJSON(jsonstr) && strings.Contains(jsonstr, "{") {
		var obj = make(map[string]interface{})
		if err = jsoniter.UnmarshalFromString(jsonstr, &obj); err != nil {
			return stringCmd, nil
		}
		return nil, jsoniter.UnmarshalFromString(jsonstr, child)
	}
	return stringCmd, nil
}

func (this *BaseRedisSet) RedisDel(child IRedis) (err error) {
	var (
		child2 IRedisSet
	)
	if child2, err = this.check(child); err != nil {
		return err
	}

	for _, v := range child2.GetRedisModels() { //同时保存所有redis模型
		child2.SetRedisType(v)
		child2.SetRedisModel(child2.GetRedisType())
		_, err = rds.ClientIns().Del(context.Background(), child2.GetRedisName()).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
