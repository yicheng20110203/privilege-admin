package model

import (
    "errors"
    "fmt"
    "github.com/jinzhu/gorm"
    "gitlab.ceibsmoment.com/c/mp/client"
    "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "math"
    "strconv"
    "time"
)

type _menuModel struct {
    DefaultModel
    Id           int32     `json:"id"`
    Path         string    `json:"path"`
    Component    string    `json:"component"`
    Title        string    `json:"title"`
    Name         string    `json:"name"`
    Icon         string    `json:"icon"`
    MenuKey      string    `json:"menu_key"`
    Level        int32     `json:"level"`
    DisplayOrder int32     `json:"display_order"`
    IsHidden     int32     `json:"is_hidden"`
    IsDelete     int32     `json:"is_delete"`
    CreateTime   time.Time `json:"create_time"`
    UpdateTime   time.Time `json:"update_time"`
}

var (
    MenuModel *_menuModel
)

func (*_menuModel) TableName() string {
    return "privilege_menu"
}

func (u *_menuModel) Find(in *pbs.PrivilegeMenuFindReq) (out *pbs.PrivilegeMenuFindResp, err error) {
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

    if where.GetPath() != "" {
        db = db.Where("path = ?", where.GetPath())
    }

    if len(where.GetComponent()) > 0 {
        db = db.Where("component IN (?)", where.GetComponent())
    }

    if len(where.GetIsHidden()) > 0 {
        db = db.Where("is_hidden IN(?)", where.GetIsHidden())
    }

    if where.GetTitle() != "" {
        db = db.Where("title = ?", where.GetTitle())
    }

    if where.GetName() != "" {
        db = db.Where("name = ?", where.GetName())
    }

    if where.GetIcon() != "" {
        db = db.Where("icon = ?", where.GetIcon())
    }

    if where.GetMenuKey() != "" {
        if where.GetFuzzyMenuKey() {
            db = db.Where("menu_key LIKE ?", where.GetMenuKey()+"%")
        } else {
            db = db.Where("menu_key = ?", where.GetMenuKey())
        }
    }

    if len(where.GetMenuKeys()) > 0 {
        db = db.Where("menu_key IN (?)", where.GetMenuKeys())
    }

    if len(where.GetLevel()) > 0 {
        db = db.Where("level IN (?)", where.GetLevel())
    }

    if len(where.GetIsDelete()) > 0 {
        db = db.Where("is_delete IN (?)", where.GetIsDelete())
    }

    var total int32
    db.Count(&total)

    if in.GetOrderBy() != "" {
        db = db.Order(in.GetOrderBy())
    } else {
        db = db.Order("display_order asc")
    }

    offset := (in.GetPage() - 1) * (in.GetSize())
    db = db.Limit(in.GetSize()).Offset(offset)

    objs := make([]*pbs.PrivilegeMenu, 0)
    out = &pbs.PrivilegeMenuFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }
    ms := make([]*_menuModel, 0)
    err = db.Find(&ms).Error
    if err == gorm.ErrRecordNotFound {
        err = nil
    }

    if err != nil {
        logger.Logger.Errorf("_menuModel.Find(%v) error: %v", *in, err)
        return nil, err
    }
    for _, v := range ms {
        objs = append(objs, &pbs.PrivilegeMenu{
            Id:           v.Id,
            Path:         v.Path,
            Component:    v.Component,
            Title:        v.Title,
            Name:         v.Name,
            Icon:         v.Icon,
            MenuKey:      v.MenuKey,
            Level:        v.Level,
            DisplayOrder: v.DisplayOrder,
            IsHidden:     v.IsHidden,
            IsDelete:     v.IsDelete,
            CreateTime:   int32(v.CreateTime.Unix()),
            UpdateTime:   int32(v.UpdateTime.Unix()),
        })
    }

    out = &pbs.PrivilegeMenuFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }

    return
}

