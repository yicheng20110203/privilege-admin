package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.ceibsmoment.com/c/mp/client"
	"gitlab.ceibsmoment.com/c/mp/global"
	"gitlab.ceibsmoment.com/c/mp/logger"
	pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
	"math"
	"time"
)

type _roleMenuModel struct {
	DefaultModel
	Id         int32     `json:"id"`
	RoleKey    string    `json:"role_key"`
	MenuKey    string    `json:"menu_key"`
	IsDelete   int32     `json:"is_delete"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

var (
	RoleMenuModel *_roleMenuModel
)

func (*_roleMenuModel) TableName() string {
	return "privilege_role_menu"
}

func (u *_roleMenuModel) Find(in *pbs.PrivilegeRoleMenuFindReq) (out *pbs.PrivilegeRoleMenuFindResp, err error) {
	db := client.GetMysqlDb()
	db = db.Table(u.TableName())

	where := in.GetWhere()
	if in.Page <= 0 {
		in.Page = GetDefaultPage()
	}
	if in.Size <= 0 {
		in.Size = GetDefaultSize()
	}

	if len(where.GetId()) > 0 {
		db = db.Where("id IN (?)", where.GetId())
	}

	if where.GetRoleKey() != "" {
		if where.GetFuzzyRoleKey() {
			db = db.Where("role_key LIKE ?", "%"+where.GetRoleKey()+"%")
		} else {
			db = db.Where("role_key = ?", where.GetRoleKey())
		}
	}

	if where.GetMenuKey() != "" {
		if where.GetFuzzyMenuKey() {
			db = db.Where("menu_key LIKE ?", "%"+where.GetMenuKey()+"%")
		} else {
			db = db.Where("menu_key = ?", where.GetMenuKey())
		}
	}

	if len(where.GetMenuKeys()) > 0 {
		db = db.Where("menu_key IN (?)", where.GetMenuKeys())
	}

	if len(where.GetRoleKey()) > 0 {
		db = db.Where("role_key IN (?)", where.GetRoleKey())
	}

	if len(where.GetIsDelete()) > 0 {
		db = db.Where("is_delete IN (?)", where.GetIsDelete())
	}

	var total int32
	db.Count(&total)

	if in.GetOrderBy() != "" {
		db = db.Order(in.GetOrderBy())
	} else {
		db = db.Order("id asc")
	}

	offset := (in.GetPage() - 1) * (in.GetSize())
	db = db.Limit(in.GetSize()).Offset(offset)

	objs := make([]*pbs.PrivilegeRoleMenu, 0)
	out = &pbs.PrivilegeRoleMenuFindResp{
		Page:      in.GetPage(),
		Size:      in.GetSize(),
		TotalSize: int32(total),
		TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
		List:      objs,
	}

	ms := make([]*_roleMenuModel, 0)
	err = db.Find(&ms).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	if err != nil {
		logger.Logger.Errorf("_roleMenuModel.Find(%v) error: %v", *in, err)
		return nil, err
	}

	for _, v := range ms {
		objs = append(objs, &pbs.PrivilegeRoleMenu{
			Id:         v.Id,
			RoleKey:    v.RoleKey,
			MenuKey:    v.MenuKey,
			IsDelete:   v.IsDelete,
			CreateTime: int32(v.CreateTime.Unix()),
			UpdateTime: int32(v.UpdateTime.Unix()),
		})
	}
	out = &pbs.PrivilegeRoleMenuFindResp{
		Page:      in.GetPage(),
		Size:      in.GetSize(),
		TotalSize: int32(total),
		TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
		List:      objs,
	}

	return
}

func (u *_roleMenuModel) Save(in *pbs.PrivilegeRoleMenuSaveReq) (out *pbs.PrivilegeRoleMenuSaveResp, err error) {
	//nowTime := int32(time.Now().Unix())
	db := client.GetMysqlDb()
	db = db.Table(u.TableName())

	if in.GetId() > 0 {
		id := in.GetId()
		in.Id = 0
		//in.UpdateTime = nowTime
		err = db.Model(&_roleMenuModel{}).Where("id = ?", id).Update(in).Error
		if err != nil {
			logger.Logger.Errorf("_roleMenuModel.Save(%v) t1 error <%v>", *in, err)
			return
		}
		return &pbs.PrivilegeRoleMenuSaveResp{Id: id}, nil
	}

	if in.GetMenuKey() == "" {
		logger.Logger.Errorf("_roleMenuModel.Save menu_key nil error: (in = %+v)", in)
		err = errors.New("menu_key can not be null")
		return
	}

	if in.GetRoleKey() == "" {
		logger.Logger.Errorf("_roleMenuModel.Save role_key nil error: (in = %+v)", in)
		err = errors.New("role_key can not be null")
		return
	}

	obj := &_roleMenuModel{
		Id:         in.GetId(),
		RoleKey:    in.GetRoleKey(),
		MenuKey:    in.GetMenuKey(),
		IsDelete:   global.IsDeleteNormal,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err = db.Create(obj).Error
	if err != nil {
		logger.Logger.Errorf("_roleMenuModel.Save(%#v) insert t4 error: <%+v>", obj, err)
		return
	}
	out = &pbs.PrivilegeRoleMenuSaveResp{
		Id: obj.Id,
	}
	return
}

func (u *_roleMenuModel) Delete(in *pbs.PrivilegeRoleMenuDeleteReq) (out *pbs.PrivilegeRoleMenuDeleteResp, err error) {
	db := client.GetMysqlDb()
	db = db.Table(u.TableName())

	where := in
	if len(where.GetId()) > 0 {
		db = db.Where("id IN (?)", where.GetId())
	}

	if where.GetRoleKey() != "" {
		if where.GetFuzzyRoleKey() {
			db = db.Where("role_key LIKE ?", "%"+where.GetRoleKey()+"%")
		} else {
			db = db.Where("role_key = ?", where.GetRoleKey())
		}
	}

	if where.GetMenuKey() != "" {
		if where.GetFuzzyMenuKey() {
			db = db.Where("menu_key LIKE ?", "%"+where.GetMenuKey()+"%")
		} else {
			db = db.Where("menu_key = ?", where.GetMenuKey())
		}
	}

	if len(where.GetMenuKeys()) > 0 {
		db = db.Where("menu_key IN (?)", where.GetMenuKeys())
	}

	if len(where.GetRoleKey()) > 0 {
		db = db.Where("role_key IN (?)", where.GetRoleKey())
	}

	if len(where.GetIsDelete()) > 0 {
		db = db.Where("is_delete IN (?)", where.GetIsDelete())
	}

	var total int32
	db.Count(&total)
	out = &pbs.PrivilegeRoleMenuDeleteResp{
		AffectRows: total,
	}
	if total == 0 {
		return
	}

	err = db.Delete(_roleMenuModel{}).Error
	if err != nil {
		logger.Logger.Error("_roleMenuModel.Delete(%+v) error:%v", *in, err)
		return nil, err
	}

	return
}
