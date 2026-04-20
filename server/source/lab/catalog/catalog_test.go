package catalog

import "testing"

func TestDefaultCatalogContainsLabRootsAndPages(t *testing.T) {
	menus := BuildMenus()
	if len(menus) == 0 {
		t.Fatal("expected menus from catalog")
	}

	required := []string{
		"lab",
		"labSimulation",
		"labComponentDemo",
		"labReusable",
		"labSimulationOverview",
		"labReusableExcelIO",
		"labReusableBizLog",
	}
	for _, menuName := range required {
		if !menuExists(menus, menuName) {
			t.Fatalf("expected menu %s", menuName)
		}
	}
}

func TestDefaultCatalogBuildsApisAndPolicies(t *testing.T) {
	apis := BuildAPIs()
	if len(apis) == 0 {
		t.Fatal("expected apis from catalog")
	}
	if !apiExists(apis, "/excelIO/import", "POST") {
		t.Fatal("expected /excelIO/import POST api")
	}
	if !apiExists(apis, "/bizLog/list", "GET") {
		t.Fatal("expected /bizLog/list GET api")
	}

	policies := BuildCasbinPolicies()
	if len(policies) == 0 {
		t.Fatal("expected casbin policies from catalog")
	}
	if !policyExists(policies, "888", "/excelIO/import", "POST") {
		t.Fatal("expected 888 policy for /excelIO/import POST")
	}
	if !policyExists(policies, "9528", "/bizLog/list", "GET") {
		t.Fatal("expected 9528 policy for /bizLog/list GET")
	}
}

func menuExists(menus []MenuItem, name string) bool {
	for _, menu := range menus {
		if menu.Name == name {
			return true
		}
	}
	return false
}

func apiExists(apis []APIItem, path, method string) bool {
	for _, api := range apis {
		if api.Path == path && api.Method == method {
			return true
		}
	}
	return false
}

func policyExists(policies []CasbinPolicyItem, role, path, method string) bool {
	for _, policy := range policies {
		if policy.RoleID == role && policy.Path == path && policy.Method == method {
			return true
		}
	}
	return false
}
