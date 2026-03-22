<template>
  <el-dialog
    :model-value="modelValue"
    :title="dialogTitle"
    :width="width"
    :destroy-on-close="destroyOnClose"
    @update:model-value="(value) => emit('update:modelValue', value)"
    @closed="handleClosed"
  >
    <el-form
      ref="formRef"
      :model="localForm"
      :rules="rules"
      :label-width="labelWidth"
      class="crud-form"
      @submit.prevent
    >
      <el-row :gutter="12">
        <el-col
          v-for="item in items"
          :key="item.prop"
          :xs="24"
          :sm="item.sm || 24"
          :md="item.md || 12"
          :lg="item.lg || item.span || 12"
        >
          <el-form-item :label="item.label" :prop="item.prop">
            <slot
              v-if="$slots[`field-${item.prop}`]"
              :name="`field-${item.prop}`"
              :model-value="localForm[item.prop]"
              :update-value="(value) => updateField(item.prop, value)"
              :item="item"
            />
            <el-input
              v-else-if="item.type === 'input' || !item.type"
              v-model="localForm[item.prop]"
              :placeholder="item.placeholder || `请输入${item.label}`"
              clearable
            />
            <el-input
              v-else-if="item.type === 'textarea'"
              v-model="localForm[item.prop]"
              type="textarea"
              :rows="item.rows || 3"
              :placeholder="item.placeholder || `请输入${item.label}`"
              maxlength="500"
              show-word-limit
            />
            <el-select
              v-else-if="item.type === 'select'"
              v-model="localForm[item.prop]"
              :placeholder="item.placeholder || `请选择${item.label}`"
              :multiple="!!item.multiple"
              clearable
              filterable
            >
              <el-option
                v-for="option in item.options || []"
                :key="`${item.prop}-${option.value}`"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
            <el-switch
              v-else-if="item.type === 'switch'"
              v-model="localForm[item.prop]"
              :active-text="item.activeText"
              :inactive-text="item.inactiveText"
            />
            <el-date-picker
              v-else-if="item.type === 'date'"
              v-model="localForm[item.prop]"
              type="date"
              value-format="YYYY-MM-DD"
              :placeholder="item.placeholder || `请选择${item.label}`"
              clearable
            />
            <el-date-picker
              v-else-if="item.type === 'dateRange'"
              v-model="localForm[item.prop]"
              type="daterange"
              value-format="YYYY-MM-DD"
              range-separator="至"
              :start-placeholder="item.startPlaceholder || '开始日期'"
              :end-placeholder="item.endPlaceholder || '结束日期'"
              clearable
            />
            <el-cascader
              v-else-if="item.type === 'cascader'"
              v-model="localForm[item.prop]"
              :options="item.options || []"
              :props="resolveCascaderProps(item)"
              :placeholder="item.placeholder || `请选择${item.label}`"
              clearable
              filterable
              collapse-tags
              collapse-tags-tooltip
            />
            <el-input
              v-else
              v-model="localForm[item.prop]"
              :placeholder="item.placeholder || `请输入${item.label}`"
              clearable
            />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">
          {{ cancelText }}
        </el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ confirmText || defaultConfirmText }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
  import { computed, ref, watch } from 'vue'

  defineOptions({
    name: 'LabCrudFormDialog'
  })

  const props = defineProps({
    modelValue: {
      type: Boolean,
      default: false
    },
    mode: {
      type: String,
      default: 'add'
    },
    formData: {
      type: Object,
      default: () => ({})
    },
    items: {
      type: Array,
      default: () => []
    },
    rules: {
      type: Object,
      default: () => ({})
    },
    loading: {
      type: Boolean,
      default: false
    },
    width: {
      type: String,
      default: '760px'
    },
    labelWidth: {
      type: String,
      default: '96px'
    },
    addTitle: {
      type: String,
      default: '新增'
    },
    editTitle: {
      type: String,
      default: '编辑'
    },
    confirmText: {
      type: String,
      default: ''
    },
    cancelText: {
      type: String,
      default: '取消'
    },
    destroyOnClose: {
      type: Boolean,
      default: true
    }
  })

  const emit = defineEmits(['update:modelValue', 'update:formData', 'submit', 'cancel', 'closed'])

  const formRef = ref(null)
  const localForm = ref({})

  const dialogTitle = computed(() => (props.mode === 'edit' ? props.editTitle : props.addTitle))
  const defaultConfirmText = computed(() => (props.mode === 'edit' ? '保存' : '创建'))

  const defaultByType = (item) => {
    if (item.defaultValue !== undefined) {
      return item.defaultValue
    }
    if (item.type === 'switch') {
      return false
    }
    if (item.type === 'select' && item.multiple) {
      return []
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
    return ''
  }

  const buildInitialForm = () => {
    const result = {}
    props.items.forEach((item) => {
      result[item.prop] = defaultByType(item)
    })
    return result
  }

  const resetLocalForm = () => {
    localForm.value = {
      ...buildInitialForm(),
      ...(props.formData || {})
    }
  }

  const updateField = (prop, value) => {
    localForm.value[prop] = value
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

  const handleCancel = () => {
    emit('update:modelValue', false)
    emit('cancel')
  }

  const handleClosed = () => {
    formRef.value?.clearValidate?.()
    emit('closed')
  }

  const handleSubmit = async () => {
    const valid = await formRef.value?.validate?.().catch(() => false)
    if (!valid) {
      return
    }
    const payload = { ...localForm.value }
    emit('update:formData', payload)
    emit('submit', payload)
  }

  watch(
    () => props.modelValue,
    (visible) => {
      if (visible) {
        resetLocalForm()
      }
    },
    { immediate: true }
  )
</script>

<style scoped>
  .crud-form {
    padding-top: 6px;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
</style>
