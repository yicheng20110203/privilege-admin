package model

import (
    "errors"
    "fmt"
    "github.com/jinzhu/gorm"
    "gitlab.ceibsmoment.com/c/mp/client"
    "gitlab.ceibsmoment.com/c/mp/library"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/mppbs"
    "math"
    "time"
)

type _privilegeAdmin struct {
    DefaultModel
    Id         int32     `json:"id"`
    LoginName  string    `json:"login_name"`
    Password   string    `json:"password"`
    Username   string    `json:"username"`
    Avatar     string    `json:"avatar"`
    Salt       string    `json:"salt"`
    DepKey     string    `json:"dep_key"`
    RoleKey    string    `json:"role_key"`
    IsAdmin    int32     `json:"is_admin"`
    CreateTime time.Time `json:"create_time"`
    UpdateTime time.Time `json:"update_time"`
}

var (
    PrivilegeAdmin *_privilegeAdmin
)

func (*_privilegeAdmin) TableName() string {
    return "privilege_admin"
}

func (u *_privilegeAdmin) Find(in *mppbs.PrivilegeAdminFindReq) (out *mppbs.PrivilegeAdminFindResp, err error) {
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

    if where.GetLoginName() != "" {
        db = db.Where("login_name = ?", where.GetLoginName())
    }

    if where.GetUsername() != "" {
        db = db.Where("username LIKE ?", "%"+where.GetUsername()+"%")
    }

    if where.GetRoleKey() != "" {
        db = db.Where("role_key LIKE ?", where.GetRoleKey()+"%")
    }

    if where.GetFilterAdmin() {
        db = db.Where("is_admin = ?", 0)
    } else {
        if where.GetIsAdmin() > 0 {
            db = db.Where("is_admin = ?", where.GetIsAdmin())
        }
    }

    var total int32
    db.Count(&total)

    if in.GetOrderBy() != "" {
        db = db.Order(in.GetOrderBy())
    } else {
        db = db.Order("create_time DESC")
    }

    offset := (in.GetPage() - 1) * (in.GetSize())
    db = db.Limit(in.GetSize()).Offset(offset)

    objs := make([]*mppbs.PrivilegeAdmin, 0)
    out = &mppbs.PrivilegeAdminFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }

    ms := make([]*_privilegeAdmin, 0)
    err = db.Find(&ms).Error
    if err == gorm.ErrRecordNotFound {
        err = nil
    }

    if err != nil {
        return nil, err
    }

    for _, v := range ms {
        objs = append(objs, &mppbs.PrivilegeAdmin{
            Id:         v.Id,
            LoginName:  v.LoginName,
            Password:   v.Password,
            Username:   v.Username,
            Avatar:     v.Avatar,
            Salt:       v.Salt,
            DepKey:     v.DepKey,
            RoleKey:    v.RoleKey,
            IsAdmin:    v.IsAdmin,
            CreateTime: int32(v.CreateTime.Unix()),
            UpdateTime: int32(v.UpdateTime.Unix()),
        })
    }

    out = &mppbs.PrivilegeAdminFindResp{
        Page:      in.GetPage(),
        Size:      in.GetSize(),
        TotalSize: int32(total),
        TotalPage: int32(math.Ceil(float64(total) / float64(in.GetSize()))),
        List:      objs,
    }

    return
}

func (u *_privilegeAdmin) getPasswordSalt(pwd string) (password string, salt string) {
    rand := time.Now().UnixNano()
    rand = rand >> 20
    salt = fmt.Sprintf("%d", rand)
    password = library.MD5(pwd + salt)
    return
}

func (u *_privilegeAdmin) Save(in *mppbs.PrivilegeAdminSaveReq) (out *mppbs.PrivilegeAdminSaveResp, err error) {
    //nowTime := int32(time.Now().Unix())
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    if in.GetId() > 0 {
        id := in.GetId()
        in.Id = 0
        //in.UpdateTime = nowTime
        // update password
        if in.GetPassword() != "" {
            password, salt := u.getPasswordSalt(in.GetPassword())
            in.Password = password
            in.Salt = salt
        }

        err = db.Model(&_privilegeAdmin{}).Where("id = ?", id).Update(in).Error
        if err != nil {
            logger.Logger.Errorf("_privilegeAdmin.Save(%v) t1 error <%v>", *in, err)
            return
        }
        return &mppbs.PrivilegeAdminSaveResp{Id: id}, nil
    }

    if in.GetLoginName() == "" || in.GetPassword() == "" {
        logger.Logger.Errorf("_privilegeAdmin.Save login_name or password empty error: (in = %+v)", in)
        err = errors.New("login_name can not be null")
        return
    }

    one, err := u.Find(&mppbs.PrivilegeAdminFindReq{
        Page: 1,
        Size: 1,
        Where: &mppbs.PrivilegeAdminFindWhere{
            LoginName: in.GetLoginName(),
        },
    })
    if err != nil {
        logger.Logger.Errorf("_privilegeAdmin.Save u.Find(%v) error: %v", in.GetLoginName(), err)
        return nil, err
    }

    if len(one.GetList()) > 0 {
        logger.Logger.Errorf("_privilegeAdmin.Save login_name(%s) has been exist error", in.GetLoginName())
        err = errors.New(fmt.Sprintf("login_name(%s) has been exist", in.GetLoginName()))
        return nil, err
    }

    password, salt := u.getPasswordSalt(in.GetPassword())
    obj := &_privilegeAdmin{
        Id:         in.GetId(),
        LoginName:  in.GetLoginName(),
        Password:   password,
        Username:   in.GetUsername(),
        Avatar:     in.GetAvatar(),
        Salt:       salt,
        DepKey:     in.GetDepKey(),
        RoleKey:    in.GetRoleKey(),
        IsAdmin:    in.GetIsAdmin(),
        CreateTime: time.Now(),
        UpdateTime: time.Now(),
    }
    err = db.Create(obj).Error
    if err != nil {
        logger.Logger.Errorf("_privilegeAdmin.Save(%#v) insert t4 error: <%+v>", obj, err)
        return
    }
    out = &mppbs.PrivilegeAdminSaveResp{
        Id: obj.Id,
    }
    return
}

func (u *_privilegeAdmin) Delete(in *mppbs.PrivilegeAdminDeleteReq) (out *mppbs.PrivilegeAdminDeleteResp, err error) {
    db := client.GetMysqlDb()
    db = db.Table(u.TableName())

    where := in

    if len(where.GetId()) > 0 {
        db = db.Where("id IN (?)", where.GetId())
    }

    if where.GetLoginName() != "" {
        db = db.Where("login_name LIKE ?", "%"+where.GetLoginName()+"%")
    }

    if where.GetUsername() != "" {
        db = db.Where("username LIKE ?", "%"+where.GetUsername()+"%")
    }

    if where.GetRoleKey() != "" {
        db = db.Where("role_key LIKE ?", where.GetRoleKey()+"%")
    }

    if where.GetIsAdmin() > 0 {
        db = db.Where("is_admin = ?", where.GetIsAdmin())
    }

    var total int32
    db.Count(&total)
    out = &mppbs.PrivilegeAdminDeleteResp{
        AffectRows: total,
    }
    if total == 0 {
        return
    }

    err = db.Delete(_privilegeAdmin{}).Error
    if err != nil {
        logger.Logger.Error("_privilegeAdmin.Delete(%+v) error:%v", *in, err)
        return nil, err
    }

    return
}
