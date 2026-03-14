import service from '@/utils/request'

export const getSystemConfig = () => {
  return service({
    url: '/system/getSystemConfig',
    method: 'post'
  })
}

export const setSystemConfig = (data) => {
  return service({
    url: '/system/setSystemConfig',
    method: 'post',
    data
  })
}

export const reloadSystem = (data) => {
  return service({
    url: '/system/reloadSystem',
    method: 'post',
    data
  })
}
