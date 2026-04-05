<!-- web/src/view/lab/reusable/biz-log.vue -->
<template>
  <div class="lab-page">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 业务操作日志</p>
        <h2>业务操作日志</h2>
        <p class="subtitle">
          通用日志表 + 异步 LogRecorder + BizLogTimeline 前端时间线组件。
          Service 层埋点一行代码，详情页嵌入组件一个标签，迁入只需复制
          utils/bizlog、model 和前端组件。
        </p>
      </div>
      <div class="hero-meta">
        <el-tag type="success">前后端</el-tag>
        <el-tag type="primary">可迁移</el-tag>
        <el-tag>通用日志表</el-tag>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="14">
        <el-card shadow="hover">
          <template #header>
            <div class="panel-header">
              <span>日志时间线 Demo</span>
              <div style="display: flex; gap: 8px; align-items: center">
                <el-input v-model="demoModule" placeholder="module" style="width: 120px" size="small" />
                <el-input v-model="demoEntityId" placeholder="entityId" style="width: 100px" size="small" />
                <el-button size="small" @click="writeTestLog" :loading="writing">写入测试日志</el-button>
                <el-button size="small" type="primary" @click="refreshTimeline">刷新</el-button>
              </div>
            </div>
          </template>
          <biz-log-timeline :key="timelineKey" :module="demoModule" :entity-id="demoEntityId" />
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="10">
        <el-card shadow="hover">
          <template #header><span class="panel-title">后端埋点示例</span></template>
          <pre class="code-block">import "github.com/your-org/server/utils/bizlog"

// 在 service 层操作完成后调用
bizlog.Record(ctx, bizlog.Entry{
    Module:       "order",
    EntityID:     strconv.FormatUint(order.ID, 10),
    Action:       "approve",
    OperatorID:   utils.GetUserID(c),
    OperatorName: utils.GetUserName(c),
    Remark:       fmt.Sprintf("订单 #%d 审批通过", order.ID),
})</pre>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="panel-title">前端组件用法</span></template>
          <pre class="code-block">&lt;biz-log-timeline
  module="order"
  :entity-id="String(orderId)"
/&gt;</pre>
          <p class="hint">组件自行请求数据，父页面零感知。支持分页（默认每页 20 条）。</p>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="panel-title">迁移说明</span></template>
          <div class="step-list">
            <div class="step-item"><span>1</span><div>执行建表或开启 AutoMigrate，确认 <code>biz_log</code> 表存在。</div></div>
            <div class="step-item"><span>2</span><div>复制 <code>server/utils/bizlog/recorder.go</code>，修改 DB 实例来源为项目自身。</div></div>
            <div class="step-item"><span>3</span><div>在 service 层关键操作后调用 <code>bizlog.Record()</code>。</div></div>
            <div class="step-item"><span>4</span><div>注册路由和 Casbin 权限，复制前端组件，在详情页引用。</div></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import BizLogTimeline from '@/components/lab/biz-log-timeline.vue'
  import service from '@/utils/request'

  defineOptions({ name: 'LabReusableBizLog' })

  const demoModule = ref('demo')
  const demoEntityId = ref('1')
  const timelineKey = ref(0)
  const writing = ref(false)

  const refreshTimeline = () => { timelineKey.value++ }

  const writeTestLog = async () => {
    writing.value = true
    try {
      const res = await service({
        url: '/bizLog/writeDemo',
        method: 'post',
        data: { module: demoModule.value, entityId: demoEntityId.value }
      })
      if (res.code === 0) {
        ElMessage.success('测试日志已写入')
        refreshTimeline()
      }
    } finally {
      writing.value = false
    }
  }
</script>

<style scoped>
  .lab-page { display: flex; flex-direction: column; gap: 16px; }
  .hero {
    display: flex; justify-content: space-between; gap: 16px;
    padding: 24px; border-radius: 16px;
    border: 1px solid #dbeafe;
    background: linear-gradient(135deg, #f0f9ff 0%, #f8fafc 100%);
  }
  .eyebrow { margin: 0 0 8px; color: #0369a1; font-size: 13px; font-weight: 700; letter-spacing: 0.08em; }
  .hero h2 { margin: 0 0 8px; color: #0f172a; font-size: 28px; }
  .subtitle { margin: 0; color: #475569; line-height: 1.7; max-width: 600px; }
  .hero-meta { display: flex; flex-wrap: wrap; align-content: flex-start; gap: 8px; }
  .panel-header { display: flex; justify-content: space-between; align-items: center; font-weight: 600; }
  .panel-title { font-weight: 600; }
  .code-block {
    margin: 0; padding: 12px 16px; border-radius: 8px;
    background: #f1f5f9; font-family: Consolas, 'Courier New', monospace;
    font-size: 12px; line-height: 1.65; white-space: pre-wrap; word-break: break-word;
  }
  .hint { margin: 8px 0 0; color: #475569; font-size: 13px; }
  .step-list { display: flex; flex-direction: column; gap: 10px; }
  .step-item { display: grid; grid-template-columns: 28px 1fr; gap: 12px; align-items: start; }
  .step-item span {
    display: inline-flex; align-items: center; justify-content: center;
    width: 28px; height: 28px; border-radius: 999px;
    background: #0ea5e9; color: #fff; font-size: 13px; font-weight: 700;
  }
</style>
