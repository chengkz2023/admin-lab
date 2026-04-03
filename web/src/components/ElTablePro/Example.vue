<template>
  <div class="example-wrap">
    <div class="summary-grid">
      <div class="summary-card">
        <div class="summary-label">当前结果</div>
        <div class="summary-value">{{ total }}</div>
        <div class="summary-tip">已按顶部关键字即时过滤</div>
      </div>
      <div class="summary-card">
        <div class="summary-label">当前分页</div>
        <div class="summary-value">{{ currentPage }}/{{ pageCount }}</div>
        <div class="summary-tip">分页由 `v-model:current-page` 驱动</div>
      </div>
      <div class="summary-card">
        <div class="summary-label">选中行数</div>
        <div class="summary-value">{{ selectedRows.length }}</div>
        <div class="summary-tip">适合接批量操作栏</div>
      </div>
    </div>

    <ElTablePro
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :columns="columns"
      :data="pagedData"
      :total="total"
      :loading="loading"
      show-detail
      :actions-config="{
        width: 150,
      }"
      border
      highlight-current-row
      @page-change="handlePageChange"
      @sort-change="handleSortChange"
      @selection-change="handleSelectionChange"
    >
      <template #toolbar>
        <div class="toolbar-wrap">
          <div>
            <h3 class="toolbar-title">服务监控列表示例</h3>
            <p class="toolbar-subtitle">
              演示工具栏插槽、自定义列插槽、展开行、详情抽屉、开关渲染器与分页联动。
            </p>
          </div>

          <div class="toolbar-actions">
            <el-input
              v-model="searchQuery"
              placeholder="按服务名、负责人、区域筛选"
              clearable
              class="search-input"
              @keyup.enter="handleToolbarSearch"
              @clear="handleToolbarSearch"
            />
            <el-button type="primary" :icon="Search" @click="handleToolbarSearch">
              搜索
            </el-button>
            <el-button :icon="Refresh" circle @click="handleToolbarRefresh" />
            <el-button
              type="warning"
              plain
              :disabled="selectedRows.length === 0"
              @click="handleBatchAction"
            >
              批量巡检
            </el-button>
          </div>
        </div>
      </template>

      <template #expand="{ row }">
        <div class="expand-wrap">
          <CompoundCell
            :row="row"
            :column="{
              items: [
                { prop: 'sourceIp', label: '源 IP', type: 'ip' },
                { prop: 'region', label: '区域', format: (val) => `区域：${val || '-'}` },
                {
                  prop: 'owner',
                  label: '负责人',
                  type: 'tag',
                  config: {
                    '张三': { type: 'primary' },
                    '李四': { type: 'success' },
                    '王五': { type: 'warning' }
                  }
                },
                { prop: 'fullRequest', label: '完整请求', type: 'json', previewText: '查看完整请求体' }
              ],
              config: {
                labelWidth: '100px',
              }
            }"
          />
        </div>
      </template>

      <template #serverName="{ row }">
        <div class="server-link" @click="viewDetails(row)">
          <el-icon :size="16" class="server-link__icon">
            <Platform />
          </el-icon>
          <div>
            <div class="server-link__title">{{ row.serverName }}</div>
            <div class="server-link__meta">{{ row.region || '未分区' }} / {{ row.owner || '未分配' }}</div>
          </div>
        </div>
      </template>

      <template #actions="{ row }">
        <div class="actions-wrap">
          <el-button type="primary" link size="small" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-popconfirm title="确定删除此行吗？" @confirm="handleDelete(row)">
            <template #reference>
              <el-button type="danger" link size="small">
                删除
              </el-button>
            </template>
          </el-popconfirm>
        </div>
      </template>
    </ElTablePro>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import { ElMessage, ElMessageBox, ElPopconfirm } from 'element-plus';
import { Platform, Search, Refresh } from '@element-plus/icons-vue';
import ElTablePro from './ElTablePro.vue';
import CompoundCell from './renderers/compoundCell.vue';

const loading = ref(false);
const searchQuery = ref('');
const selectedRows = ref([]);
const currentPage = ref(1);
const pageSize = ref(10);

