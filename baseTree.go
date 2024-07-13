package gmodel

import (
	"fmt"
	"strings"
)

type ITree interface {
	IDb

	GetPid() int64
	GetPath() string
	DSonNum(tree ITree) int64
}

type BaseTree struct {
	BaseDbRedisSet

	//Children []ITree `orm:"-" json:"children"`
}

func (this *BaseTree) IsRoot(child ITree) bool {
	return child.GetPid() == 0
}

func (this *BaseTree) DSonNum(child ITree) int64 {
	count, _ := this.DbCount(child, "pid", child.GetId())
	return count
}

func (this *BaseTree) MkSonPath(child ITree) string {
	if child.GetId() == 0 {
		return ","
	} else {
		return fmt.Sprintf("%s%d,", child.GetPath(), child.GetId())
	}
}

func (this *BaseTree) HasDSon(child ITree, dson ITree) bool {
	return child.GetId() == dson.GetPid()
}

//func (this *BaseTree) HasGSon(child ITree,gsonOrId interface{}) bool {
//	var (
//		id   int64
//		ok    bool
//		err   error
//	)
//	if id, ok = gsonOrId.(int64); ok {
//
//		admin = &Admin{Id: uid}
//
//		if err = admin.DbRedisGet(admin, false); err != nil {
//			return false
//		}
//		return this.Id == admin.Pid
//	} else if admin, ok = uidOrAdmin.(*Admin); ok {
//		return this.Id == admin.Pid
//	}
//	return false
//}

// depend this.id
func (this *BaseTree) HasSon(child ITree, son ITree) bool {
	if has := this.HasDSon(child, son); has {
		return true
	} else if strings.Index(son.GetPath(), fmt.Sprintf(",%d,", child.GetId())) != -1 {
		return true
	}
	return false
}

//// get all son ids,contain dson and gson
//func (this *BaseTree) GetSonIds(child ITree) ([]int64, error) {
//	var (
//		sonIds []int64
//	)
//	sonIds = make([]int64, 0)
//	err := this.DbPage(child, "t.id", "", "", 1, -1, "t.path__contains", fmt.Sprintf("%%,%d,%%", child.GetId()))
//	if err != nil {
//		return nil, err
//	}
//	for _, v := range list {
//		id, _ := util.StringToInt64(v["id"].(string))
//		sonIds = append(sonIds, id)
//	}
//	return sonIds, nil
//}
//
//func (this *BaseTree) GetDSonIds(child ITree) ([]int64, error) {
//	var (
//		sonIds []int64
//	)
//	sonIds = make([]int64, 0)
//	list, _, err := this.DbPage(child, "t.id", "", "", 1, -1, "t.pid", child.GetId())
//	if err != nil {
//		return nil, err
//	}
//	for _, v := range list {
//		id, _ := util.StringToInt64(v["id"].(string))
//		sonIds = append(sonIds, id)
//	}
//	return sonIds, nil
//}
//
//func (this *BaseTree) GetParent(child ITree, parent ITree) (err error) {
//	if result, err := this.DbGet(parent, "t.id", child.GetPid()); err != nil {
//		return err
//	} else {
//		parent = result.(ITree)
//	}
//	return
//}
//
//func (this *BaseTree) GetRoot(child ITree, root ITree) error {
//	pathArr, err := util.ArraySplit2Int64(util.StrRmHeadEnd(child.GetPath()), ",")
//	if err != nil {
//		return err
//	}
//	if len(pathArr) == 0 {
//		root = child
//	} else if result, err := this.DbGet(root, "t.id", pathArr[0]); err != nil {
//		return err
//	} else {
//		root = result.(ITree)
//	}
//	return nil
//}

// fmt:,1,2,3, 包含自己的id
func (this *BaseTree) GetWholePath(child ITree) string {
	return fmt.Sprintf("%s%d,", child.GetPath(), child.GetId())
}
