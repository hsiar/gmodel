package gmodel

import (
	"gorm.io/gorm"
)

// 支持同一模型多redis-name保存
type IDbRedisSet interface {
	IDb
	IRedisSet
	RedisFromDbCondSet(whereParams []interface{}, orderBy string)
	RedisFromDbCondGet() ([]interface{}, string)

	DbRedisGet(child IDbRedisSet, force2Redis bool) error
	DbRedisSave(child IDbRedisSet, isAdd ...bool) error
	DbRedisDel(child IDbRedisSet) error
}

// 自动过期的redis缓存数据模型 by cc 2020-3-31 01:26:39
// 保存数据类型为set类型
// 具体模型类要实现
type BaseDbRedisSet struct {
	Base
	//这里不能写其它字段名
	BaseDb
	BaseRedisSet

	//私有成员
	//redisFromDbCond RedisFromDbCond
	whereParams []interface{}
	orderBy     string
}

func (this *BaseDbRedisSet) RedisFromDbCondSet(whereParams []interface{}, orderBy string) {
	this.whereParams = whereParams
	this.orderBy = orderBy
}
func (this *BaseDbRedisSet) RedisFromDbCondGet() ([]interface{}, string) {
	return this.whereParams, this.orderBy
}

func (this *BaseDbRedisSet) Init(m *BaseRedisSetModel) {
	this.SetRedisName(m.Name)
	this.SetRedisValue(m.Value)
	this.SetRedisExpire(m.Expire)
	this.RedisFromDbCondSet(m.DbWhere, m.DbOrder)
}

// overload BaseRedisSetV2.SetRedisType for son struct linking operator by cc 2020-5-9 01:24:17
func (this *BaseDbRedisSet) SetRMT(child IDbRedisSet, typ int8) IDbRedisSet {
	child.SetRedisType(typ)
	return child
}

//redis数据从数据库获取的条件
//默认为[]interface{}{"id", child.GetId()}
//func (this *BaseRedis) RedisFromDbCond() []interface{} {
//	return nil
//}

// 从redis中取不到,从数据库取只能通过ID参数
// 依赖:this.id
func (this *BaseDbRedisSet) DbRedisGet(child IDbRedisSet, force2Redis bool) error {
	var (
		err     error
		cond    []interface{}
		orderBy string
		//cmd     *redis.StringCmd
	)
	//this.SetRedisModel()
	if _, err = this.RedisGet(child); err != nil { //出错从数据库取
		cond, orderBy = child.RedisFromDbCondGet()
		if err = this.DbFirst(child, orderBy, cond...); err != nil {
			return err
		} else if force2Redis { //强制保存到redis缓存
			_ = this.RedisSet(child)
		}
		return nil
	}
	return nil
}

// 同时插入/更新数据库和redis
// 依赖:child数据必须为全字段数据
// 保存数据类型为set类型,具体模型类要实现RedisName(),RedisSetValue()两个接口,具体参数user模型实现
func (this *BaseDbRedisSet) DbRedisSave(child IDbRedisSet, isAdd ...bool) (err error) {
	return child.Client().Db.Transaction(func(tx *gorm.DB) error {
		if child.GetId() == 0 || (len(isAdd) > 0 && isAdd[0]) {
			if err = tx.Create(child).Error; err != nil { //新增到数据库
				return err
			}
		} else {
			if err = tx.Model(child).Updates(child).Error; err != nil { //更新数据库
				return err
			}
		}
		if err = this.RedisSet(child); err != nil {
			return err
		}
		//返回 nil 提交事务
		return nil
	})
}

// 删除数据
func (this *BaseDbRedisSet) DbRedisDel(child IDbRedisSet) error {
	return child.Client().Db.Transaction(func(tx *gorm.DB) error {
		if err := this.RedisDel(child); err != nil {
			return err
		} else if err = tx.Delete(child).Error; err != nil {
			return err
		}
		return nil
	})
}

// type CheckOwnerFn func() bool
//func (this *BaseDbRedisSetV2) DbRedisOwnerCheck(child IDbRedisSetV2, checkFn func() bool) (err error) {
//	if err = child.DbRedisGet(child, true); err != nil {
//		return errors.New("not found")
//	} else if !checkFn() {
//		return errors.New("not owner")
//	}
//	return nil
//}
