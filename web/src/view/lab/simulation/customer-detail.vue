<template>
  <div v-loading="loading" class="page-wrap">
    <el-empty v-if="!loading && !detail" description="未加载到仿真数据" />
    <template v-else-if="detail">
      <div class="topbar">
        <div class="breadcrumb">
          <template v-for="(item, index) in detail.breadcrumb" :key="item">
            <span :class="{ cur: index === detail.breadcrumb.length - 1 }">{{ item }}</span>
            <span v-if="index !== detail.breadcrumb.length - 1" class="sep">›</span>
          </template>
        </div>
        <span class="status-badge">{{ detail.status }}</span>
      </div>

      <div class="header-card">
        <div class="avatar">{{ detail.customer.shortName }}</div>
        <div>
          <div class="h1">{{ detail.customer.name }}</div>
          <div class="hmeta">
            <span>客户编号 {{ detail.customer.customerNo }}</span>
            <span class="hdot" />
            <span class="mono">{{ detail.customer.unifiedCode }}</span>
            <span class="hdot" />
            <span class="tag tag-blue">{{ detail.customer.serviceType }}</span>
            <span class="hdot" />
            <span>开通 {{ detail.customer.serviceStart }}</span>
          </div>
        </div>
      </div>

      <div class="metrics-row">
        <div v-for="metric in detail.metrics" :key="metric.label" class="metric">
          <div class="mlabel">{{ metric.label }}</div>
          <div class="mvalue">{{ metric.value }}</div>
          <div class="msub">{{ metric.sub }}</div>
        </div>
      </div>

      <div class="two-col">
        <div class="card" style="margin-bottom: 0">
          <div class="card-head"><span class="card-title">基本信息</span></div>
          <div class="card-body">
            <div v-for="item in detail.basicInfo" :key="item.key" class="kv-row">
              <span class="kv-key">{{ item.key }}</span>
              <span class="kv-val">
                <span v-if="item.type === 'tag-blue'" class="tag tag-blue">{{ item.value }}</span>
                <span v-else :class="{ mono: item.type === 'mono' }">{{ item.value }}</span>
              </span>
            </div>
          </div>
        </div>
        <div class="card" style="margin-bottom: 0">
          <div class="card-head"><span class="card-title">网络信息安全责任人</span></div>
          <div class="card-body">
            <div v-for="item in detail.securityOwner" :key="item.key" class="kv-row">
              <span class="kv-key">{{ item.key }}</span>
              <span class="kv-val" :class="{ mono: item.type === 'mono' }">{{ item.value }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="sec-label">机房 / 节点资源占用</div>
      <div class="card">
        <table class="rooms-table-wrap">
          <thead>
            <tr>
              <th style="width: 200px">机房 / 节点</th>
              <th style="width: 150px">占用机架</th>
              <th>接入信息</th>
              <th style="width: 80px">详情</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="room in detail.roomNodes" :key="room.code">
              <tr class="room-row-main">
                <td>
                  <div class="room-name-cell">{{ room.name }}</div>
                  <div class="room-sub-cell">{{ room.subName }}</div>
                  <div style="margin-top: 5px; display: flex; gap: 4px">
                    <span v-for="tag in room.tags" :key="`${room.code}-${tag}`" :class="['tag', roomTagClass(tag)]">
                      {{ tag }}
                    </span>
                  </div>
                </td>
                <td>
                  <div class="rack-chips">
                    <span v-for="rack in room.racks" :key="`${room.code}-${rack}`" class="rack-chip">{{ rack }}</span>
                  </div>
                </td>
                <td>
                  <div class="conn-list">
                    <div v-for="line in room.connections" :key="`${room.code}-${line.label}`" class="conn-item">
                      <span class="conn-label">{{ line.label }}</span>{{ line.content }}
                    </div>
                  </div>
                </td>
                <td>
                  <button class="expand-btn" @click="toggleRoom(room.code)">
                    {{ isRoomOpen(room.code) ? '▾ 收起' : '▸ 展开' }}
                  </button>
                </td>
              </tr>
              <tr class="room-row-expand">
                <td colspan="4">
                  <div class="room-expand-inner" :class="{ open: isRoomOpen(room.code) }">
                    <div v-for="group in room.details" :key="`${room.code}-${group.label}`" class="room-detail-group">
                      <div class="room-detail-label">{{ group.label }}</div>
                      <div v-for="item in group.items" :key="`${room.code}-${group.label}-${item.key}`" class="room-detail-kv">
                        <span class="room-detail-key">{{ item.key }}</span>
                        <span class="room-detail-val" :class="{ mono: item.type === 'mono' }">{{ item.value }}</span>
                      </div>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <div class="sec-label">客户 IP 段明细</div>
      <div class="card">
        <table class="iptable">
          <colgroup>
            <col style="width: 44px">
            <col style="width: 130px">
            <col style="width: 44px">
            <col style="width: 52px">
            <col style="width: 68px">
            <col style="width: 88px">
            <col style="width: 84px">
            <col>
          </colgroup>
          <thead>
            <tr>
              <th>来源</th>
              <th>IP 段（起始 — 终止）</th>
              <th>掩码</th>
              <th>地址数</th>
              <th>使用方式</th>
              <th>所属机房 / 节点</th>
              <th>分配时间</th>
              <th>NAT 映射</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="segment in visibleIPSegments" :key="segment.range">
              <tr>
                <td><span :class="sourceClass(segment.source)">{{ segment.source }}</span></td>
                <td class="mono td-p">{{ segment.range }}</td>
                <td class="mono">{{ segment.mask }}</td>
                <td>{{ segment.addressNum }}</td>
                <td><span :class="usageClass(segment.usage)">{{ segment.usage }}</span></td>
                <td>{{ segment.nodeCode }}</td>
                <td>{{ segment.allocatedAt }}</td>
                <td>
                  <button
                    v-if="segment.natMappings?.length"
                    class="nat-btn"
                    @click="toggleNat(segment.range)"
                  >
                    {{ isNatOpen(segment.range) ? '▾ 收起' : '▸ 有映射' }}
                  </button>
                  <span v-else class="nat-empty">无</span>
                </td>
              </tr>
              <tr v-if="segment.natMappings?.length" class="nat-expand-row">
                <td colspan="8">
                  <div class="nat-expand-detail" :class="{ open: isNatOpen(segment.range) }">
                    <div v-for="line in segment.natMappings" :key="`${segment.range}-${line.from}`" class="nat-line">
                      <span>{{ line.from }}</span>
                      <span class="nat-arr">→</span>
                      <span>{{ line.to }}</span>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
        <button v-if="hiddenIPCount > 0" class="show-more-btn" @click="toggleHiddenIP">
          {{ showAllIP ? '收起 ▴' : `展示更多 ${hiddenIPCount} 段 ▾` }}
        </button>
      </div>

      <div class="sec-label">应用服务</div>
      <div class="card">
        <table class="app-table">
          <thead>
            <tr>
              <th>服务类型</th>
              <th>许可证 / 备案号</th>
              <th>接入方式</th>
              <th>服务内容</th>
              <th>域名</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in detail.applications" :key="item.domain">
              <td><span :class="['tag', appTagClass(item.serviceType)]">{{ item.serviceType }}</span></td>
              <td class="mono">{{ item.permit }}</td>
              <td>{{ item.accessMode }}</td>
              <td class="td-p">{{ item.content }}</td>
              <td><span class="ip-pill">{{ item.domain }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getCustomerDetailSimulation } from '@/api/customerDetailSimulation'

  defineOptions({
    name: 'LabSimulationCustomerDetail'
  })

  const loading = ref(false)
  const detail = ref(null)
  const openedRooms = ref([])
  const openedNATs = ref([])
  const showAllIP = ref(false)

  const hiddenIPCount = computed(() => {
    return detail.value?.ipSegments?.filter((item) => item.hidden).length || 0
  })

  const visibleIPSegments = computed(() => {
    const all = detail.value?.ipSegments || []
    if (showAllIP.value) {
      return all
    }
    return all.filter((item) => !item.hidden)
  })

  const roomTagClass = (tag) => {
    if (tag === 'IDC') {
      return 'tag-teal'
    }
    if (tag === 'ISP') {
      return 'tag-blue'
    }
    return ''
  }

  const sourceClass = (source) => {
    return source === 'IDC' ? 'src-idc' : 'src-isp'
  }

  const usageClass = (usage) => {
    if (usage === '静态') {
      return 'use-static'
    }
    if (usage === '动态') {
      return 'use-dynamic'
    }
    return 'use-leased'
  }

  const appTagClass = (serviceType) => {
    if (serviceType === '内部应用') {
      return ''
    }
    return 'tag-blue'
  }

  const isRoomOpen = (code) => openedRooms.value.includes(code)

  const isNatOpen = (range) => openedNATs.value.includes(range)

  const toggleRoom = (code) => {
    if (isRoomOpen(code)) {
      openedRooms.value = openedRooms.value.filter((item) => item !== code)
      return
    }
    openedRooms.value = [...openedRooms.value, code]
  }

  const toggleNat = (range) => {
    if (isNatOpen(range)) {
      openedNATs.value = openedNATs.value.filter((item) => item !== range)
      return
    }
    openedNATs.value = [...openedNATs.value, range]
  }

  const toggleHiddenIP = () => {
    showAllIP.value = !showAllIP.value
  }

  const loadDetail = async () => {
    loading.value = true
    try {
      const res = await getCustomerDetailSimulation()
      detail.value = res.data || null
    } catch (error) {
      ElMessage.error('客户详情仿真数据加载失败')
    } finally {
      loading.value = false
    }
  }

  onMounted(async () => {
    await loadDetail()
  })
</script>

<style scoped>
  * {
    box-sizing: border-box;
  }

  .page-wrap {
    background: #f5f7fb;
    padding: 20px;
    border-radius: 14px;
  }

  .topbar {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 20px;
  }

  .breadcrumb {
    font-size: 13px;
    color: #8f98a7;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .cur {
    color: #2b3444;
    font-weight: 500;
  }

  .sep {
    font-size: 11px;
  }

  .status-badge {
    margin-left: auto;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 500;
    background: #e1f5ee;
    color: #0f6e56;
  }

  .header-card {
    background: #fff;
    border: 1px solid #e4e7ee;
    border-radius: 14px;
    padding: 16px 20px;
    margin-bottom: 12px;
    display: flex;
    align-items: flex-start;
    gap: 14px;
  }

  .avatar {
    width: 46px;
    height: 46px;
    border-radius: 10px;
    background: #e6f1fb;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 15px;
    font-weight: 500;
    color: #0c447c;
    flex-shrink: 0;
  }

  .h1 {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 4px;
  }

  .hmeta {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 12px;
    color: #687285;
    flex-wrap: wrap;
  }

  .hdot {
    width: 3px;
    height: 3px;
    border-radius: 50%;
    background: #cad0db;
    flex-shrink: 0;
  }

  .metrics-row {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 10px;
    margin-bottom: 12px;
  }

  .metric {
    background: #eff3fa;
    border-radius: 10px;
    padding: 11px 14px;
  }

  .mlabel {
    font-size: 11px;
    color: #8d97a8;
    margin-bottom: 3px;
  }

  .mvalue {
    font-size: 17px;
    font-weight: 500;
  }

  .msub {
    font-size: 11px;
    color: #8d97a8;
    margin-top: 2px;
  }

  .two-col {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    margin-bottom: 12px;
  }

  .card {
    background: #fff;
    border: 1px solid #e4e7ee;
    border-radius: 14px;
    overflow: hidden;
    margin-bottom: 12px;
  }

  .card:last-child {
    margin-bottom: 0;
  }

  .card-head {
    padding: 11px 16px;
    border-bottom: 1px solid #e4e7ee;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .card-title {
    font-size: 13px;
    font-weight: 500;
  }

  .card-body {
    padding: 14px 16px;
  }

  .kv-row {
    display: flex;
    align-items: baseline;
    padding: 5px 0;
    border-bottom: 1px solid #e4e7ee;
    gap: 8px;
  }

  .kv-row:last-child {
    border-bottom: none;
  }

  .kv-key {
    font-size: 12px;
    color: #8d97a8;
    width: 88px;
    flex-shrink: 0;
  }

  .kv-val {
    font-size: 13px;
    flex: 1;
    color: #2b3444;
  }

  .tag {
    display: inline-block;
    padding: 1px 8px;
    border-radius: 20px;
    font-size: 11px;
    font-weight: 500;
    background: #eff3fa;
    color: #687285;
    border: 1px solid #d8deea;
  }

  .tag-blue {
    background: #e6f1fb;
    color: #185fa5;
    border-color: #b5d4f4;
  }

  .tag-teal {
    background: #e1f5ee;
    color: #0f6e56;
    border-color: #9fe1cb;
  }

  .mono {
    font-family: Menlo, Consolas, Monaco, monospace;
    font-size: 12px;
  }

  .sec-label {
    font-size: 11px;
    font-weight: 500;
    color: #8d97a8;
    letter-spacing: 0.05em;
    padding: 4px 0 8px;
    margin-top: 4px;
  }

  .rooms-table-wrap {
    width: 100%;
    border-collapse: collapse;
  }

  .rooms-table-wrap th {
    font-size: 11px;
    font-weight: 500;
    color: #8d97a8;
    text-align: left;
    padding: 7px 12px;
    border-bottom: 1px solid #e4e7ee;
    background: #f7f9fd;
    white-space: nowrap;
  }

  .rooms-table-wrap td {
    font-size: 12px;
    padding: 0;
    border-bottom: 1px solid #e4e7ee;
    vertical-align: top;
  }

  .rooms-table-wrap tr:last-child > td {
    border-bottom: none;
  }

  .room-row-main td {
    padding: 10px 12px;
    border-bottom: none;
  }

  .room-row-expand td {
    padding: 0;
    border-bottom: 1px solid #e4e7ee;
  }

  .room-name-cell {
    color: #2b3444;
    font-weight: 500;
    font-size: 13px;
  }

  .room-sub-cell {
    font-size: 11px;
    color: #8d97a8;
    margin-top: 2px;
  }

  .rack-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .rack-chip {
    font-family: Menlo, Consolas, Monaco, monospace;
    font-size: 11px;
    font-weight: 500;
    padding: 2px 7px;
    border-radius: 4px;
    background: #e1f5ee;
    color: #0f6e56;
    border: 1px solid #9fe1cb;
  }

  .conn-list {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .conn-item {
    font-size: 12px;
    color: #687285;
  }

  .conn-item .conn-label {
    color: #8d97a8;
    font-size: 11px;
    margin-right: 4px;
  }

  .expand-btn {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 11px;
    color: #8d97a8;
    background: #f7f9fd;
    border: 1px solid #d8deea;
    cursor: pointer;
    font-family: inherit;
    white-space: nowrap;
  }

  .expand-btn:hover {
    color: #2b3444;
    background: #fff;
  }

  .room-expand-inner {
    display: none;
    padding: 10px 12px;
    background: #f7f9fd;
    border-top: 1px solid #e4e7ee;
  }

  .room-expand-inner.open {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }

  .room-detail-label {
    font-size: 11px;
    font-weight: 500;
    color: #8d97a8;
    margin-bottom: 5px;
  }

  .room-detail-kv {
    display: flex;
    gap: 6px;
    font-size: 12px;
    padding: 2px 0;
  }

  .room-detail-key {
    color: #8d97a8;
    width: 56px;
    flex-shrink: 0;
  }

  .room-detail-val {
    color: #2b3444;
  }

  .iptable {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
  }

  .iptable th {
    font-size: 11px;
    font-weight: 500;
    color: #8d97a8;
    text-align: left;
    padding: 6px 10px;
    border-bottom: 1px solid #e4e7ee;
    background: #f7f9fd;
    white-space: nowrap;
  }

  .iptable td {
    font-size: 12px;
    padding: 7px 10px;
    border-bottom: 1px solid #e4e7ee;
    vertical-align: middle;
    color: #687285;
  }

  .iptable tr:last-child td {
    border-bottom: none;
  }

  .iptable tr.nat-expand-row td {
    padding: 0;
    border-bottom: none;
  }

  .iptable tr:not(.nat-expand-row):hover td {
    background: #f7f9fd;
  }

  .td-p {
    color: #2b3444;
  }

  .src-idc,
  .src-isp {
    display: inline-block;
    padding: 1px 6px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 500;
  }

  .src-idc {
    background: #e1f5ee;
    color: #0f6e56;
  }

  .src-isp {
    background: #e6f1fb;
    color: #185fa5;
  }

  .use-static,
  .use-dynamic,
  .use-leased {
    font-size: 11px;
    font-weight: 500;
    padding: 1px 7px;
    border-radius: 10px;
  }

  .use-static {
    background: #e6f1fb;
    color: #185fa5;
  }

  .use-dynamic {
    background: #faeeda;
    color: #854f0b;
  }

  .use-leased {
    background: #eeedfe;
    color: #534ab7;
  }

  .nat-btn {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 1px 7px;
    border-radius: 10px;
    font-size: 11px;
    font-weight: 500;
    background: #faeeda;
    color: #854f0b;
    border: none;
    cursor: pointer;
    font-family: inherit;
    white-space: nowrap;
  }

  .nat-empty {
    color: #9ca3af;
    font-size: 11px;
  }

  .nat-expand-detail {
    display: none;
    padding: 6px 10px 7px 22px;
    background: #f7f9fd;
    border-top: 1px solid #e4e7ee;
  }

  .nat-expand-detail.open {
    display: block;
  }

  .nat-line {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 2px 0;
    font-family: Menlo, Consolas, Monaco, monospace;
    font-size: 11px;
    color: #687285;
  }

  .nat-arr {
    color: #8d97a8;
  }

  .show-more-btn {
    display: block;
    width: 100%;
    padding: 7px;
    font-size: 12px;
    color: #8d97a8;
    background: #f7f9fd;
    border: none;
    border-top: 1px solid #e4e7ee;
    cursor: pointer;
    font-family: inherit;
    text-align: center;
  }

  .show-more-btn:hover {
    color: #2b3444;
  }

  .app-table {
    width: 100%;
    border-collapse: collapse;
  }

  .app-table th {
    font-size: 11px;
    font-weight: 500;
    color: #8d97a8;
    text-align: left;
    padding: 6px 10px;
    border-bottom: 1px solid #e4e7ee;
    background: #f7f9fd;
  }

  .app-table td {
    font-size: 12px;
    padding: 7px 10px;
    border-bottom: 1px solid #e4e7ee;
    color: #687285;
    vertical-align: middle;
  }

  .app-table tr:last-child td {
    border-bottom: none;
  }

  .ip-pill {
    font-family: Menlo, Consolas, Monaco, monospace;
    font-size: 11px;
    padding: 1px 6px;
    background: #eff3fa;
    border: 1px solid #d8deea;
    border-radius: 4px;
    color: #687285;
  }

  @media (max-width: 1024px) {
    .metrics-row {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .two-col {
      grid-template-columns: 1fr;
    }

    .room-expand-inner.open {
      grid-template-columns: 1fr;
    }
  }
</style>
