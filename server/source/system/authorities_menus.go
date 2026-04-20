package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Keep role-menu binding at the end of system seed stage so all menus
// (including lab seeds from other source packages) are already present.
const initOrderMenuAuthority = system.InitOrderSystem + 900

type initMenuAuthority struct{}

func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false
}

func (i *initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	authorities, err := loadAuthorities(ctx, db)
	if err != nil {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}

	allMenus, err := loadMenus(ctx, db)
	if err != nil {
		return ctx, errors.Wrap(err, "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "为超级管理员分配菜单失败")
	}
	if err = db.Model(&authorities[1]).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "为普通用户分配菜单失败")
	}
	if err = db.Model(&authorities[2]).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "为测试角色分配菜单失败")
	}

	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}

	var totalMenus int64
	if err := db.Model(&sysModel.SysBaseMenu{}).Count(&totalMenus).Error; err != nil {
		return false
	}
	if totalMenus == 0 {
		return false
	}

	requiredAuthorityIDs := []uint{888, 9528, 8881}
	for _, authorityID := range requiredAuthorityIDs {
		auth := &sysModel.SysAuthority{}
		ret := db.Model(auth).
			Where("authority_id = ?", authorityID).Preload("SysBaseMenus").First(auth)
		if ret == nil || ret.Error != nil {
			return false
		}
		if int64(len(auth.SysBaseMenus)) < totalMenus {
			return false
		}
	}
	return true
}

func loadAuthorities(ctx context.Context, db *gorm.DB) ([]sysModel.SysAuthority, error) {
	initAuth := &initAuthority{}
	if authorities, ok := ctx.Value(initAuth.InitializerName()).([]sysModel.SysAuthority); ok && len(authorities) >= 3 {
		return authorities, nil
	}

	requiredIDs := []uint{888, 9528, 8881}
	authorities := make([]sysModel.SysAuthority, 0, len(requiredIDs))
	for _, authorityID := range requiredIDs {
		var authority sysModel.SysAuthority
		if err := db.Where("authority_id = ?", authorityID).First(&authority).Error; err != nil {
			return nil, err
		}
		authorities = append(authorities, authority)
	}
	return authorities, nil
}

func loadMenus(ctx context.Context, db *gorm.DB) ([]sysModel.SysBaseMenu, error) {
	if menus, ok := ctx.Value(new(initMenu).InitializerName()).([]sysModel.SysBaseMenu); ok && len(menus) > 0 {
		return menus, nil
	}

	var menus []sysModel.SysBaseMenu
	if err := db.Order("sort asc, id asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return menus, nil
}
