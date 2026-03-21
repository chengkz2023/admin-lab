/**
 * 前端静态菜单（脚手架默认）
 * 当 config.useStaticMenu 为 true 时使用，无需后端 getMenu 接口
 */
import { asyncRouterHandle } from '@/utils/asyncRouter'

const staticMenus = [
  {
    path: 'admin',
    name: 'superAdmin',
    meta: { title: '超级管理员', icon: 'user' },
    component: 'view/superAdmin/index.vue',
    children: [
      { path: 'authority', name: 'authority', meta: { title: '角色管理', icon: 'avatar' }, component: 'view/superAdmin/authority/authority.vue' },
      { path: 'menu', name: 'menu', meta: { title: '菜单管理', icon: 'tickets', keepAlive: true }, component: 'view/superAdmin/menu/menu.vue' },
      { path: 'api', name: 'api', meta: { title: 'API管理', icon: 'platform', keepAlive: true }, component: 'view/superAdmin/api/api.vue' },
      { path: 'user', name: 'user', meta: { title: '用户管理', icon: 'coordinate' }, component: 'view/superAdmin/user/user.vue' },
      { path: 'dictionary', name: 'dictionary', meta: { title: '字典管理', icon: 'notebook' }, component: 'view/superAdmin/dictionary/sysDictionary.vue' },
      { path: 'operation', name: 'operation', meta: { title: '操作历史', icon: 'pie-chart' }, component: 'view/superAdmin/operation/sysOperationRecord.vue' },
      { path: 'sysParams', name: 'sysParams', meta: { title: '参数管理', icon: 'compass' }, component: 'view/superAdmin/params/sysParams.vue' }
    ]
  },
  {
    path: 'lab',
    name: 'lab',
    meta: { title: '实验室', icon: 'data-analysis' },
    component: 'view/lab/index.vue',
    children: [
      {
        path: 'simulation',
        name: 'labSimulation',
        meta: { title: '需求仿真', icon: 'document' },
        component: 'view/routerHolder.vue',
        children: [
          { path: 'overview', name: 'labSimulationOverview', meta: { title: '概览', icon: 'tickets' }, component: 'view/lab/simulation/overview.vue' },
          { path: 'base-data-io', name: 'labSimulationBaseDataIO', meta: { title: '基础数据导入导出仿真', icon: 'document-copy' }, component: 'view/lab/simulation/base-data-io.vue' }
        ]
      },
      {
        path: 'component-demo',
        name: 'labComponentDemo',
        meta: { title: '组件示例', icon: 'magic-stick' },
        component: 'view/routerHolder.vue',
        children: [
          { path: 'overview', name: 'labComponentDemoOverview', meta: { title: '概览', icon: 'tickets' }, component: 'view/lab/component-demo/overview.vue' }
        ]
      },
      {
        path: 'reusable',
        name: 'labReusable',
        meta: { title: '复用组件', icon: 'files' },
        component: 'view/routerHolder.vue',
        children: [
          { path: 'overview', name: 'labReusableOverview', meta: { title: '概览', icon: 'tickets' }, component: 'view/lab/reusable/overview.vue' },
          { path: 'excel-io', name: 'labReusableExcelIO', meta: { title: 'Excel 实验面板', icon: 'document-copy' }, component: 'view/lab/reusable/excel-io.vue' },
          { path: 'list-query-bar', name: 'labReusableListQueryBar', meta: { title: '列表查询栏', icon: 'search' }, component: 'view/lab/reusable/list-query-bar.vue' },
          { path: 'reliable-upload', name: 'labReusableReliableUpload', meta: { title: '可靠上报框架', icon: 'upload-filled' }, component: 'view/lab/reusable/reliable-upload.vue' }
        ]
      }
    ]
  }
]

export function getStaticRouter() {
  const baseRouter = [
    {
      path: '/layout',
      name: 'layout',
      component: 'view/layout/index.vue',
      meta: { title: '底层layout' },
      children: [...staticMenus]
    }
  ]
  asyncRouterHandle(baseRouter)
  return baseRouter
}

export { staticMenus }
