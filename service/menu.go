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

type _menuService struct {
}

var (
    MenuService *_menuService
)

// GetAllMenu
// 获取所有菜单
func (s *_menuService) GetAllMenu(c *gin.Context) (resp *pbs.Menu, err error) {
    baseMenu, err := model.MenuModel.Find(&pbs.PrivilegeMenuFindReq{
        Page: 1,
        Size: def.MaxUrlNum,
        Where: &pbs.PrivilegeMenuFindWhere{
            Level:    []int32{1},
            IsDelete: []int32{def.IsDeleteNormal},
            IsHidden: []int32{def.MenuIsHiddenFalse},
        },
    })
    if err != nil {
        logger.Logger.Errorf("_menuService.GetAllMenu model.MenuModel.Find<baseMenu> error: %v", err)
        return
    }

    if baseMenu.GetTotalSize() == 0 {
        logger.Logger.Errorf("_menuService.GetAllMenu model.MenuModel.Find<baseMenu> error: %s", "未配置前台菜单")
        err = errors.New("未配置前台路由")
        return nil, err
    }

    var wg sync.WaitGroup
    wg.Add(len(baseMenu.GetList()))

    var items []*pbs.MenuItem
    var ids []int32
    var itemMaps = make(map[int32][]*pbs.MenuItem)
    var mutex sync.RWMutex
    for _, v := range baseMenu.GetList() {
        ids = append(ids, v.GetId())
        go func(wg *sync.WaitGroup, v *pbs.PrivilegeMenu) {
            mInfo, err := s.recvCalMenu(wg, v)
            defer wg.Done()
            if err != nil {
                logger.Logger.Errorf("_menuService.GetAllMenu s.recvCalMenu(wg, v) error: %v", err)
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
    resp = &pbs.Menu{
        Items: items,
    }

    return
}

func (s *_menuService) recvCalMenu(wg *sync.WaitGroup, menu *pbs.PrivilegeMenu) (*pbs.MenuItem, error) {
    var err error
    var cur *pbs.MenuItem
    var menuItem = &pbs.MenuItem{}
    cur = &pbs.MenuItem{
        Id:           menu.GetId(),
        MenuKey:      menu.GetMenuKey(),
        Level:        menu.GetLevel(),
        DisplayOrder: menu.GetDisplayOrder(),
        Path:         menu.GetPath(),
        Component:    menu.GetComponent(),
        // Redirect: "",
        Name:  menu.GetName(),
        Title: menu.GetTitle(),
        Icon:  menu.GetIcon(),
    }
    menuItem = cur
    where := &pbs.PrivilegeMenuFindWhere{
        MenuKey:      menu.GetMenuKey(),
        FuzzyMenuKey: true,
        Level: []int32{
            int32(len(menu.GetMenuKey())/3 + 1),
        },
        IsDelete: []int32{def.IsDeleteNormal},
        IsHidden: []int32{def.MenuIsHiddenFalse},
    }
    menus, err := model.MenuModel.Find(&pbs.PrivilegeMenuFindReq{
        Page:    1,
        Size:    def.MaxUrlNum,
        Where:   where,
        OrderBy: "display_order ASC",
    })

    if err != nil {
        logger.Logger.Errorf("_menuService.recvCalMenu model.MenuModel.Find where = %v, error = %v", where, err)
        return menuItem, nil
    }

    if len(menus.GetList()) == 0 {
        return menuItem, err
    }

    for _, v := range menus.GetList() {
        next, err := s.recvCalMenu(wg, v)
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
func (s *_menuService) GetMenu(c *gin.Context, userId int32) (resp *pbs.Menu, err error) {
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
        logger.Logger.Errorf("_menuService.getCommonAdminMenu model.PrivilegeAdmin.Find(userId = %d) error: %v", userId, err)
        return nil, err
    }
    if user.GetTotalSize() == 0 {
        logger.Logger.Errorf("_menuService.getCommonAdminMenu model.PrivilegeAdmin.Find(userId = %d) error: %v", userId, "管理员不存在")
        return nil, errors.New("管理员不存在")
    }

    roleKey := user.GetList()[0].GetRoleKey()
    if roleKey == "" {
        err = errors.New("未分配任何角色")
        return nil, err
    }
    isAdmin := RoleService.CheckIsAdmin(roleKey)
    if isAdmin {
        return s.GetAllMenu(c)
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
                logger.Logger.Errorf("_menuService.getCommonAdminMenu model.RoleMenuModel.Find(role_key = %s) error: %v", v, err)
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
            IsHidden: []int32{def.MenuIsHiddenFalse},
        },
        OrderBy: "level ASC, display_order ASC",
    })
    if err != nil {
        logger.Logger.Errorf("_menuService.getCommonAdminMenu model.MenuModel.Find(%v) error: %v", menuKeys, err)
        return nil, err
    }
    if menus.GetTotalSize() == 0 {
        return &pbs.Menu{
            Items: []*pbs.MenuItem{},
        }, nil
    }

    // 处理菜单
    return s.process(menus.GetList())
}

// process
func (s *_menuService) process(menus []*pbs.PrivilegeMenu) (resp *pbs.Menu, err error) {
    menuItems := make([]*pbs.MenuItem, 0)
    // 处理一级菜单
    for _, v := range menus {
        if v.GetLevel() == 1 {
            menuItems = append(menuItems, &pbs.MenuItem{
                Id:           v.GetId(),
                MenuKey:      v.GetMenuKey(),
                Level:        v.GetLevel(),
                DisplayOrder: v.GetDisplayOrder(),
                Path:         v.GetPath(),
                Component:    v.GetComponent(),
                Name:         v.GetName(),
                Title:        v.GetTitle(),
                Icon:         v.GetIcon(),
                BasePath:     "/" + strings.Trim(v.GetPath(), "/"),
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
                        menuItems[index].Children = make([]*pbs.MenuItem, 0)
                    }
                    menuItems[index].Children = append(menuItems[index].Children, &pbs.MenuItem{
                        Id:           v.GetId(),
                        MenuKey:      v.GetMenuKey(),
                        Level:        v.GetLevel(),
                        DisplayOrder: v.GetDisplayOrder(),
                        Path:         v.GetPath(),
                        Component:    v.GetComponent(),
                        Name:         v.GetName(),
                        Title:        v.GetTitle(),
                        Icon:         v.GetIcon(),
                        BasePath:     v2.GetBasePath() + "/" + strings.Trim(v.GetPath(), "/"),
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
                        logger.Logger.Errorf("_menuService.process 菜单<%v>没有配置父级节点", *v)
                        continue
                    }

                    for idx2, v3 := range v2.GetChildren() {
                        if v3.GetMenuKey() == v.GetMenuKey()[:6] {
                            if menuItems[index].Children[idx2].Children == nil {
                                menuItems[index].Children[idx2].Children = make([]*pbs.MenuItem, 0)
                            }
                            menuItems[index].Children[idx2].Children = append(menuItems[index].Children[idx2].Children, &pbs.MenuItem{
                                Id:           v.GetId(),
                                MenuKey:      v.GetMenuKey(),
                                Level:        v.GetLevel(),
                                DisplayOrder: v.GetDisplayOrder(),
                                Path:         v.GetPath(),
                                Component:    v.GetComponent(),
                                Name:         v.GetName(),
                                Title:        v.GetTitle(),
                                Icon:         v.GetIcon(),
                                BasePath:     v3.GetBasePath() + "/" + strings.Trim(v.GetPath(), "/"),
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
                        logger.Logger.Errorf("_menuService.process 二级菜单<%v>没有配置父级节点", *v)
                        continue
                    }

                    for idx2, v3 := range v2.GetChildren() {
                        if v3.GetMenuKey() == v.GetMenuKey()[:6] {
                            if menuItems[index].Children[idx2].Children == nil || len(menuItems[index].Children[idx2].Children) == 00 {
                                logger.Logger.Errorf("_menuService.process 三级菜单<%v>没有配置父级节点", *v)
                                continue
                            }

                            for idx3, v4 := range v3.GetChildren() {
                                if v4.GetMenuKey() == v.GetMenuKey()[:9] {
                                    if menuItems[index].Children[idx2].Children[idx3].Children == nil {
                                        menuItems[index].Children[idx2].Children[idx3].Children = make([]*pbs.MenuItem, 0)
                                    }
                                    menuItems[index].Children[idx2].Children[idx3].Children = append(menuItems[index].Children[idx2].Children[idx3].Children, &pbs.MenuItem{
                                        Id:           v.GetId(),
                                        MenuKey:      v.GetMenuKey(),
                                        Level:        v.GetLevel(),
                                        DisplayOrder: v.GetDisplayOrder(),
                                        Path:         v.GetPath(),
                                        Component:    v.GetComponent(),
                                        Name:         v.GetName(),
                                        Title:        v.GetTitle(),
                                        Icon:         v.GetIcon(),
                                        BasePath:     v4.GetBasePath() + "/" + strings.Trim(v.GetPath(), "/"),
                                    })
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    // 重制菜单url
    menuItems = s.resetMenuUrl(menuItems)
    resp = &pbs.Menu{
        Items: menuItems,
    }
    return
}

// resetMenuUrl
// 重置菜单路由，以最底下菜单路由为准
func (s *_menuService) resetMenuUrl(items []*pbs.MenuItem) []*pbs.MenuItem {
    for k, item := range items {
        items[k] = s.dynamicCalMenuUrl(item)
        if item.GetChildren() != nil && len(item.GetChildren()) > 0 {
            item.Children = s.resetMenuUrl(item.GetChildren())
            items[k] = item
        }
    }

    return items
}

// dynamicCalMenuUrl
// 动态计算菜单最低级路由
func (s *_menuService) dynamicCalMenuUrl(item *pbs.MenuItem) *pbs.MenuItem {
    if item.GetChildren() == nil || len(item.GetChildren()) == 0 {
        item.Redirect = item.GetBasePath()
        return item
    }
    item.Redirect = item.GetBasePath() + "/" + strings.Trim(item.GetChildren()[0].GetPath(), "/")
    return s.dynamicCalMenuUrl(item.GetChildren()[0])
}
