import service from '@/utils/request'

export const getExcelTemplateOptions = () => {
  return service({
    url: '/excelIO/templates',
    method: 'get'
  })
}

export const downloadExcelTemplate = (templateKey) => {
  return service({
    url: '/excelIO/template',
    method: 'get',
    params: templateKey ? { templateKey } : undefined,
    responseType: 'blob'
  })
}

export const exportExcelSample = () => {
  return service({
    url: '/excelIO/export',
    method: 'get',
    responseType: 'blob'
  })
}

export const importExcelFile = (data) => {
  return service({
    url: '/excelIO/import',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
