package gmodel

// 同一模型可设置多个Redis-name-value的实现
// 实现类参考:PrizeRedPacketSeckill
type IRedisSet interface {
	IRedis
	//set专属
	SetRedisName(name string)
	GetRedisName() string
	SetRedisValue(value interface{}) //当存储的值为实体类对象本身json串时,必须由实体类实现,[解决插入数据库后,该值的id为0],参考seller_group
	GetRedisValue() interface{}
	SetRedisExpire(duration int64) //IRedisSetV2
	GetRedisExpire() int64

	//解决 当value为对象本身json串时,保存redis同时重新取对象的json串的办法
	//1.基类BaseRedisSetV2建redisType私有成员 +
	//2.接口IRedisSetV2增加SetRedisType,GetReidsType两个方法 +
	//3.保存redis时,即基类RedisSet方法里保存前重新调用实体类的SetRedisModel +
	//4.实体类重载SetRedisModel方法,参数value变为私有成员redisType值,由里面指定真正的redisValue值 +

	SetRedisType(typ int8) IRedisSet
	GetRedisType() int8

	SetRedisModel(_type int8) IRedisSet //设置单个redis模型类型,实体类必须实现的方法
	GetRedisModels() []int8             //获取所有redis模型类型,实体类必须实现的方法
}
