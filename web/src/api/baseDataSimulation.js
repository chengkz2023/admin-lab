import service from '@/utils/request'

export const getBaseDataTemplateOptions = () => {
  return service({
    url: '/baseDataSimulation/templates',
    method: 'get'
  })
}

export const downloadBaseDataTemplate = (templateKey) => {
  return service({
    url: '/baseDataSimulation/template',
    method: 'get',
    params: templateKey ? { templateKey } : undefined,
    responseType: 'blob'
  })
}

export const exportBaseData = (templateKey) => {
  return service({
    url: '/baseDataSimulation/export',
    method: 'get',
    params: templateKey ? { templateKey } : undefined,
    responseType: 'blob'
  })
}

export const importBaseData = (data) => {
  return service({
    url: '/baseDataSimulation/import',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
