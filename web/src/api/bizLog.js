// web/src/api/bizLog.js
import service from '@/utils/request'

export const getBizLogList = (params) =>
  service({ url: '/bizLog/list', method: 'get', params })

export const writeBizLogDemo = (data) =>
  service({ url: '/bizLog/writeDemo', method: 'post', data })
