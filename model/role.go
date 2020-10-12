package model

import (
    "errors"
    "fmt"
    "github.com/jinzhu/gorm"
    "gitlab.ceibsmoment.com/c/mp/client"
    "gitlab.ceibsmoment.com/c/mp/logger"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "math"
    "strconv"
    "time"
)

type _roleModel struct {
    DefaultModel
    Id         int32     `json:"id"`
    Name       string    `json:"name"`
    Desc       string    `json:"desc"`
    RoleKey    string    `json:"role_key"`
    Status     int32     `json:"status"`
    IsDelete   int32     `json:"is_delete"`
    CreateTime time.Time `json:"create_time"`
    UpdateTime time.Time `json:"update_time"`
}

var (
    RoleModel *_roleModel
)

func (*_roleModel) TableName() string {
    return "privilege_role"
}

func (u *_roleModel) Find(in *pbs.PrivilegeRoleFindReq) (out *pbs.PrivilegeRoleFindResp, err error) {
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

    if where.GetDesc() != "" {
        db = db.Where("desc LIKE ?", "%"+where.GetDesc()+"%")
    }

    if len(where.GetStatus()) > 0 {
        db = db.Where("status IN (?)", where.GetStatus())
    }

    if len(where.GetRoleKeys()) > 0 {
        db = db.Where("role_key IN (?)", where.GetRoleKeys())
    }

    if where.GetName() != "" {
        if where.GetFuzzyName() {
            db = db.Where("name LIKE ?", "%"+where.GetName()+"%")
        } else {
            db = db.Where("name = ?", where.GetName())
        }
    }

    var total int32
    db.Count(&total)

    if in.GetOrderBy() != "" {
        db = db.Order(in.GetOrderBy())
    }

    offset := (in.GetPage() - 1) * (in.GetSize())
    db = db.Limit(in.GetSize()).Offset(offset)

    objs := make([]*pbs.PrivilegeRole, 0)
    out = &pbs.PrivilegeRoleFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }
    ms := make([]*_roleModel, 0)
    err = db.Find(&ms).Error
    if err == gorm.ErrRecordNotFound {
        err = nil
    }

    if err != nil {
        logger.Logger.Errorf("_roleModel.Find(%v) error: %v", *in, err)
        return nil, err
    }

    for _, v := range ms {
        objs = append(objs, &pbs.PrivilegeRole{
            Id:         v.Id,
            Name:       v.Name,
            Desc:       v.Desc,
            RoleKey:    v.RoleKey,
            Status:     v.Status,
            CreateTime: int32(v.CreateTime.Unix()),
            UpdateTime: int32(v.UpdateTime.Unix()),
            IsDelete:   v.IsDelete,
        })
    }
    out = &pbs.PrivilegeRoleFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }

    return
}

