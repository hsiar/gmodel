package gmodel

import (
	"fmt"
	"github.com/hsiar/gmodel/corm"
	"gorm.io/gorm"
)

type BaseDb struct {
	//c *corm.Client `gorm:"-"`
}

// 默认返回key=default的crom实例
//func (this *BaseDb) SetClient(c ...*corm.Client) {
//	if len(c) > 0 {
//		this.c = c[0]
//	} else {
//		this.c = corm.ClientIns()
//	}
//}

func (this *BaseDb) Client() *corm.Client { //默认返回主数据库
	return corm.ClientIns()
	//return this.c
}

//func (this *BaseDb) TableName(name string) string {
//	return this.Client().Config.TablePrefix + name
//}

func (this *BaseDb) Db() *gorm.DB {
	return this.Client().Db
}

func (this *BaseDb) Alias(child IDb, key ...string) *gorm.DB {
	alias := "t"
	if len(key) > 0 {
		alias = key[0]
	}
	return this.Db().Table(fmt.Sprintf("%s as %s", child.TableName(), alias))
}

func (this *BaseDb) Tx(child IDb, fn func(tx *gorm.DB) error) error {
	return this.Db().Transaction(fn)
}

//func (this BaseDb) Debug() {
//	this.Client().Db.Debug()
//}

// child必须为指针类型，下同
func (this *BaseDb) DbGet(child IDb, wheres ...any) (err error) {
	if len(wheres) == 0 {
		wheres = make([]any, 2)
		wheres = append(wheres, "id", child.GetId())
	}
	return this.Alias(child).First(child, wheres...).Error
}

func (this *BaseDb) DbFirst(child IDb, orderBy string, conds ...any) (err error) {
	return this.Alias(child).Order(orderBy).Take(child, conds...).Error
}

// list示例：*[]*User,*[]map[string]any todo 支持范型
func (this *BaseDb) DbPage(child IDb, list any, keys, join, order string, page, limit int, wheres ...any) (count int64, err error) {
	offset := (page - 1) * limit
	//return this.Db().Select(keys).Joins(join).Order(order).Limit(limit).Offset(offset).Find(list, conds...).Error
	//return this.Db().Model(child).Select(keys).Joins(join).Order(order).Limit(limit).Offset(offset).Find(list, conds...).Error
	//return this.Db().Table("z_user as t").Select(keys).Joins(join).Order(order).Limit(limit).Offset(offset).Find(list, conds...).Error
	//get count
	if count, err = this.DbCount(child, wheres...); err != nil {
		return
	}
	err = this.Alias(child).Select(keys).Joins(join).Order(order).Limit(limit).Offset(offset).Find(list, wheres...).Error
	return
}

func (this *BaseDb) DbAdd(child IDb) (err error) {
	return this.Client().Db.Create(child).Error
}

func (this *BaseDb) DbUpdate(child IDb, keys ...string) (err error) {
	if len(keys) > 0 {
		return this.Client().Db.Model(child).Select(keys).Updates(child).Error
	} else {
		return this.Client().Db.Model(child).Updates(child).Error
	}
}

func (this *BaseDb) DbSave(child IDb) (err error) {
	return this.Client().Db.Save(child).Error
}

func (this *BaseDb) DbDel(child IDb) (err error) {
	return this.Client().Db.Delete(child).Error
}

func (this *BaseDb) DbCount(child IDb, wheres ...any) (count int64, err error) {
	db := this.Alias(child).Model(child)
	for i := 0; i < len(wheres); i += 2 {
		db.Where(wheres[i], wheres[i+1])
	}
	if err = db.Count(&count).Error; err != nil {
		return 0, err
	}
	return
}

func (this *BaseDb) DbExist(child IDb, wheres ...any) bool {
	count, _ := this.DbCount(child, wheres...)
	return count > 0
}
