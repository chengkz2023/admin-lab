<template>
  <div class="flex justify-between items-center gap-2 w-full">
    <el-cascader
      v-if="pathIsSelect"
      v-model="activeComponent"
      placeholder="请选择文件路径"
      :options="pathOptions"
      filterable
      class="!w-full"
      clearable
      @change="emitChange"
    />
    <el-input
      v-else
      v-model="tempPath"
      placeholder="页面:view/xxx/xx.vue"
      @change="emitChange"
    />
    <el-button @click="togglePathIsSelect">
      {{ pathIsSelect ? '手动输入' : '快捷选择' }}
    </el-button>
  </div>
</template>

<script setup>
  import { onMounted, ref, watch } from 'vue'
  import pathInfo from '@/pathInfo.json'

  const props = defineProps({
    component: {
      type: String,
      default: ''
    }
  })

  const emits = defineEmits(['change'])

  const pathOptions = ref([])
  const tempPath = ref('')
  const activeComponent = ref([])
  const pathIsSelect = ref(true)

  const togglePathIsSelect = () => {
    if (pathIsSelect.value) {
      tempPath.value = activeComponent.value?.join('/') || ''
    } else {
      activeComponent.value = tempPath.value?.split('/') || []
    }

    pathIsSelect.value = !pathIsSelect.value
    emitChange()
  }

  function convertToCascaderOptions(data) {
    const result = []

    for (const path in data) {
      const label = data[path]
      const parts = path.split('/').filter(Boolean)
      const startIndex = parts[0] === 'src' ? 1 : 0
      let currentLevel = result

      for (let i = startIndex; i < parts.length; i++) {
        const part = parts[i]
        let node = currentLevel.find((item) => item.value === part)

        if (!node) {
          node = {
            value: part,
            label: part,
            children: []
          }
          currentLevel.push(node)
        }

        if (i === parts.length - 1) {
          node.label = label
          delete node.children
        }

        currentLevel = node.children || []
      }
    }

    return result
  }

  const initCascader = (value) => {
    if (value === '') {
      pathIsSelect.value = true
      return
    }

    if (pathInfo[`/src/${value}`]) {
      activeComponent.value = value.split('/').filter(Boolean)
      tempPath.value = ''
      pathIsSelect.value = true
      return
    }

    tempPath.value = value
    activeComponent.value = []
    pathIsSelect.value = false
  }

  const emitChange = () => {
    emits(
      'change',
      pathIsSelect.value ? activeComponent.value?.join('/') : tempPath.value
    )
  }

  watch(
    () => props.component,
    (value) => {
      initCascader(value)
    }
  )

  onMounted(() => {
    pathOptions.value = convertToCascaderOptions(pathInfo)
    initCascader(props.component)
  })
</script>

<style scoped lang="scss"></style>
