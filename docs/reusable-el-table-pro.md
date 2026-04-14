# ElTablePro

## 模块归类

- 分区：复用组件
- 类型：前端通用表格增强组件
- 收录位置：`web/src/components/ElTablePro`

## 背景与目标

`ElTablePro` 解决的不是“普通列表页能不能显示数据”，而是复杂后台表格在多个项目里反复出现的这些问题：

- 单元格类型多，页面里充满重复格式化逻辑
- 同一张表既有状态、进度、流量，也有 JSON、链接、开关等异构字段
- 需要多级表头、展开行、详情查看、操作列等复合交互
- 列表页分页请求、搜索参数拼装、返回结果转换逻辑重复

在 `admin-lab` 中收录它，是为了先沉淀一套可迁移的复杂表格底座，后续可以按业务继续裁剪。

## 当前收录范围

已收录文件：

- `web/src/components/ElTablePro/ElTablePro.vue`
- `web/src/components/ElTablePro/renderers/*`
- `web/src/components/ElTablePro/useElTablePro/useElTablePro.js`
- `web/src/components/ElTablePro/Example.vue`
- `web/src/components/ElTablePro/index.js`
- `web/src/hooks/useElTablePro.js`

配套展示入口：

- 前端页面：`实验室 / 复用组件 / ElTablePro`

## 组件能力

当前组件已覆盖：

- 多级表头递归渲染
- 常见后台字段渲染器
- 顶部工具栏插槽
- 展开行插槽
- 固定操作列与详情查看
- 分页器绑定
- Hook 方式的数据加载封装

典型渲染类型包括：

- `status`
- `progress`
- `compound`
- `size`
- `ip`
- `protocol`
- `port`
- `tag`
- `icon`
- `switch`
- `unit`
- `url`
- `json`
- `datetime`

## 最小接入方式

推荐接入方式：

1. 页面维护查询表单
2. 使用 `useElTablePro` 统一处理分页和请求
3. 使用 `columns` 描述列结构和渲染类型
4. 按需补充 `toolbar / expand / actions` 插槽

## 迁移到内网时建议同步的内容

建议至少同步这些文件：

- `web/src/components/ElTablePro`
- `web/src/hooks/useElTablePro.js`
- 使用该组件的业务示例页

如果内网项目目录规范不同，可保留组件目录不变，只调整导出路径。

## 当前边界

这次在 `admin-lab` 里补齐的是：

- 组件目录收录
- Hook 导出路径
- 实验室展示页
- 菜单路由入口
- 迁移说明文档

当前还没有做的是：

- 抽成与现有 `lab/table-pro.vue` 统一协议的一套兼容层
- 提供真实后端分页接口示例
- 对全部渲染器补齐单测或独立文档页

如果后续继续沉淀，建议下一步补：

- 列配置协议文档
- 渲染器扩展规范
- 与真实接口联动的分页示例
