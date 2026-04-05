<template>
  <el-tag :type="tagType">{{ label }}</el-tag>
</template>

<script setup>
  import { ref, watch } from 'vue'
  import { useDictionaryStore } from '@/pinia/modules/dictionary'

  defineOptions({ name: 'DictTag' })

  const props = defineProps({
    dictType: { type: String, required: true },
    value: { required: true }
  })

  const VALID_TYPES = ['success', 'warning', 'danger', 'info', '']
  const dictStore = useDictionaryStore()
  const label = ref(String(props.value))
  const tagType = ref('')

  watch(
    () => [props.dictType, props.value],
    async ([type, val]) => {
      if (!type) return
      try {
        const items = await dictStore.getDictionary(type, 1)
        const matched = (items || []).find((i) => String(i.value) === String(val))
        if (matched) {
          label.value = matched.label
          tagType.value = VALID_TYPES.includes(matched.extend) ? matched.extend : ''
        } else {
          label.value = val != null ? String(val) : ''
          tagType.value = ''
        }
      } catch {
        label.value = val != null ? String(val) : ''
        tagType.value = ''
      }
    },
    { immediate: true }
  )
</script>