func (u *_roleModel) Save(in *pbs.PrivilegeRoleSaveReq) (out *pbs.PrivilegeRoleSaveResp, err error) {
    //nowTime := int32(time.Now().Unix())
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    if in.GetId() > 0 {
        id := in.GetId()
        in.Id = 0

        var (
            m1 = &_roleModel{}
            m2 = &_roleModel{}
        )
        err = db.Model(&_roleModel{}).Where("id = ?", in.Id).First(&m1).Error
        if err == gorm.ErrRecordNotFound {
            err = nil
            return &pbs.PrivilegeRoleSaveResp{}, err
        }
        if err != nil {
            logger.Logger.Errorf("_roleModel.Save(%v) find same one error: %v", err)
            return nil, err
        }

        if in.Name != "" {
            err = db.Model(&_roleModel{}).Where("LENGTH(`role_key`) = ? AND `id` NOT IN(?) AND `name` = ?", len(m1.RoleKey), id, in.Name).First(&m2).Error
            if err != gorm.ErrRecordNotFound && err != nil {
                logger.Logger.Errorf("_roleModel.Save(%v) find same two error: %v", err)
                return nil, err
            }

            if err == nil && m2 != nil {
                err = errors.New("同级目录下已存在角色<" + in.Name + ">")
                logger.Logger.Errorf("_roleModel.Save() 同级目录已存在角色名: %s", in.Name)
                return nil, err
            }
        }

        err = db.Model(&_roleModel{}).Where("id = ?", id).Update(in).Error
        if err != nil {
            logger.Logger.Errorf("_roleModel.Save(%v) t1 error <%v>", *in, err)
            return
        }
        return &pbs.PrivilegeRoleSaveResp{Id: id}, nil
    }

    if in.GetName() == "" {
        logger.Logger.Errorf("_roleModel.Save name empty error: (in = %+v)", in)
        err = errors.New("name can not be null")
        return
    }

    roleKey, _ := u.GetRoleKey(in.GetRoleKey())

    var search = &_roleModel{}
    err = db.Model(&_roleModel{}).Where("LENGTH(`role_key`) = ? AND `name` = ?", len(roleKey), in.Name).First(&search).Error
    if err != gorm.ErrRecordNotFound && err != nil {
        logger.Logger.Errorf("_roleModel.Save insert find same name record error: (in = %+v) error: %v", *in, err)
        return nil, err
    }
    if err == nil && search != nil && search.Id > 0 {
        err = errors.New(fmt.Sprintf("已存在同名角色<%s>", in.Name))
        return nil, err
    }

    obj := &_roleModel{
        Id:         in.GetId(),
        Name:       in.GetName(),
        Desc:       in.GetDesc(),
        RoleKey:    roleKey,
        Status:     in.GetStatus(),
        CreateTime: time.Now(),
        UpdateTime: time.Now(),
        IsDelete:   in.GetIsDelete(),
    }
    err = db.Create(obj).Error
    if err != nil {
        logger.Logger.Errorf("_roleModel.Save(%#v) insert t4 error: <%+v>", obj, err)
        return
    }
    out = &pbs.PrivilegeRoleSaveResp{
        Id: obj.Id,
    }
    return
}

// 计算role_key
func (u *_roleModel) GetRoleKey(menuKey string) (key string, err error) {
    m := &_roleModel{}
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    if menuKey == "" {
        err = db.Order("role_key DESC").First(m).Error

        if err != nil && err != gorm.ErrRecordNotFound {
            logger.Logger.Errorf("_roleModel.GetRoleKey(%v) error: %v", menuKey, err)
            return "", err
        }

        if err == gorm.ErrRecordNotFound {
            // 搜索一级菜单
            err = db.Where("LENGTH(`key`) = 3").Order("`key` DESC").First(m).Error
            if err != nil && err != gorm.ErrRecordNotFound {
                logger.Logger.Errorf("_roleModel.GetRoleKey(%v) second find error: %v", menuKey, err)
                return "", err
            }

            if err == gorm.ErrRecordNotFound {
                key = "100"
                return key, nil
            }
        }

        k := m.RoleKey
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

    err = db.Where("role_key LIKE ?", menuKey+"%").Where("LENGTH(role_key) = ?", len(menuKey)+3).Order("role_key DESC").First(m).Error
    if err != nil && err != gorm.ErrRecordNotFound {
        logger.Logger.Errorf("_roleModel.GetRoleKey(%v) error: %v", menuKey, err)
        return "", err
    }
    if err == gorm.ErrRecordNotFound {
        key = menuKey + "001"
        return key, nil
    }

    k := m.RoleKey
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

func (u *_roleModel) Delete(in *pbs.PrivilegeRoleDeleteReq) (out *pbs.PrivilegeRoleDeleteResp, err error) {
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    where := in

    if len(where.GetId()) > 0 {
        db = db.Where("id IN (?)", where.GetId())
    }

    if len(where.GetStatus()) > 0 {
        db = db.Where("status IN (?)", where.GetStatus())
    }

    if len(where.GetRoleKeys()) > 0 {
        db = db.Where("role_key IN (?)", where.GetRoleKeys())
    }

    if where.GetDesc() != "" {
        db = db.Where("desc LIKE ?", "%"+where.GetDesc()+"%")
    }

    if where.GetName() != "" {
        if where.GetFuzzyName() {
            db = db.Where("name LIKE ?", "%"+where.GetName()+"%")
        } else {
            db = db.Where("name = ?", where.GetName())
        }
    }

    var total int32
    db.Count(&total)
    out = &pbs.PrivilegeRoleDeleteResp{
        AffectRows: total,
    }
    if total == 0 {
        return
    }

    err = db.Delete(_roleModel{}).Error
    if err != nil {
        logger.Logger.Error("_roleModel.Delete(%+v) error:%v", *in, err)
        return nil, err
    }

    return
}
