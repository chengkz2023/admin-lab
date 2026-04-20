# Lab Source 种子说明

`server/source/lab` 现在统一负责所有 Lab 相关种子数据：

- Lab 菜单
- Lab API
- Lab Casbin 权限

## 如何新增一个 Lab 页面

1. 在 `server/source/lab/catalog/` 下新增一个文件（例如 `page_xxx.go`）。
2. 在 `init()` 中通过 `registerPage(...)` 注册一个 `PageSpec`。
3. 补全菜单信息（`Area`、`Menu.Path`、`Menu.Name`、`Menu.Component`、`Menu.Sort`、`Title`、`Icon`）。
4. 如果该页面有后端接口，补充 `APIs`。

完成后，种子层会自动处理：

- 菜单初始化
- API 初始化
- 默认角色策略（`888`、`8881`、`9528`）

不再需要额外修改 `system/menu.go`、`authorities_menus.go`、`api_*.go`、`casbin_*.go`。