const rawTableData = ref([
  {
    id: 1,
    serverName: 'DB-Master-1',
    status: 'up',
    memoryUsage: 48.6,
    packetLoss: 79.8,
    uploadRate: 183.6835,
    downloadSize: 346882529206.6157,
    lastSeen: '2025-06-07 12:00:00',
    nodeType: 'edge',
    serviceType: 'database',
    requestCount: 134322349251,
    region: '华北',
    sourceIp: '133.74.144.187',
    logFile: 'https://example.com/logs/app-1.log',
    owner: '李四',
    shortRequest: '{"page":1,"limit":10}',
    fullRequest: {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Request-ID': 'req-1'
      },
      body: {
        serverName: 'DB-Master-1',
        action: 'restart'
      }
    },
    protocol: 6,
    port: '18080',
    envCode: 'prod',
    featureEnabled: 1,
    toggleLocked: false,
    toggleShouldFail: false,
    toggleNote: '关闭时会二次确认，异步保存后保留状态',
  },
  {
    id: 5,
    serverName: 'ETL-Worker-5',
    status: 'down',
    memoryUsage: 96.7,
    packetLoss: 0.56,
    uploadRate: 175775860.7806006,
    downloadSize: 103568046748.29298,
    lastSeen: 1750011156636.395,
    nodeType: 'edge',
    serviceType: 'database',
    requestCount: 508813.3154338964,
    region: '华南',
    sourceIp: '157.14.151.61',
    logFile: 'https://example.com/logs/app-5.log',
    owner: '王五',
    shortRequest: '{"page":1,"limit":10}',
    fullRequest: {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Request-ID': 'req-5'
      },
      body: {
        serverName: 'ETL-Worker-5',
        action: 'restart'
      }
    },
    protocol: 17,
    port: 3306,
    envCode: 'staging',
    featureEnabled: 0,
    toggleLocked: true,
    toggleShouldFail: false,
    toggleNote: '该行被 disabled 规则锁定，不允许切换',
  },
  {
    id: 9,
    serverName: 'Cache-Replica-9',
    status: 'pending',
    memoryUsage: 34.2,
    packetLoss: 1.35,
    uploadRate: 128576.4,
    downloadSize: 2458291201,
    lastSeen: '2025-06-07 13:45:10',
    nodeType: 'core',
    serviceType: 'cache',
    requestCount: 36210500,
    region: '华东',
    sourceIp: '10.23.9.88',
    logFile: 'https://example.com/logs/cache-9.log',
    owner: '张三',
    shortRequest: '{"sync":false,"cache":"warmup"}',
    fullRequest: {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-Request-ID': 'req-9'
      },
      body: {
        serverName: 'Cache-Replica-9',
        action: 'flush'
      }
    },
    protocol: 6,
    port: 6379,
    envCode: 'prod',
    featureEnabled: 0,
    toggleLocked: false,
    toggleShouldFail: true,
    toggleNote: '模拟接口失败，切换后会自动回滚',
  },
  {
    id: 12,
    serverName: 'Gateway-Core-12',
    status: 'restarting',
    memoryUsage: 81.2,
    packetLoss: 12.4,
    uploadRate: 32684512,
    downloadSize: 9876543210,
    lastSeen: '2025-06-07 14:10:00',
    nodeType: 'core',
    serviceType: 'gateway',
    requestCount: 872310000,
    region: '华中',
    sourceIp: '172.16.30.12',
    logFile: 'https://example.com/logs/gateway-12.log',
    owner: '李四',
    shortRequest: '{"path":"/open/api","qps":1024}',
    fullRequest: {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'X-Request-ID': 'req-12'
      },
      body: {
        serverName: 'Gateway-Core-12',
        action: 'reload'
      }
    },
    protocol: 6,
    port: 8443,
    envCode: 'gray',
    featureEnabled: 1,
    toggleLocked: false,
    toggleShouldFail: false,
    toggleNote: '最近一次切换：已开启',
  }
]);

const wait = (ms) => new Promise((resolve) => {
  window.setTimeout(resolve, ms);
});

