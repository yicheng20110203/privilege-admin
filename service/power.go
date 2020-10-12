// WARNING 最多支持4级菜
// 每级菜单最多支持999个
package service

import (
    "errors"
    "github.com/gin-gonic/gin"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/model"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "strings"
    "sync"
)

type _powerService struct {
}

var (
    PowerService *_powerService
)

// GetAllMenu
// 获取所有权限
func (s *_powerService) GetAllPower(c *gin.Context) (resp *pbs.Power, err error) {
    baseMenu, err := model.MenuModel.Find(&pbs.PrivilegeMenuFindReq{
        Page: 1,
        Size: def.MaxUrlNum,
        Where: &pbs.PrivilegeMenuFindWhere{
            Level:    []int32{1},
            IsDelete: []int32{def.IsDeleteNormal},
        },
    })
    if err != nil {
        logger.Logger.Errorf("_powerService.GetAllMenu model.MenuModel.Find<baseMenu> error: %v", err)
        return
    }

    if baseMenu.GetTotalSize() == 0 {
        logger.Logger.Errorf("_powerService.GetAllMenu model.MenuModel.Find<baseMenu> error: %s", "未配置前台菜单")
        err = errors.New("未配置前台路由")
        return nil, err
    }

    var wg sync.WaitGroup
    wg.Add(len(baseMenu.GetList()))

    var items []*pbs.PowerItem
    var ids []int32
    var itemMaps = make(map[int32][]*pbs.PowerItem)
    var mutex sync.RWMutex
    for _, v := range baseMenu.GetList() {
        ids = append(ids, v.GetId())
        go func(wg *sync.WaitGroup, v *pbs.PrivilegeMenu) {
            mInfo, err := s.recvCallPower(wg, v)
            defer wg.Done()
            if err != nil {
                logger.Logger.Errorf("_powerService.GetAllMenu s.recvCallPower(wg, v) error: %v", err)
                return
            }
            mutex.Lock()
            itemMaps[v.GetId()] = append(itemMaps[v.GetId()], mInfo)
            defer mutex.Unlock()
        }(&wg, v)
    }
    wg.Wait()

    for _, id := range ids {
        if v, ok := itemMaps[id]; ok {
            items = append(items, v...)
        }
    }
    resp = &pbs.Power{
        Items: items,
    }

    return
}

func (s *_powerService) recvCallPower(wg *sync.WaitGroup, menu *pbs.PrivilegeMenu) (*pbs.PowerItem, error) {
    var err error
    var cur *pbs.PowerItem
    var menuItem = &pbs.PowerItem{}
    cur = &pbs.PowerItem{
        MenuKey:      menu.GetMenuKey(),
        Level:        menu.GetLevel(),
        DisplayOrder: menu.GetDisplayOrder(),
        Name:         menu.GetName(),
        Title:        menu.GetTitle(),
    }
    menuItem = cur
    where := &pbs.PrivilegeMenuFindWhere{
        MenuKey:      menu.GetMenuKey(),
        FuzzyMenuKey: true,
        Level:        []int32{int32(len(menu.GetMenuKey())/3 + 1)},
        IsDelete:     []int32{def.IsDeleteNormal},
    }
    menus, err := model.MenuModel.Find(&pbs.PrivilegeMenuFindReq{
        Page:    1,
        Size:    def.MaxUrlNum,
        Where:   where,
        OrderBy: "display_order ASC",
    })

    if err != nil {
        logger.Logger.Errorf("_powerService.recvCallPower model.MenuModel.Find where = %v, error = %v", where, err)
        return menuItem, nil
    }

    if len(menus.GetList()) == 0 {
        return menuItem, err
    }

    for _, v := range menus.GetList() {
        next, err := s.recvCallPower(wg, v)
        if err == nil && next != nil {
            cur.Children = append(cur.GetChildren(), next)
            continue
        }
    }

    return menuItem, err
}

