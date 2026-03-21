<template>
  <el-card class="query-bar-card" shadow="never">
    <el-form :label-width="labelWidth" @submit.prevent>
      <el-row :gutter="12">
        <el-col
          v-for="item in visibleItems"
          :key="item.prop"
          :xs="24"
          :sm="item.sm || 12"
          :md="item.md || 8"
          :lg="item.lg || item.span || 6"
        >
          <el-form-item :label="item.label">
            <slot
              v-if="$slots[`field-${item.prop}`]"
              :name="`field-${item.prop}`"
              :model-value="innerModel[item.prop]"
              :update-value="(value) => updateField(item.prop, value)"
              :item="item"
            />
            <el-input
              v-else-if="item.type === 'input' || !item.type"
              :model-value="innerModel[item.prop]"
              :placeholder="item.placeholder || `请输入${item.label}`"
              clearable
              @update:model-value="(value) => updateField(item.prop, value)"
              @keyup.enter="submitOnEnter ? handleSearch() : null"
            />
            <el-select
              v-else-if="item.type === 'select'"
              :model-value="innerModel[item.prop]"
              :placeholder="item.placeholder || `请选择${item.label}`"
              clearable
              filterable
              :multiple="!!item.multiple"
              @update:model-value="(value) => updateField(item.prop, value)"
            >
              <el-option
                v-for="option in item.options || []"
                :key="`${item.prop}-${option.value}`"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
            <el-cascader
              v-else-if="item.type === 'cascader'"
              :model-value="innerModel[item.prop]"
              :options="item.options || []"
              :props="resolveCascaderProps(item)"
              :placeholder="item.placeholder || `请选择${item.label}`"
              clearable
              filterable
              collapse-tags
              collapse-tags-tooltip
              @update:model-value="(value) => updateField(item.prop, value || getDefaultValue(item))"
            />
            <el-date-picker
              v-else-if="item.type === 'dateRange'"
              :model-value="innerModel[item.prop]"
              :type="item.pickerType || 'daterange'"
              :start-placeholder="item.startPlaceholder || '开始日期'"
              :end-placeholder="item.endPlaceholder || '结束日期'"
              range-separator="至"
              value-format="YYYY-MM-DD"
              clearable
              @update:model-value="(value) => updateField(item.prop, value || [])"
            />
            <el-input
              v-else
              :model-value="innerModel[item.prop]"
              :placeholder="item.placeholder || `请输入${item.label}`"
              clearable
              @update:model-value="(value) => updateField(item.prop, value)"
              @keyup.enter="submitOnEnter ? handleSearch() : null"
            />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <div class="query-bar-actions">
      <div class="action-left">
        <slot name="actions-left" />
      </div>
      <div class="action-right">
        <el-button type="primary" :loading="loading" @click="handleSearch">
          {{ searchText }}
        </el-button>
        <el-button @click="handleReset">
          {{ resetText }}
        </el-button>
        <el-button v-if="showExport" @click="emitExport">
          {{ exportText }}
        </el-button>
        <el-button v-if="needToggle" link type="primary" @click="expanded = !expanded">
          {{ expanded ? collapseText : expandText }}
        </el-button>
        <slot name="actions-right" />
      </div>
    </div>
  </el-card>
</template>

<script setup>
  import { computed, ref } from 'vue'

  defineOptions({
    name: 'LabListQueryBar'
  })

  const props = defineProps({
    modelValue: {
      type: Object,
      default: () => ({})
    },
    items: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    maxVisible: {
      type: Number,
      default: 4
    },
    showExport: {
      type: Boolean,
      default: false
    },
    labelWidth: {
      type: String,
      default: '88px'
    },
    searchText: {
      type: String,
      default: '查询'
    },
    resetText: {
      type: String,
      default: '重置'
    },
    exportText: {
      type: String,
      default: '导出'
    },
    expandText: {
      type: String,
      default: '展开更多'
    },
    collapseText: {
      type: String,
      default: '收起'
    },
    submitOnEnter: {
      type: Boolean,
      default: true
    }
  })

  const emit = defineEmits(['update:modelValue', 'search', 'reset', 'export'])

  const expanded = ref(false)

  const needToggle = computed(() => props.items.length > props.maxVisible)
  const visibleItems = computed(() => {
    if (expanded.value || !needToggle.value) {
      return props.items
    }
    return props.items.slice(0, props.maxVisible)
  })

  const innerModel = computed(() => props.modelValue || {})

  const updateField = (prop, value) => {
    emit('update:modelValue', {
      ...innerModel.value,
      [prop]: value
    })
  }

  const getDefaultValue = (item) => {
    if (item.defaultValue !== undefined) {
      return item.defaultValue
    }
    if (item.type === 'dateRange') {
      return []
    }
    if (item.type === 'cascader' && item.multiple) {
      return []
    }
    if (item.type === 'cascader') {
      return item.emitPath === false ? '' : []
    }
    if (item.type === 'select' && item.multiple) {
      return []
    }
    return ''
  }

  const resolveCascaderProps = (item) => {
    const userProps = item.cascaderProps || {}
    return {
      ...userProps,
      multiple: item.multiple ?? userProps.multiple ?? false,
      checkStrictly: item.checkStrictly ?? userProps.checkStrictly ?? false,
      emitPath: item.emitPath ?? userProps.emitPath ?? true
    }
  }

  const handleSearch = () => {
    emit('search', { ...innerModel.value })
  }

  const handleReset = () => {
    const resetModel = {}
    props.items.forEach((item) => {
      resetModel[item.prop] = getDefaultValue(item)
    })
    emit('update:modelValue', resetModel)
    emit('reset', resetModel)
  }

  const emitExport = () => {
    emit('export', { ...innerModel.value })
  }
</script>

<style scoped>
  .query-bar-card {
    border-radius: 12px;
  }

  .query-bar-actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
    margin-top: 4px;
  }

  .action-left,
  .action-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  @media (max-width: 768px) {
    .query-bar-actions {
      align-items: flex-start;
      flex-direction: column;
    }
  }
</style>