func (u *_menuModel) Save(in *pbs.PrivilegeMenuSaveReq) (out *pbs.PrivilegeMenuSaveResp, err error) {
    //nowTime := int32(time.Now().Unix())
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    if in.GetId() > 0 {
        id := in.GetId()
        in.Id = 0

        var (
            m1 = &_menuModel{}
            m2 = &_menuModel{}
        )
        err = db.Model(&_menuModel{}).Where("id = ?", in.Id).First(&m1).Error
        if err == gorm.ErrRecordNotFound {
            err = nil
            return &pbs.PrivilegeMenuSaveResp{}, err
        }
        if err != nil {
            logger.Logger.Errorf("_menuModel.Save(%v) find same one error: %v", err)
            return nil, err
        }

        if in.Name != "" {
            err = db.Model(&_menuModel{}).Where("LENGTH(`menu_key`) = ? AND `id` NOT IN(?) AND `name` = ?", len(m1.MenuKey), id, in.Name).First(&m2).Error
            if err != gorm.ErrRecordNotFound && err != nil {
                logger.Logger.Errorf("_menuModel.Save(%v) find same two error: %v", err)
                return nil, err
            }

            if err == nil && m2 != nil {
                err = errors.New("同级目录下已存在菜单<" + in.Name + ">")
                logger.Logger.Errorf("_menuModel.Save() 同级目录已存在菜单: %s", in.Name)
                return nil, err
            }
        }

        err = db.Model(&_menuModel{}).Where("id = ?", id).Update(in).Error
        if err != nil {
            logger.Logger.Errorf("_menuModel.Save(%v) t1 error <%v>", *in, err)
            return
        }
        return &pbs.PrivilegeMenuSaveResp{Id: id}, nil
    }

    if in.GetPath() == "" || in.GetComponent() == "" || in.GetTitle() == "" || in.GetName() == "" {
        logger.Logger.Errorf("_menuModel.Save [path,component,title,name] empty error: (in = %+v)", in)
        err = errors.New("[path,component,title,name] can not be null")
        return
    }

    menuKey, _ := u.GetMenuKey(in.GetMenuKey())
    var search = &_menuModel{}
    err = db.Model(&_menuModel{}).Where("LENGTH(`menu_key`) = ? AND `name` = ?", len(menuKey), in.Name).First(&search).Error
    if err != gorm.ErrRecordNotFound && err != nil {
        logger.Logger.Errorf("_menuModel.Save insert find same name record error: (in = %+v) error: %v", *in, err)
        return nil, err
    }
    if err == nil && search != nil && search.Id > 0 {
        err = errors.New(fmt.Sprintf("已存在同名菜单<%s>", in.Name))
        return nil, err
    }
    obj := &_menuModel{
        Id:           in.GetId(),
        Path:         in.GetPath(),
        Component:    in.GetComponent(),
        Title:        in.GetTitle(),
        Name:         in.GetName(),
        Icon:         in.GetIcon(),
        MenuKey:      menuKey,
        Level:        int32(len(menuKey) / 3),
        DisplayOrder: in.GetDisplayOrder(),
        IsHidden:     in.GetIsHidden(),
        IsDelete:     global.IsDeleteNormal,
        CreateTime:   time.Now(),
        UpdateTime:   time.Now(),
    }
    err = db.Create(obj).Error
    if err != nil {
        logger.Logger.Errorf("_menuModel.Save(%#v) insert t4 error: <%+v>", obj, err)
        return
    }
    out = &pbs.PrivilegeMenuSaveResp{
        Id: obj.Id,
    }
    return
}

// 计算menu_key
func (u *_menuModel) GetMenuKey(menuKey string) (key string, err error) {
    m := &_menuModel{}
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    if menuKey == "" {
        err = db.Where("LENGTH(`menu_key`) = 3").Order("menu_key DESC").First(m).Error

        if err != nil && err != gorm.ErrRecordNotFound {
            logger.Logger.Errorf("_menuModel.GetMenuKey(%v) error: %v", menuKey, err)
            return "", err
        }

        if err == gorm.ErrRecordNotFound {
            key = "100"
            return key, nil
        }

        k := m.MenuKey
        o, _ := strconv.Atoi(k)
        o = o + 1
        k = fmt.Sprintf("%d", o)
        if len(k)%3 != 0 {
            for i := 0; i < len(k)%3; i++ {
                k = "0" + k
            }
        }
        key = k
        return key, err
    }

    err = db.Where("menu_key LIKE ?", menuKey+"%").Where("LENGTH(menu_key) = ?", len(menuKey)+3).Order("menu_key DESC").First(m).Error
    if err != nil && err != gorm.ErrRecordNotFound {
        logger.Logger.Errorf("_menuModel.GetMenuKey(%v) error: %v", menuKey, err)
        return "", err
    }
    if err == gorm.ErrRecordNotFound {
        key = menuKey + "001"
        return key, nil
    }

    k := m.MenuKey
    o, _ := strconv.Atoi(k[len(k)-3:])
    o = o + 1
    k = fmt.Sprintf("%d", o)
    if len(k)%3 != 0 {
        for i := 0; i < len(k)%3; i++ {
            k = "0" + k
        }
    }
    key = menuKey + k
    return
}

func (u *_menuModel) Delete(in *pbs.PrivilegeMenuDeleteReq) (out *pbs.PrivilegeMenuDeleteResp, err error) {
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    where := in

    if len(where.GetId()) > 0 {
        db = db.Where("id IN (?)", where.GetId())
    }

    if where.GetPath() != "" {
        db = db.Where("path = ?", where.GetPath())
    }

    if len(where.GetComponent()) > 0 {
        db = db.Where("component IN (?)", where.GetComponent())
    }

    if len(where.GetIsHidden()) > 0 {
        db = db.Where("is_hidden IN(?)", where.GetIsHidden())
    }

    if where.GetTitle() != "" {
        db = db.Where("title = ?", where.GetTitle())
    }

    if where.GetName() != "" {
        db = db.Where("name = ?", where.GetName())
    }

    if where.GetIcon() != "" {
        db = db.Where("icon = ?", where.GetIcon())
    }

    if where.GetMenuKey() != "" {
        if where.GetFuzzyMenuKey() {
            db = db.Where("menu_key LIKE ?", where.GetMenuKey()+"%")
        } else {
            db = db.Where("menu_key = ?", where.GetMenuKey())
        }
    }

    if len(where.GetMenuKeys()) > 0 {
        db = db.Where("menu_key IN (?)", where.GetMenuKeys())
    }

    if len(where.GetLevel()) > 0 {
        db = db.Where("level IN (?)", where.GetLevel())
    }

    if len(where.GetIsDelete()) > 0 {
        db = db.Where("is_delete IN (?)", where.GetIsDelete())
    }

    var total int32
    db.Count(&total)
    out = &pbs.PrivilegeMenuDeleteResp{
        AffectRows: total,
    }
    if total == 0 {
        return
    }

    err = db.Delete(_menuModel{}).Error
    if err != nil {
        logger.Logger.Error("_menuModel.Delete(%+v) error:%v", *in, err)
        return nil, err
    }

    return
}
