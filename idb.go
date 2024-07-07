package gmodel

import "github.com/hsiar/gmodel/corm"

type IDb interface {
	IBase

	//基类实现
	SetClient(c *corm.Client)
	Client() *corm.Client

	//实体类实现

	TableName() string

	//Dao() *Dao             //[基类默认实现]
	//IsHash() bool          //[基类默认实现] 是否分表
	//HashTableName() string //[基类默认实现] 该方法仅供hash分表使用,未分表的模型,默认返回TableName()

	//DbPageAllFld(child IDb, order string, page, limit int, params ...interface{}) ([]orm.Params, int, error)

}
