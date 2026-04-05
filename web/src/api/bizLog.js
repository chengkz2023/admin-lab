// web/src/api/bizLog.js
import service from '@/utils/request'

export const getBizLogList = (params) =>
  service({ url: '/bizLog/list', method: 'get', params })
