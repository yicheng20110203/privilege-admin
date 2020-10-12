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

type _menuBackUrlModel struct {
    DefaultModel
    Id         int32     `json:"id"`
    MenuKey    string    `json:"menu_key"`
    BackUrl    string    `json:"back_url"`
    Desc       string    `json:"desc"`
    IsDelete   int32     `json:"is_delete"`
    CreateTime time.Time `json:"create_time"`
    UpdateTime time.Time `json:"update_time"`
}

var (
    MenuBackUrlModel *_menuBackUrlModel
)

func (*_menuBackUrlModel) TableName() string {
    return "privilege_menu_back_url"
}

func (u *_menuBackUrlModel) Find(in *pbs.PrivilegeMenuBackUrlFindReq) (out *pbs.PrivilegeMenuBackUrlFindResp, err error) {
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

    if where.GetBackUrl() != "" {
        if where.GetFuzzyBackUrl() {
            db = db.Where("back_url LIKE ?", "%"+where.GetBackUrl()+"%")
        } else {
            db = db.Where("back_url = ?", where.GetBackUrl())
        }
    }

    if where.GetDesc() != "" {
        db = db.Where("desc LIKE ?", "%"+where.GetDesc()+"%")
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

    objs := make([]*pbs.PrivilegeMenuBackUrl, 0)
    out = &pbs.PrivilegeMenuBackUrlFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }

    ms := make([]*_menuBackUrlModel, 0)
    err = db.Find(&ms).Error
    if err == gorm.ErrRecordNotFound {
        err = nil
    }

    if err != nil {
        logger.Logger.Errorf("_menuBackUrlModel.Find(%v) error: %v", *in, err)
        return nil, err
    }

    for _, v := range ms {
        objs = append(objs, &pbs.PrivilegeMenuBackUrl{
            Id:         v.Id,
            MenuKey:    v.MenuKey,
            BackUrl:    v.BackUrl,
            Desc:       v.Desc,
            IsDelete:   v.IsDelete,
            CreateTime: int32(v.CreateTime.Unix()),
            UpdateTime: int32(v.UpdateTime.Unix()),
        })
    }
    out = &pbs.PrivilegeMenuBackUrlFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }

    return
}

func (u *_menuBackUrlModel) Save(in *pbs.PrivilegeMenuBackUrlSaveReq) (out *pbs.PrivilegeMenuBackUrlSaveResp, err error) {
    //nowTime := int32(time.Now().Unix())
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    if in.GetId() > 0 {
        id := in.GetId()
        in.Id = 0
        //in.UpdateTime = nowTime
        err = db.Model(&_menuBackUrlModel{}).Where("id = ?", id).Update(in).Error
        if err != nil {
            logger.Logger.Errorf("_menuBackUrlModel.Save(%v) t1 error <%v>", *in, err)
            return
        }
        return &pbs.PrivilegeMenuBackUrlSaveResp{Id: id}, nil
    }

    if in.GetBackUrl() == "" {
        logger.Logger.Errorf("_menuBackUrlModel.Save back_url error: (in = %+v)", in)
        err = errors.New("back_url can not be null")
        return
    }

    obj := &_menuBackUrlModel{
        Id:         in.GetId(),
        MenuKey:    in.GetMenuKey(),
        BackUrl:    in.GetBackUrl(),
        Desc:       in.GetDesc(),
        IsDelete:   global.IsDeleteNormal,
        CreateTime: time.Now(),
        UpdateTime: time.Now(),
    }
    err = db.Create(obj).Error
    if err != nil {
        logger.Logger.Errorf("_menuBackUrlModel.Save(%#v) insert t4 error: <%+v>", obj, err)
        return
    }
    out = &pbs.PrivilegeMenuBackUrlSaveResp{
        Id: obj.Id,
    }
    return
}

func (u *_menuBackUrlModel) Delete(in *pbs.PrivilegeMenuBackUrlDeleteReq) (out *pbs.PrivilegeMenuBackUrlDeleteResp, err error) {
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    where := in

    if len(where.GetId()) > 0 {
        db = db.Where("id IN (?)", where.GetId())
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

    if where.GetBackUrl() != "" {
        if where.GetFuzzyBackUrl() {
            db = db.Where("back_url LIKE ?", "%"+where.GetBackUrl()+"%")
        } else {
            db = db.Where("back_url = ?", where.GetBackUrl())
        }
    }

    if where.GetDesc() != "" {
        db = db.Where("desc LIKE ?", "%"+where.GetDesc()+"%")
    }

    if len(where.GetIsDelete()) > 0 {
        db = db.Where("is_delete IN (?)", where.GetIsDelete())
    }

    var total int32
    db.Count(&total)
    out = &pbs.PrivilegeMenuBackUrlDeleteResp{
        AffectRows: total,
    }
    if total == 0 {
        return
    }

    err = db.Delete(_menuBackUrlModel{}).Error
    if err != nil {
        logger.Logger.Error("_menuBackUrlModel.Delete(%+v) error:%v", *in, err)
        return nil, err
    }

    return
}
