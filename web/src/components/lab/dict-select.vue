<!-- web/src/components/lab/dict-select.vue -->
<template>
  <el-select v-bind="$attrs" :loading="loading" :disabled="loading || $attrs.disabled">
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    />
  </el-select>
</template>

<script setup>
  import { ref, watch } from 'vue'
  import { useDictionaryStore } from '@/pinia/modules/dictionary'

  defineOptions({ name: 'DictSelect', inheritAttrs: false })

  const props = defineProps({
    dictType: { type: String, required: true },
    depth:    { type: Number, default: 1 }
  })

  const dictStore = useDictionaryStore()
  const options = ref([])
  const loading = ref(false)

  watch(
    () => props.dictType,
    async (val) => {
      if (!val) return
      loading.value = true
      try {
        const items = await dictStore.getDictionary(val, props.depth)
        options.value = items || []
      } finally {
        loading.value = false
      }
    },
    { immediate: true }
  )
</script>
