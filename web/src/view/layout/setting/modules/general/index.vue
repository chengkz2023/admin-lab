<template>
  <div class="gva-theme-font">
    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">系统信息</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div class="flex items-center justify-between border-b border-gray-200 py-3 dark:border-gray-600">
              <span class="gva-theme-text-sub font-medium">前端框架</span>
              <span class="gva-theme-text-main font-mono font-semibold">Vue 3</span>
            </div>
            <div class="flex items-center justify-between border-b border-gray-200 py-3 dark:border-gray-600">
              <span class="gva-theme-text-sub font-medium">UI 组件</span>
              <span class="gva-theme-text-main font-mono font-semibold">Element Plus</span>
            </div>
            <div class="flex items-center justify-between border-b border-gray-200 py-3 dark:border-gray-600">
              <span class="gva-theme-text-sub font-medium">构建工具</span>
              <span class="gva-theme-text-main font-mono font-semibold">Vite</span>
            </div>
            <div class="flex items-center justify-between border-b border-gray-200 py-3 dark:border-gray-600">
              <span class="gva-theme-text-sub font-medium">当前主题</span>
              <span class="gva-theme-text-main font-mono font-semibold">{{ config.darkMode }}</span>
            </div>
            <div class="flex items-center justify-between py-3">
              <span class="gva-theme-text-sub font-medium">浏览器</span>
              <span class="gva-theme-text-main font-mono font-semibold">{{ browserInfo }}</span>
            </div>
            <div class="flex items-center justify-between py-3">
              <span class="gva-theme-text-sub font-medium">屏幕分辨率</span>
              <span class="gva-theme-text-main font-mono font-semibold">{{ screenResolution }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">配置管理</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <div class="space-y-5">
            <div class="gva-theme-card-white flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-xl border border-red-200 bg-red-50 text-xl text-red-600 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
                  R
                </div>
                <div>
                  <h4 class="gva-theme-text-main text-sm font-semibold">重置界面配置</h4>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">恢复为当前脚手架默认显示设置</p>
                </div>
              </div>
              <el-button
                type="danger"
                size="small"
                class="rounded-lg font-medium transition-all duration-150 ease-in-out hover:-translate-y-0.5"
                @click="handleResetConfig"
              >
                重置配置
              </el-button>
            </div>

            <div class="gva-theme-card-white flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-xl border border-blue-200 bg-blue-50 text-xl text-blue-600 dark:border-blue-800 dark:bg-blue-900/20 dark:text-blue-400">
                  E
                </div>
                <div>
                  <h4 class="gva-theme-text-main text-sm font-semibold">导出配置</h4>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">导出当前界面设置为 JSON 文件</p>
                </div>
              </div>
              <el-button
                type="primary"
                size="small"
                class="rounded-lg font-medium transition-all duration-150 ease-in-out hover:-translate-y-0.5"
                :style="{ backgroundColor: config.primaryColor, borderColor: config.primaryColor }"
                @click="handleExportConfig"
              >
                导出配置
              </el-button>
            </div>

            <div class="gva-theme-card-white flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="flex h-12 w-12 items-center justify-center rounded-xl border border-green-200 bg-green-50 text-xl text-green-600 dark:border-green-800 dark:bg-green-900/20 dark:text-green-400">
                  I
                </div>
                <div>
                  <h4 class="gva-theme-text-main text-sm font-semibold">导入配置</h4>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">从 JSON 文件恢复界面设置</p>
                </div>
              </div>
              <el-upload
                ref="uploadRef"
                :auto-upload="false"
                :show-file-list="false"
                accept=".json"
                @change="handleImportConfig"
              >
                <el-button
                  type="success"
                  size="small"
                  class="rounded-lg font-medium transition-all duration-150 ease-in-out hover:-translate-y-0.5"
                >
                  导入配置
                </el-button>
              </el-upload>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mb-10">
      <div class="gva-theme-section-header">
        <div class="gva-theme-divider"></div>
        <span class="gva-theme-section-title">关于项目</span>
        <div class="gva-theme-divider"></div>
      </div>

      <div class="gva-theme-section-content">
        <div class="gva-theme-card-bg">
          <div class="flex items-start gap-5">
            <div class="flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-xl border border-gray-200 bg-white shadow-sm dark:border-gray-600 dark:bg-gray-700">
              <Logo />
            </div>
            <div class="flex-1">
              <h4 class="gva-theme-text-main mb-3 text-xl font-semibold">admin-lab</h4>
              <p class="gva-theme-text-sub mb-2 text-sm leading-relaxed">
                当前项目已整理为内部使用的全栈后台脚手架，保留系统管理能力和可持续扩展的前后端结构。
              </p>
              <p class="gva-theme-text-sub text-sm leading-relaxed">
                如果后续还要继续补品牌信息、官网地址或团队介绍，可以再在这里按你的需要微调。
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/pinia'
import Logo from '@/components/logo/index.vue'

defineOptions({
  name: 'GeneralSettings'
})

const appStore = useAppStore()
const { config } = storeToRefs(appStore)
const uploadRef = ref()

const browserInfo = ref('')
const screenResolution = ref('')

onMounted(() => {
  const userAgent = navigator.userAgent
  if (userAgent.includes('Chrome')) {
    browserInfo.value = 'Chrome'
  } else if (userAgent.includes('Firefox')) {
    browserInfo.value = 'Firefox'
  } else if (userAgent.includes('Safari')) {
    browserInfo.value = 'Safari'
  } else if (userAgent.includes('Edge')) {
    browserInfo.value = 'Edge'
  } else {
    browserInfo.value = 'Unknown'
  }

  screenResolution.value = `${screen.width}x${screen.height}`
})

const handleResetConfig = async () => {
  try {
    await ElMessageBox.confirm('确定要重置当前界面配置吗？', '重置配置', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    appStore.resetConfig()
    ElMessage.success('配置已重置')
  } catch {
    // user cancelled
  }
}

const handleExportConfig = () => {
  const configData = JSON.stringify(config.value, null, 2)
  const blob = new Blob([configData], { type: 'application/json' })
  const url = URL.createObjectURL(blob)

  const link = document.createElement('a')
  link.href = url
  link.download = `admin-lab-config-${new Date().toISOString().split('T')[0]}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  ElMessage.success('配置已导出')
}

const handleImportConfig = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const importedConfig = JSON.parse(e.target.result)
      Object.keys(importedConfig).forEach((key) => {
        if (key in config.value) {
          config.value[key] = importedConfig[key]
        }
      })
      ElMessage.success('配置已导入')
    } catch {
      ElMessage.error('配置文件格式错误')
    }
  }
  reader.readAsText(file.raw)
}
</script>
