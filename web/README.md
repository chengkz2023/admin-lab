# Admin Scaffold

基于 Vue 3 + Vite + Element Plus + Pinia 的管理后台前端脚手架，由开源项目精简而成，便于在此基础上开发自己的业务。

## 技术栈

- **Vue 3** + **Vite** + **Vue Router**
- **Element Plus** + **UnoCSS**
- **Pinia**（状态管理）
- **Axios**（请求封装）

## 本地运行

```bash
npm install
npm run dev
```

默认使用 **静态前端菜单** 和 **Mock 登录**，无需后端即可登录（任意用户名、密码即可进入）。

## 配置说明

在 `src/core/config.js` 中可修改：

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `appName` | 应用名称（登录页、标题等） | `'Admin Scaffold'` |
| `useStaticMenu` | 是否使用前端静态菜单（不请求后端 getMenu） | `true` |
| `useMockLogin` | 是否启用 Mock 登录（无后端时使用） | `true` |

- 接入真实后端时：将 `useStaticMenu`、`useMockLogin` 设为 `false`，并配置 `.env.development` 等环境变量中的接口地址。

## 目录与扩展方式

- **静态菜单**：在 `src/router/staticRoutes.js` 的 `staticMenus` 中增删改菜单项，每项需包含 `path`、`name`、`meta.title`、`component`（如 `'view/xxx/yyy.vue'`）。
- **页面视图**：在 `src/view/` 下按模块建目录与 `.vue` 文件，再在静态菜单或后端菜单中注册。
- **接口**：在 `src/api/` 下按模块封装，请求基址在 `src/utils/request.js` 中配置。
- **布局**：主布局在 `src/view/layout/`（侧栏、顶栏、标签页、设置等），可按需调整。

## 构建与预览

```bash
npm run build
npm run preview
```

## 说明

- 当前为纯净脚手架版本，已移除原项目的业务页面与品牌信息，仅保留登录、布局、首页与示例页。
- 需要动态菜单与鉴权时，关闭 `useStaticMenu` / `useMockLogin` 并接入自己的后端接口即可。
