package lab

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/source/lab/catalog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderLabMenu = system.InitOrderSystem + 26

type initLabMenu struct{}

func init() {
	system.RegisterInit(initOrderLabMenu, &initLabMenu{})
}

func (i *initLabMenu) InitializerName() string {
	return "sys_base_menu_lab"
}

func (i *initLabMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initLabMenu) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initLabMenu) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	menus := catalog.BuildMenus()
	menuNameMap := make(map[string]uint, len(menus))
	for _, menu := range menus {
		parentID := uint(0)
		if menu.ParentName != "" {
			var hasParent bool
			parentID, hasParent = menuNameMap[menu.ParentName]
			if !hasParent {
				var parent sysModel.SysBaseMenu
				if err := db.Where("name = ?", menu.ParentName).First(&parent).Error; err != nil {
					return ctx, errors.Wrapf(err, "missing parent menu %s for %s", menu.ParentName, menu.Name)
				}
				parentID = parent.ID
				menuNameMap[menu.ParentName] = parentID
			}
		}

		saved, err := ensureMenu(db, sysModel.SysBaseMenu{
			MenuLevel: uint(menu.MenuLevel),
			Hidden:    menu.Hidden,
			ParentId:  parentID,
			Path:      menu.Path,
			Name:      menu.Name,
			Component: menu.Component,
			Sort:      menu.Sort,
			Meta: sysModel.Meta{
				Title: menu.Title,
				Icon:  menu.Icon,
			},
		})
		if err != nil {
			return ctx, err
		}
		menuNameMap[saved.Name] = saved.ID
	}
	return ctx, nil
}

func (i *initLabMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	for _, menu := range catalog.BuildMenus() {
		if errors.Is(db.Where("name = ?", menu.Name).First(&sysModel.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func ensureMenu(db *gorm.DB, menu sysModel.SysBaseMenu) (sysModel.SysBaseMenu, error) {
	var existing sysModel.SysBaseMenu
	err := db.Where("name = ?", menu.Name).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if createErr := db.Create(&menu).Error; createErr != nil {
			return sysModel.SysBaseMenu{}, createErr
		}
		return menu, nil
	}
	if err != nil {
		return sysModel.SysBaseMenu{}, err
	}

	existing.MenuLevel = menu.MenuLevel
	existing.Hidden = menu.Hidden
	existing.ParentId = menu.ParentId
	existing.Path = menu.Path
	existing.Component = menu.Component
	existing.Sort = menu.Sort
	existing.Meta = menu.Meta

	if saveErr := db.Save(&existing).Error; saveErr != nil {
		return sysModel.SysBaseMenu{}, saveErr
	}
	return existing, nil
}
