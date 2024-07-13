package gmodel

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hsiar/gbase"
	jsoniter "github.com/json-iterator/go"
)

type Base struct { /*这里不能写字段名*/
}

func (this *Base) ToString(child IBase) string {
	str, _ := jsoniter.MarshalToString(child)
	return str
}

func (this *Base) ToCMap(child IBase) (gbase.Map, error) {
	m := gbase.Map{}
	if bytes, err := jsoniter.Marshal(child); err != nil {
		return nil, err
	} else if err = jsoniter.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ToCMapExclude
//
//keys:xx,xx,xx
func (this *Base) ToCMapExclude(child IBase, keys ...string) (cm gbase.Map, err error) {
	if cm, err = this.ToCMap(child); err != nil {
		return
	}
	cm.RemoveKeys(keys...)
	return

}

func (this *Base) ToCMapMustExclude(child IBase, keys ...string) (cm gbase.Map) {
	var err error
	if cm, err = this.ToCMap(child); err != nil {
		hlog.Errorf("Base.ToCMapMustExclude err:%s", err.Error())
		panic("Base.ToCMapMustExclude err")
	}
	cm.RemoveKeys(keys...)
	return

}

func (this *Base) ToCMapInclude(child IBase, keys ...string) (cm gbase.Map, err error) {
	var (
		all gbase.Map
	)
	cm = gbase.Map{}
	if all, err = this.ToCMap(child); err != nil {
		return
	}
	for _, v := range keys {
		if all.Exist(v) {
			cm[v] = all[v]
		}
	}
	return
}

func (this *Base) ToCMapMustInclude(child IBase, keys ...string) (cm gbase.Map) {
	var (
		all gbase.Map
		err error
	)
	cm = gbase.Map{}
	if all, err = this.ToCMap(child); err != nil {
		panic("Base.ToCMapMustInclude err:" + err.Error())
	}
	for _, v := range keys {
		if all.Exist(v) {
			cm[v] = all[v]
		} else {
			panic(fmt.Sprintf("Base.ToCMapMustInclude not found key:%s in data", v))
		}
	}
	return
}

// map转struct by cc 2020-3-11 10:09:04
func (this *Base) FromMap(child IBase, params *orm.Params) error {
	var (
		jsonBytes []byte
		err       error
	)
	if jsonBytes, err = jsoniter.Marshal(params); err != nil {
		return err
	} else if err = jsoniter.Unmarshal(jsonBytes, child); err != nil {
		return err
	}
	return nil
}

// from map/struct
func (this *Base) FromX(child IBase, params interface{}) error {
	var (
		jsonBytes []byte
		err       error
	)
	if jsonBytes, err = jsoniter.Marshal(params); err != nil {
		return err
	} else if err = jsoniter.Unmarshal(jsonBytes, child); err != nil {
		return err
	}
	return nil
}

// string转到child对象
func (this *Base) FromString(child IBase, jsonstr string) error {
	return jsoniter.UnmarshalFromString(jsonstr, child)
}
