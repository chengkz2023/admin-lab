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
  import { ref, onMounted } from 'vue'
  import { useDictionaryStore } from '@/pinia/modules/dictionary'

  defineOptions({ name: 'DictSelect', inheritAttrs: false })

  const props = defineProps({
    dictType: { type: String, required: true }
  })

  const dictStore = useDictionaryStore()
  const options = ref([])
  const loading = ref(false)

  onMounted(async () => {
    loading.value = true
    try {
      const items = await dictStore.getDictionary(props.dictType, 1)
      options.value = items || []
    } finally {
      loading.value = false
    }
  })
</script>