const filteredData = computed(() => {
  const keyword = searchQuery.value.trim().toLowerCase();
  if (!keyword) {
    return rawTableData.value;
  }
  return rawTableData.value.filter((row) => {
    return [row.serverName, row.owner, row.region, row.sourceIp]
      .filter(Boolean)
      .some((field) => String(field).toLowerCase().includes(keyword));
  });
});

const total = computed(() => filteredData.value.length);

const pageCount = computed(() => {
  return Math.max(1, Math.ceil(total.value / pageSize.value));
});

const pagedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return filteredData.value.slice(start, start + pageSize.value);
});

const resetCurrentPageIfNeeded = () => {
  if (currentPage.value > pageCount.value) {
    currentPage.value = pageCount.value;
  }
};

const columns = ref([
  { type: 'expand', width: 55, fixed: 'left' },
  { type: 'selection', width: 55, fixed: 'left' },
  { type: 'index', label: '序号', width: 70, fixed: 'left' },
  { type: 'id', prop: 'id', label: 'ID', width: 60, fixed: 'left' },
  { type: 'datetime', prop: 'lastSeen', label: '最后心跳', width: 160, config: { format: 'YYYY/MM/DD HH:mm' } },
  { prop: 'serverName', label: '服务名', minWidth: 180, fixed: 'left' },
  { prop: 'owner', label: '负责人', width: 100, formatter: (_row, _column, value) => value || '-' },
  {
    label: '服务状态',
    children: [
      {
        type: 'status',
        prop: 'status',
        label: '状态',
        width: 100,
        config: {
          up: { label: '运行中', type: 'success' },
          down: { label: '异常', type: 'danger' },
          pending: { label: '待同步', type: 'warning' },
          restarting: { label: '重启中', type: 'warning', isPing: true }
        }
      },
      {
        type: 'progress',
        prop: 'memoryUsage',
        label: '内存占用',
        width: 120,
        config: { valueType: 'decimal', thresholds: { warning: 60, danger: 85 } }
      },
      {
        type: 'progress',
        prop: 'memoryUsage',
        label: '圆形占比',
        width: 140,
        align: 'center',
        config: { type: 'circle' }
      },
    ]
  },
  {
    type: 'compound',
    prop: 'networkMetrics',
    label: '网络指标',
    minWidth: 180,
    sortBy: 'uploadRate',
    items: [
      { type: 'signal', prop: 'packetLoss', label: '丢包率', align: 'right', config: { inverse: true } },
      { type: 'size', prop: 'uploadRate', label: '上行流量', align: 'center' },
    ],
  },
  { type: 'size', prop: 'uploadRate', label: '上行流量', width: 120 },
  { type: 'ip', prop: 'sourceIp', label: '源 IP', width: 150, align: 'center' },
  { type: 'protocol', prop: 'protocol', label: '协议', width: 100 },
  { type: 'port', prop: 'port', label: '端口', width: 100 },
  {
    type: 'tag',
    prop: 'nodeType',
    label: '节点类型',
    width: 100,
    config: {
      core: { label: '核心', type: 'danger' },
      edge: { label: '边缘', type: 'success' },
    },
  },
  {
    type: 'options',
    prop: 'envCode',
    label: '环境',
    width: 100,
    options: [
      { label: '生产', value: 'prod', type: 'danger' },
      { label: '预发', value: 'staging', type: 'warning' },
      { label: '灰度', value: 'gray', type: 'primary' }
    ]
  },
  { type: 'icon', prop: 'serviceType', label: '服务类型', width: 110 },
  {
    type: 'switch',
    prop: 'featureEnabled',
    label: '功能开关',
    width: 120,
    align: 'center',
    config: {
      activeValue: 1,
      inactiveValue: 0,
      inlinePrompt: true,
      activeText: '开',
      inactiveText: '关',
      disabled: ({ row }) => row.toggleLocked,
      beforeChange: async ({ row, nextValue }) => {
        if (nextValue !== 0) return true;

        try {
          await ElMessageBox.confirm(
            `确认关闭 ${row.serverName} 的功能开关吗？`,
            '关闭确认',
            {
              type: 'warning',
              confirmButtonText: '确认关闭',
              cancelButtonText: '取消',
            }
          );
          return true;
        } catch (_error) {
          return false;
        }
      },
      onChange: async ({ row, nextValue }) => {
        await wait(700);

        if (row.toggleShouldFail) {
          ElMessage.error(`${row.serverName} 模拟保存失败，已自动回滚`);
          throw new Error('mock failure');
        }

        row.toggleNote = nextValue === 1
          ? '最近一次切换：已开启'
          : '最近一次切换：已关闭';
        ElMessage.success(`${row.serverName} 开关已${nextValue === 1 ? '开启' : '关闭'}`);
        return true;
      },
    },
  },
  { prop: 'toggleNote', label: '开关说明', minWidth: 220 },
  {
    type: 'unit',
    prop: 'requestCount',
    label: '请求量',
    align: 'right',
    width: 100,
    config: {
      unitBase: 1000,
      units: ['', 'K', 'M', 'B'],
      decimals: 1,
    },
  },
  { type: 'url', prop: 'logFile', label: '日志地址', minWidth: 160 },
  { type: 'json', prop: 'shortRequest', label: '请求参数', width: 150, previewText: '查看参数' },
]);

