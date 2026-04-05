<!-- web/src/components/lab/biz-log-timeline.vue -->
<template>
  <div v-loading="loading">
    <el-empty v-if="!loading && !list.length" description="暂无操作记录" :image-size="60" />
    <el-timeline v-else>
      <el-timeline-item
        v-for="item in list"
        :key="item.id"
        :timestamp="formatTime(item.createdAt)"
        placement="top"
        type="primary"
      >
        <div class="log-item">
          <span class="operator">{{ item.operatorName }}</span>
          <el-tag size="small" style="margin: 0 6px">{{ item.action }}</el-tag>
          <span class="remark">{{ item.remark }}</span>
        </div>
      </el-timeline-item>
    </el-timeline>
    <div v-if="total > pageSize" class="pagination">
      <el-pagination
        small
        layout="prev, pager, next"
        :total="total"
        :page-size="pageSize"
        v-model:current-page="page"
        @current-change="load"
      />
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { getBizLogList } from '@/api/bizLog'

  defineOptions({ name: 'BizLogTimeline' })

  const props = defineProps({
    module: { type: String, required: true },
    entityId: { type: String, required: true }
  })

  const loading = ref(false)
  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = 20

  const formatTime = (ts) => {
    if (!ts) return ''
    return new Date(ts).toLocaleString('zh-CN', { hour12: false })
  }

  const load = async () => {
    loading.value = true
    try {
      const res = await getBizLogList({ module: props.module, entityId: props.entityId, page: page.value, pageSize })
      if (res.code === 0) {
        list.value = res.data.list || []
        total.value = res.data.total || 0
      }
    } finally {
      loading.value = false
    }
  }

  onMounted(load)
</script>

<style scoped>
  .log-item { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
  .operator { font-weight: 600; color: #0f172a; }
  .remark { color: #475569; }
  .pagination { margin-top: 12px; display: flex; justify-content: center; }
</style>
