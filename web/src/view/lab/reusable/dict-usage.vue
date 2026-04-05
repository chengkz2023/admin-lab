<!-- web/src/view/lab/reusable/dict-usage.vue -->
<template>
  <div class="lab-page">
    <div class="hero">
      <div>
        <p class="eyebrow">复用组件 / 字典消费端</p>
        <h2>字典消费组件</h2>
        <p class="subtitle">
          基于 useDictionaryStore，提供 DictSelect、DictTag 和 getLabel
          三种消费方式，字典数据按需加载并自动缓存，迁入只需复制两个组件。
        </p>
      </div>
      <div class="hero-meta">
        <el-tag type="success">前端组件</el-tag>
        <el-tag type="primary">可迁移</el-tag>
        <el-tag>零后端依赖</el-tag>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover">
          <template #header><span class="card-title">DictSelect — 下拉选择</span></template>
          <p class="hint">用法：<code>&lt;dict-select dict-type="sex" v-model="val" /&gt;</code></p>
          <div class="demo-row">
            <dict-select dict-type="sex" v-model="selectVal" style="width: 200px" placeholder="请选择" />
            <span class="result">当前值：{{ selectVal ?? '—' }}</span>
          </div>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="card-title">DictTag — Tag 展示</span></template>
          <p class="hint">用法：<code>&lt;dict-tag dict-type="sex" :value="1" /&gt;</code></p>
          <div class="demo-row" style="gap: 8px; flex-wrap: wrap">
            <dict-tag dict-type="sex" :value="1" />
            <dict-tag dict-type="sex" :value="2" />
            <dict-tag dict-type="sex" :value="999" />
          </div>
          <p class="hint" style="margin-top: 8px">
            颜色来自字典值 <code>extend</code> 字段，在字典管理后台配置为
            <code>success / warning / danger / info</code> 即可生效。
          </p>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card shadow="hover">
          <template #header><span class="card-title">getLabel — 纯文字翻译</span></template>
          <p class="hint">直接调用 store 方法，适合在 js/模板插值中翻译枚举值。</p>
          <pre class="code-block">import { useDictionaryStore } from '@/pinia/modules/dictionary'

const dictStore = useDictionaryStore()
const items = await dictStore.getDictionary('sex', 1)
const label = items.find(i => String(i.value) === String(val))?.label ?? val</pre>
        </el-card>

        <el-card shadow="hover">
          <template #header><span class="card-title">迁移说明</span></template>
          <div class="step-list">
            <div class="step-item"><span>1</span><div>确认内网项目已有 <code>useDictionaryStore</code>（gin-vue-admin 标准内置）。</div></div>
            <div class="step-item"><span>2</span><div>复制 <code>components/lab/dict-select.vue</code> 和 <code>dict-tag.vue</code>，全局注册或按需引入。</div></div>
            <div class="step-item"><span>3</span><div>在字典管理后台为每个字典值的 <code>extend</code> 字段填入颜色标记（可选）。</div></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import DictSelect from '@/components/lab/dict-select.vue'
  import DictTag from '@/components/lab/dict-tag.vue'

  defineOptions({ name: 'LabReusableDictUsage' })

  const selectVal = ref(null)
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
  .card-title { font-weight: 600; }
  .hint { margin: 0 0 12px; color: #475569; font-size: 13px; }
  .demo-row { display: flex; align-items: center; gap: 16px; }
  .result { color: #0f172a; font-size: 14px; }
  .code-block {
    margin: 0; padding: 12px 16px; border-radius: 8px;
    background: #f1f5f9; font-family: Consolas, 'Courier New', monospace;
    font-size: 13px; line-height: 1.65; white-space: pre-wrap; word-break: break-word;
  }
  .step-list { display: flex; flex-direction: column; gap: 10px; }
  .step-item { display: grid; grid-template-columns: 28px 1fr; gap: 12px; align-items: start; }
  .step-item span {
    display: inline-flex; align-items: center; justify-content: center;
    width: 28px; height: 28px; border-radius: 999px;
    background: #0ea5e9; color: #fff; font-size: 13px; font-weight: 700;
  }
</style>
