<template>
  <div class="h-screen w-screen bg-gray-50 text-slate-700 dark:bg-slate-800 dark:text-slate-500">
    <iframe
      v-if="reloadFlag"
      id="gva-base-load-dom"
      class="gva-body-h w-full border-t border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
      :src="url"
    ></iframe>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, reactive, watchEffect } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { emitter } from '@/utils/bus.js'
import useResponsive from '@/hooks/responsive'
import { useUserStore } from '@/pinia/modules/user'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'GvaLayoutIframe'
})

useResponsive(true)

const appStore = useAppStore()
const { isDark } = storeToRefs(appStore)
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const font = reactive({
  color: 'rgba(0, 0, 0, .15)'
})

watchEffect(() => {
  font.color = isDark.value ? 'rgba(255,255,255, .15)' : 'rgba(0, 0, 0, .15)'
})

const url = route.query.url || 'about:blank'

onMounted(() => {
  emitter.on('reload', reload)
  if (userStore.loadingInstance) {
    userStore.loadingInstance.close()
  }
})

const reloadFlag = ref(true)
let reloadTimer = null

const reload = async () => {
  if (reloadTimer) {
    window.clearTimeout(reloadTimer)
  }
  reloadTimer = window.setTimeout(async () => {
    if (route.meta.keepAlive) {
      reloadFlag.value = false
      await nextTick()
      reloadFlag.value = true
    } else {
      const title = route.meta.title
      router.push({ name: 'Reload', params: { title } })
    }
  }, 400)
}
</script>
