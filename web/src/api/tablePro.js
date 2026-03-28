import service from '@/utils/request'

export const getTableProPage = (data) => {
  return service({
    url: '/tablePro/page',
    method: 'post',
    data
  })
}

export const exportTablePro = (data) => {
  return service({
    url: '/tablePro/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
}
