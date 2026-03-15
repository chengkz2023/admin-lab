package system

import (
	"context"

	. "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i *initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	parentMenus := []SysBaseMenu{
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 1, Meta: Meta{Title: "超级管理员", Icon: "user"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "lab", Name: "lab", Component: "view/lab/index.vue", Sort: 2, Meta: Meta{Title: "实验室", Icon: "data-analysis"}},
	}

	menuNameMap := make(map[string]uint)
	allMenus := make([]SysBaseMenu, 0, len(parentMenus)+10)

	for _, menu := range parentMenus {
		savedMenu, saveErr := ensureMenu(db, menu)
		if saveErr != nil {
			return ctx, errors.Wrap(saveErr, SysBaseMenu{}.TableName()+"父级菜单初始化失败!")
		}
		menuNameMap[savedMenu.Name] = savedMenu.ID
		allMenus = append(allMenus, savedMenu)
	}

	childMenus := []SysBaseMenu{
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: Meta{Title: "角色管理", Icon: "avatar"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "api", Name: "api", Component: "view/superAdmin/api/api.vue", Sort: 3, Meta: Meta{Title: "API管理", Icon: "platform", KeepAlive: true}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: Meta{Title: "用户管理", Icon: "coordinate"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "dictionary", Name: "dictionary", Component: "view/superAdmin/dictionary/sysDictionary.vue", Sort: 5, Meta: Meta{Title: "字典管理", Icon: "notebook"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: Meta{Title: "操作历史", Icon: "pie-chart"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "sysParams", Name: "sysParams", Component: "view/superAdmin/params/sysParams.vue", Sort: 7, Meta: Meta{Title: "参数管理", Icon: "compass"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["lab"], Path: "simulation", Name: "labSimulation", Component: "view/lab/simulation/index.vue", Sort: 1, Meta: Meta{Title: "需求仿真", Icon: "document"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["lab"], Path: "component-demo", Name: "labComponentDemo", Component: "view/lab/component-demo/index.vue", Sort: 2, Meta: Meta{Title: "组件示例", Icon: "magic-stick"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["lab"], Path: "reusable", Name: "labReusable", Component: "view/lab/reusable/index.vue", Sort: 3, Meta: Meta{Title: "复用组件", Icon: "files"}},
	}

	for _, menu := range childMenus {
		savedMenu, saveErr := ensureMenu(db, menu)
		if saveErr != nil {
			return ctx, errors.Wrap(saveErr, SysBaseMenu{}.TableName()+"子菜单初始化失败!")
		}
		allMenus = append(allMenus, savedMenu)
	}

	next = context.WithValue(ctx, i.InitializerName(), allMenus)
	return next, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "admin").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("name = ?", "lab").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("name = ?", "labSimulation").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("name = ?", "labComponentDemo").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	if errors.Is(db.Where("name = ?", "labReusable").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func ensureMenu(db *gorm.DB, menu SysBaseMenu) (SysBaseMenu, error) {
	var existing SysBaseMenu
	err := db.Where("name = ?", menu.Name).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if createErr := db.Create(&menu).Error; createErr != nil {
			return SysBaseMenu{}, createErr
		}
		return menu, nil
	}
	if err != nil {
		return SysBaseMenu{}, err
	}

	existing.MenuLevel = menu.MenuLevel
	existing.Hidden = menu.Hidden
	existing.ParentId = menu.ParentId
	existing.Path = menu.Path
	existing.Component = menu.Component
	existing.Sort = menu.Sort
	existing.Meta = menu.Meta

	if saveErr := db.Save(&existing).Error; saveErr != nil {
		return SysBaseMenu{}, saveErr
	}
	return existing, nil
}