// GetMenu functions:
// 获取超管菜单
// 获取普通管理员菜单
func (s *_powerService) GetPower(c *gin.Context, userId int32) (resp *pbs.Power, err error) {
    user, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
        Page: 1,
        Size: 1,
        Where: &pbs.PrivilegeAdminFindWhere{
            Id: []int32{
                userId,
            },
        },
    })
    if err != nil {
        logger.Logger.Errorf("_powerService.getCommonAdminMenu model.PrivilegeAdmin.Find(userId = %d) error: %v", userId, err)
        return nil, err
    }
    if user.GetTotalSize() == 0 {
        logger.Logger.Errorf("_powerService.getCommonAdminMenu model.PrivilegeAdmin.Find(userId = %d) error: %v", userId, "管理员不存在")
        return nil, errors.New("管理员不存在")
    }

    roleKey := user.GetList()[0].GetRoleKey()
    if roleKey == "" {
        err = errors.New("未分配任何角色")
        return nil, err
    }
    isAdmin := RoleService.CheckIsAdmin(roleKey)
    if isAdmin {
        return s.GetAllPower(c)
    }

    roleKeys := strings.Split(roleKey, ",")
    menuKeys := make([]string, 0)
    var wg sync.WaitGroup
    wg.Add(len(roleKeys))
    for _, v := range roleKeys {
        go func(v string) {
            // 普通管理员，非超管
            roleMenuMaps, err := model.RoleMenuModel.Find(&pbs.PrivilegeRoleMenuFindReq{
                Page: 1,
                Size: def.MaxUrlNum,
                Where: &pbs.PrivilegeRoleMenuFindWhere{
                    RoleKey:      v,
                    FuzzyRoleKey: true,
                    IsDelete:     []int32{def.IsDeleteNormal},
                },
            })
            if err != nil {
                logger.Logger.Errorf("_powerService.getCommonAdminMenu model.RoleMenuModel.Find(role_key = %s) error: %v", v, err)
                return
            }
            for _, l := range roleMenuMaps.GetList() {
                menuKeys = append(menuKeys, l.GetMenuKey())
            }
            defer wg.Done()
        }(v)
    }
    wg.Wait()

    // 获取排序菜单
    menus, err := model.MenuModel.Find(&pbs.PrivilegeMenuFindReq{
        Page: 1,
        Size: def.MaxUrlNum,
        Where: &pbs.PrivilegeMenuFindWhere{
            MenuKeys: menuKeys,
            IsDelete: []int32{def.IsDeleteNormal},
        },
        OrderBy: "level ASC, display_order ASC",
    })
    if err != nil {
        logger.Logger.Errorf("_powerService.getCommonAdminMenu model.MenuModel.Find(%v) error: %v", menuKeys, err)
        return nil, err
    }
    if menus.GetTotalSize() == 0 {
        return &pbs.Power{
            Items: []*pbs.PowerItem{},
        }, nil
    }

    // 处理菜单
    return s.process(menus.GetList())
}

