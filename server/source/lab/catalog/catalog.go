package catalog

import (
	"fmt"
	"sort"
)

const (
	AreaSimulation    = "simulation"
	AreaComponentDemo = "component-demo"
	AreaReusable      = "reusable"
)

var defaultRoleIDs = []string{"888", "8881", "9528"}

type MenuItem struct {
	MenuLevel  int
	Hidden     bool
	ParentName string
	Path       string
	Name       string
	Component  string
	Sort       int
	Title      string
	Icon       string
}

type APIItem struct {
	APIGroup    string
	Method      string
	Path        string
	Description string
}

type CasbinPolicyItem struct {
	RoleID string
	Path   string
	Method string
}

type PageSpec struct {
	Key            string
	Area           string
	Menu           MenuItem
	APIs           []APIItem
	DefaultRoleIDs []string
}

var pageRegistry []PageSpec

func registerPage(page PageSpec) {
	if page.Area != AreaSimulation && page.Area != AreaComponentDemo && page.Area != AreaReusable {
		panic(fmt.Sprintf("unknown area %s for page %s", page.Area, page.Key))
	}
	if page.Menu.Name == "" {
		panic(fmt.Sprintf("missing menu name for page %s", page.Key))
	}
	if page.Menu.Path == "" {
		panic(fmt.Sprintf("missing menu path for page %s", page.Key))
	}
	for _, existing := range pageRegistry {
		if existing.Key == page.Key {
			panic(fmt.Sprintf("duplicate page key %s", page.Key))
		}
		if existing.Menu.Name == page.Menu.Name {
			panic(fmt.Sprintf("duplicate menu name %s", page.Menu.Name))
		}
	}
	pageRegistry = append(pageRegistry, page)
}

func BuildMenus() []MenuItem {
	menus := make([]MenuItem, 0, len(pageRegistry)+4)
	menus = append(menus,
		MenuItem{
			MenuLevel:  0,
			Hidden:     false,
			ParentName: "",
			Path:       "lab",
			Name:       "lab",
			Component:  "view/lab/index.vue",
			Sort:       2,
			Title:      "实验室",
			Icon:       "data-analysis",
		},
		MenuItem{
			MenuLevel:  1,
			Hidden:     false,
			ParentName: "lab",
			Path:       "simulation",
			Name:       "labSimulation",
			Component:  "view/routerHolder.vue",
			Sort:       1,
			Title:      "需求仿真",
			Icon:       "document",
		},
		MenuItem{
			MenuLevel:  1,
			Hidden:     false,
			ParentName: "lab",
			Path:       "component-demo",
			Name:       "labComponentDemo",
			Component:  "view/routerHolder.vue",
			Sort:       2,
			Title:      "组件示例",
			Icon:       "magic-stick",
		},
		MenuItem{
			MenuLevel:  1,
			Hidden:     false,
			ParentName: "lab",
			Path:       "reusable",
			Name:       "labReusable",
			Component:  "view/routerHolder.vue",
			Sort:       3,
			Title:      "复用组件",
			Icon:       "files",
		},
	)

	for _, page := range sortedPages() {
		menu := page.Menu
		menu.MenuLevel = 2
		menu.Hidden = false
		menu.ParentName = areaParentMenuName(page.Area)
		menus = append(menus, menu)
	}

	sort.SliceStable(menus, func(i, j int) bool {
		if menus[i].MenuLevel != menus[j].MenuLevel {
			return menus[i].MenuLevel < menus[j].MenuLevel
		}
		if menus[i].ParentName != menus[j].ParentName {
			return menus[i].ParentName < menus[j].ParentName
		}
		if menus[i].Sort != menus[j].Sort {
			return menus[i].Sort < menus[j].Sort
		}
		return menus[i].Name < menus[j].Name
	})
	return menus
}

func BuildAPIs() []APIItem {
	apis := make([]APIItem, 0)
	seen := make(map[string]struct{})
	for _, page := range sortedPages() {
		for _, api := range page.APIs {
			key := api.Path + "::" + api.Method
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			apis = append(apis, api)
		}
	}
	sort.SliceStable(apis, func(i, j int) bool {
		if apis[i].Path != apis[j].Path {
			return apis[i].Path < apis[j].Path
		}
		return apis[i].Method < apis[j].Method
	})
	return apis
}

func BuildCasbinPolicies() []CasbinPolicyItem {
	policies := make([]CasbinPolicyItem, 0)
	seen := make(map[string]struct{})

	for _, page := range sortedPages() {
		roleIDs := page.DefaultRoleIDs
		if len(roleIDs) == 0 {
			roleIDs = defaultRoleIDs
		}
		for _, api := range page.APIs {
			for _, roleID := range roleIDs {
				key := roleID + "::" + api.Path + "::" + api.Method
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				policies = append(policies, CasbinPolicyItem{
					RoleID: roleID,
					Path:   api.Path,
					Method: api.Method,
				})
			}
		}
	}

	sort.SliceStable(policies, func(i, j int) bool {
		if policies[i].RoleID != policies[j].RoleID {
			return policies[i].RoleID < policies[j].RoleID
		}
		if policies[i].Path != policies[j].Path {
			return policies[i].Path < policies[j].Path
		}
		return policies[i].Method < policies[j].Method
	})
	return policies
}

func sortedPages() []PageSpec {
	pages := make([]PageSpec, len(pageRegistry))
	copy(pages, pageRegistry)
	sort.SliceStable(pages, func(i, j int) bool {
		if pages[i].Area != pages[j].Area {
			return pages[i].Area < pages[j].Area
		}
		if pages[i].Menu.Sort != pages[j].Menu.Sort {
			return pages[i].Menu.Sort < pages[j].Menu.Sort
		}
		return pages[i].Menu.Name < pages[j].Menu.Name
	})
	return pages
}

func areaParentMenuName(area string) string {
	switch area {
	case AreaSimulation:
		return "labSimulation"
	case AreaComponentDemo:
		return "labComponentDemo"
	case AreaReusable:
		return "labReusable"
	default:
		panic(fmt.Sprintf("unknown area %s", area))
	}
}