const handlePageChange = () => {
  resetCurrentPageIfNeeded();
};

const handleSortChange = ({ prop, order }) => {
  ElMessage.info(`排序事件已触发：列=${prop || '-'}，顺序=${order || '默认'}`);
};

const handleSelectionChange = (val) => {
  selectedRows.value = val || [];
};

const handleToolbarSearch = () => {
  currentPage.value = 1;
  ElMessage.info(`已按关键字筛选：${searchQuery.value || '全部'}`);
};

const handleToolbarRefresh = async () => {
  loading.value = true;
  await wait(400);
  loading.value = false;
  resetCurrentPageIfNeeded();
  ElMessage.success('演示数据已刷新');
};

const handleBatchAction = () => {
  ElMessage.success(`已提交 ${selectedRows.value.length} 条服务进行批量巡检`);
};

const viewDetails = (row) => ElMessage.info(`查看详情：${row.serverName}`);
const handleEdit = (row) => ElMessage.success(`编辑：${row.serverName}`);
const handleDelete = (row) => {
  rawTableData.value = rawTableData.value.filter((item) => item.id !== row.id);
  selectedRows.value = selectedRows.value.filter((item) => item.id !== row.id);
  resetCurrentPageIfNeeded();
  ElMessage.success(`已删除：${row.serverName}`);
};
</script>

<style scoped>
.example-wrap {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.summary-card {
  padding: 16px 18px;
  border: 1px solid #dbeafe;
  border-radius: 14px;
  background: linear-gradient(135deg, #eff6ff 0%, #f8fafc 100%);
}

.summary-label {
  font-size: 12px;
  color: #2563eb;
  font-weight: 700;
  letter-spacing: 0.06em;
}

.summary-value {
  margin-top: 8px;
  font-size: 28px;
  line-height: 1.1;
  font-weight: 700;
  color: #0f172a;
}

.summary-tip {
  margin-top: 6px;
  font-size: 13px;
  color: #475569;
}

.toolbar-wrap {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.toolbar-title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}

.toolbar-subtitle {
  margin: 6px 0 0;
  font-size: 13px;
  color: #64748b;
}

.toolbar-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.search-input {
  width: 280px;
}

.expand-wrap {
  padding: 8px 4px;
}

.server-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.server-link__icon {
  color: #2563eb;
}

.server-link__title {
  font-weight: 700;
  color: #0f172a;
}

.server-link__meta {
  font-size: 12px;
  color: #64748b;
}

.server-link:hover .server-link__title {
  color: #2563eb;
  text-decoration: underline;
}

.actions-wrap {
  display: inline-flex;
  align-items: center;
}

::v-deep(.el-table .el-table__header-wrapper th) {
  @apply p-2 !important;
}

::v-deep(.el-table .el-table__row td) {
  @apply p-2 !important;
}

@media (max-width: 900px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }

  .search-input {
    width: 100%;
  }
}
</style>