// process
func (s *_powerService) process(menus []*pbs.PrivilegeMenu) (resp *pbs.Power, err error) {
    menuItems := make([]*pbs.PowerItem, 0)
    // 处理一级菜单
    for _, v := range menus {
        if v.GetLevel() == 1 {
            menuItems = append(menuItems, &pbs.PowerItem{
                MenuKey:      v.GetMenuKey(),
                Level:        v.GetLevel(),
                DisplayOrder: v.GetDisplayOrder(),
                Name:         v.GetName(),
                Title:        v.GetTitle(),
            })
            continue
        }
    }

    // 处理二级菜单
    for _, v := range menus {
        if v.GetLevel() == 2 {
            for index, v2 := range menuItems {
                if v2.GetMenuKey() == v.GetMenuKey()[:3] {
                    if menuItems[index].Children == nil {
                        menuItems[index].Children = make([]*pbs.PowerItem, 0)
                    }
                    menuItems[index].Children = append(menuItems[index].Children, &pbs.PowerItem{
                        MenuKey:      v.GetMenuKey(),
                        Level:        v.GetLevel(),
                        DisplayOrder: v.GetDisplayOrder(),
                        Name:         v.GetName(),
                        Title:        v.GetTitle(),
                    })
                }
            }
        }
    }

    // 处理三级级菜单
    for _, v := range menus {
        if v.GetLevel() == 3 {
            for index, v2 := range menuItems {
                if v2.GetMenuKey() == v.GetMenuKey()[:3] {
                    if menuItems[index].Children == nil || len(menuItems[index].Children) == 0 {
                        logger.Logger.Errorf("_powerService.process 菜单<%v>没有配置父级节点", *v)
                        continue
                    }

                    for idx2, v3 := range v2.GetChildren() {
                        if v3.GetMenuKey() == v.GetMenuKey()[:6] {
                            if menuItems[index].Children[idx2].Children == nil {
                                menuItems[index].Children[idx2].Children = make([]*pbs.PowerItem, 0)
                            }
                            menuItems[index].Children[idx2].Children = append(menuItems[index].Children[idx2].Children, &pbs.PowerItem{
                                MenuKey:      v.GetMenuKey(),
                                Level:        v.GetLevel(),
                                DisplayOrder: v.GetDisplayOrder(),
                                Name:         v.GetName(),
                                Title:        v.GetTitle(),
                            })
                        }
                    }
                }
            }
        }
    }

    // 处理四级级菜单
    for _, v := range menus {
        if v.GetLevel() == 4 {
            for index, v2 := range menuItems {
                if v2.GetMenuKey() == v.GetMenuKey()[:3] {
                    if menuItems[index].Children == nil || len(menuItems[index].Children) == 0 {
                        logger.Logger.Errorf("_powerService.process 二级菜单<%v>没有配置父级节点", *v)
                        continue
                    }

                    for idx2, v3 := range v2.GetChildren() {
                        if v3.GetMenuKey() == v.GetMenuKey()[:6] {
                            if menuItems[index].Children[idx2].Children == nil || len(menuItems[index].Children[idx2].Children) == 00 {
                                logger.Logger.Errorf("_powerService.process 三级菜单<%v>没有配置父级节点", *v)
                                continue
                            }

                            for idx3, v4 := range v3.GetChildren() {
                                if v4.GetMenuKey() == v.GetMenuKey()[:9] {
                                    if menuItems[index].Children[idx2].Children[idx3].Children == nil {
                                        menuItems[index].Children[idx2].Children[idx3].Children = make([]*pbs.PowerItem, 0)
                                    }
                                    menuItems[index].Children[idx2].Children[idx3].Children = append(menuItems[index].Children[idx2].Children[idx3].Children, &pbs.PowerItem{
                                        MenuKey:      v.GetMenuKey(),
                                        Level:        v.GetLevel(),
                                        DisplayOrder: v.GetDisplayOrder(),
                                        Name:         v.GetName(),
                                        Title:        v.GetTitle(),
                                    })
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    resp = &pbs.Power{
        Items: menuItems,
    }
    return
}

func (*_powerService) GetPowerByRoleKey(c *gin.Context, roleKey string) (resp map[string]string, err error) {
    resp = make(map[string]string)
    powers, err := model.RoleMenuModel.Find(&pbs.PrivilegeRoleMenuFindReq{
        Page: 1,
        Size: def.MaxUrlNum,
        Where: &pbs.PrivilegeRoleMenuFindWhere{
            RoleKey:      roleKey,
            FuzzyRoleKey: true,
            IsDelete:     []int32{def.IsDeleteNormal},
        },
    })
    if err != nil {
        logger.Logger.Error("_powerService.GetPowerByRoleKey model.RoleMenuModel.Find roleKey = %s,  error: %+v", roleKey, err)
        return resp, err
    }

    if powers.GetTotalSize() == 0 {
        return
    }

    for _, v := range powers.GetList() {
        resp[v.GetMenuKey()] = v.GetMenuKey()
    }

    return
}
