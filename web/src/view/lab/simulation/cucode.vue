<template>
    <div class="p-6">
      <!-- 页面标题 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold">代码表管理</h1>
        <p class="text-sm mt-1">系统基础代码数据维护</p>
      </div>
  
      <!-- 分类 Tab 切换 -->
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane name="sa">
          <template #label>
            <span class="flex items-center gap-1.5 py-1">
              <el-icon><Lock /></el-icon>
              数安
            </span>
          </template>
        </el-tab-pane>
        <el-tab-pane name="xa">
          <template #label>
            <span class="flex items-center gap-1.5 py-1">
              <el-icon><Key /></el-icon>
              信安
            </span>
          </template>
        </el-tab-pane>
        <el-tab-pane name="wa">
          <template #label>
            <span class="flex items-center gap-1.5 py-1">
              <el-icon><Monitor /></el-icon>
              网安
            </span>
          </template>
        </el-tab-pane>
      </el-tabs>
  
      <!-- 主内容区：左右布局 -->
      <div class="flex gap-4 mt-4">
        <!-- 左侧：大类列表 -->
        <div class="w-56 shrink-0">
          <el-card :body-style="{ padding: 0 }">
            <template #header>
              <span class="text-sm font-semibold">代码大类</span>
            </template>
            <ul>
              <li
                v-for="category in currentCategories"
                :key="category.key"
                @click="selectCategory(category)"
                class="flex items-center justify-between px-4 py-3 cursor-pointer border-b last:border-b-0"
                :class="{ 'active-item': selectedCategory?.key === category.key }"
              >
                <div class="flex items-center gap-2">
                  <el-icon class="text-xs"><Document /></el-icon>
                  <span class="text-sm">{{ category.label }}</span>
                </div>
                <el-tag size="small" effect="plain">{{ getItemCount(category.key) }}</el-tag>
              </li>
            </ul>
          </el-card>
        </div>
  
        <!-- 右侧：子类数据表格 -->
        <div class="flex-1 min-w-0">
          <el-card>
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="font-semibold">
                    {{ selectedCategory ? selectedCategory.label : '请选择代码大类' }}
                  </span>
                  <span v-if="selectedCategory" class="text-xs">
                    ({{ selectedCategory.field }})
                  </span>
                  <el-divider direction="vertical" />
                  <span class="text-xs">2026-01-01 01:01:01</span>
                </div>
                <div v-if="selectedCategory" class="flex items-center gap-2">
                  <el-input
                    v-model="searchKeyword"
                    placeholder="搜索名称..."
                    size="small"
                    clearable
                    class="w-48"
                  >
                    <template #prefix><el-icon><Search /></el-icon></template>
                  </el-input>
                  <el-button type="primary" size="small" :icon="Plus">新增</el-button>
                </div>
              </div>
            </template>
  
            <!-- 空状态 -->
            <el-empty v-if="!selectedCategory" description="请在左侧选择代码大类" />
  
            <!-- 表格 -->
            <template v-else>
              <el-table
                :data="filteredSubItems"
                stripe
                v-loading="tableLoading"
                class="w-full"
                empty-text="暂无数据"
              >
                <el-table-column type="index" label="#" width="55" align="center" />
                <el-table-column prop="id" label="ID" width="80" align="center" />
                <el-table-column prop="mc" label="名称" min-width="140" />
                <el-table-column
                  v-if="selectedCategory.key === 'fwnr'"
                  prop="fl"
                  label="父类"
                  width="120"
                >
                  <template #default="{ row }">{{ row.fl || '-' }}</template>
                </el-table-column>
                <el-table-column prop="bz" label="备注" min-width="180">
                  <template #default="{ row }">{{ row.bz || '-' }}</template>
                </el-table-column>
                <el-table-column prop="sfyx" label="是否有效" width="100" align="center">
                  <template #default="{ row }">
                    <el-switch v-model="row.sfyx" size="small" disabled />
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="130" align="center" fixed="right">
                  <template #default>
                    <el-button type="primary" link size="small">编辑</el-button>
                    <el-divider direction="vertical" />
                    <el-button type="danger" link size="small">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
  
              <!-- 分页 -->
              <div class="mt-4 flex justify-end">
                <el-pagination
                  v-model:current-page="currentPage"
                  v-model:page-size="pageSize"
                  :total="filteredSubItems.length"
                  :page-sizes="[10, 20, 50]"
                  layout="total, sizes, prev, pager, next"
                  small
                  background
                />
              </div>
            </template>
          </el-card>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, computed, watch } from 'vue'
  import { Lock, Key, Monitor, Document, Search, Plus } from '@element-plus/icons-vue'
  
  const activeTab = ref('sa')
  
  const categoryMap = {
    sa: [
      { key: 'zjlx', label: '证件类型', field: 'zjlx' },
      { key: 'dwsx', label: '单位属性', field: 'dwsx' },
      { key: 'fwnr', label: '服务内容', field: 'fwnr' },
    ],
    xa: [
      { key: 'jrfs', label: '接入方式', field: 'jrfs' },
      { key: 'dllx', label: '代理类型', field: 'dllx' },
      { key: 'fwnr', label: '服务内容', field: 'fwnr' },
    ],
    wa: [
      { key: 'jfxz', label: '机房性质', field: 'jfxz' },
      { key: 'gzlx', label: '监测/过滤规则类型', field: 'gzlx' },
      { key: 'wfwgqk', label: '违法违规情况', field: 'wfwgqk' },
      { key: 'fwnr', label: '服务内容', field: 'fwnr' },
    ],
  }
  
  const currentCategories = computed(() => categoryMap[activeTab.value] || [])
  
  const mockData = {
    zjlx: [
      { id: '01', mc: '居民身份证', bz: '中华人民共和国居民身份证', sfyx: true },
      { id: '02', mc: '护照', bz: '中华人民共和国护照', sfyx: true },
      { id: '03', mc: '港澳居民来往内地通行证', bz: '俗称"回乡证"', sfyx: true },
      { id: '04', mc: '台湾居民来往大陆通行证', bz: '俗称"台胞证"', sfyx: true },
      { id: '05', mc: '外国人永久居留证', bz: '', sfyx: false },
    ],
    dwsx: [
      { id: '01', mc: '国有企业', bz: '', sfyx: true },
      { id: '02', mc: '民营企业', bz: '', sfyx: true },
      { id: '03', mc: '外资企业', bz: '', sfyx: true },
      { id: '04', mc: '政府机关', bz: '', sfyx: true },
      { id: '05', mc: '事业单位', bz: '', sfyx: true },
      { id: '06', mc: '社会团体', bz: '', sfyx: false },
    ],
    jrfs: [
      { id: '01', mc: '专线接入', bz: 'DDN/光纤专线', sfyx: true },
      { id: '02', mc: '宽带接入', bz: 'ADSL/FTTB', sfyx: true },
      { id: '03', mc: '无线接入', bz: '4G/5G/WiFi', sfyx: true },
      { id: '04', mc: '拨号接入', bz: '传统PSTN', sfyx: false },
    ],
    dllx: [
      { id: '01', mc: 'HTTP代理', bz: '', sfyx: true },
      { id: '02', mc: 'SOCKS代理', bz: '', sfyx: true },
      { id: '03', mc: '透明代理', bz: '', sfyx: true },
      { id: '04', mc: '反向代理', bz: '', sfyx: true },
    ],
    jfxz: [
      { id: '01', mc: '自建机房', bz: '企业自建自用', sfyx: true },
      { id: '02', mc: '托管机房', bz: '托管于IDC', sfyx: true },
      { id: '03', mc: '云机房', bz: '公有云/私有云', sfyx: true },
    ],
    gzlx: [
      { id: '01', mc: 'ICP备案过滤', bz: '', sfyx: true },
      { id: '02', mc: '关键词过滤', bz: '', sfyx: true },
      { id: '03', mc: 'IP黑名单', bz: '', sfyx: true },
      { id: '04', mc: '域名监测', bz: '', sfyx: true },
      { id: '05', mc: '内容合规检测', bz: '', sfyx: true },
    ],
    wfwgqk: [
      { id: '01', mc: '传播违法信息', bz: '', sfyx: true },
      { id: '02', mc: '未备案经营', bz: '', sfyx: true },
      { id: '03', mc: '侵犯知识产权', bz: '', sfyx: true },
      { id: '04', mc: '网络诈骗', bz: '', sfyx: true },
      { id: '05', mc: '散布谣言', bz: '', sfyx: false },
    ],
    fwnr: [
      { id: '0101', mc: '互联网接入服务', fl: '基础服务', bz: '', sfyx: true },
      { id: '0102', mc: '互联网数据中心', fl: '基础服务', bz: 'IDC', sfyx: true },
      { id: '0201', mc: '即时通讯服务', fl: '增值服务', bz: '', sfyx: true },
      { id: '0202', mc: '网络存储服务', fl: '增值服务', bz: '', sfyx: true },
      { id: '0203', mc: '信息发布服务', fl: '增值服务', bz: '', sfyx: true },
      { id: '0301', mc: '互联网新闻信息服务', fl: '内容服务', bz: '', sfyx: false },
    ],
  }
  
  const selectedCategory = ref(null)
  const searchKeyword = ref('')
  const tableLoading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  
  function selectCategory(category) {
    selectedCategory.value = category
    searchKeyword.value = ''
    currentPage.value = 1
    tableLoading.value = true
    setTimeout(() => { tableLoading.value = false }, 300)
  }
  
  function handleTabChange() {
    selectedCategory.value = null
    searchKeyword.value = ''
  }
  
  function getItemCount(key) {
    return mockData[key]?.length || 0
  }
  
  const filteredSubItems = computed(() => {
    if (!selectedCategory.value) return []
    const items = mockData[selectedCategory.value.key] || []
    if (!searchKeyword.value) return items
    return items.filter(item =>
      item.mc.includes(searchKeyword.value) || (item.bz && item.bz.includes(searchKeyword.value))
    )
  })
  
  watch(activeTab, () => {
    selectedCategory.value = null
  })
  </script>
  
  <style scoped>
  /* 仅用 Element Plus CSS 变量标记选中态，不硬编码颜色 */
  .active-item {
    background-color: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
    border-left: 2px solid var(--el-color-primary);
  }
  </style>