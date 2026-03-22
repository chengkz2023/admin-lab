import service from '@/utils/request'

export const getSecurityDashboardPanel = (params) => {
  return service({
    url: '/securityDashboard/panel',
    method: 'get',
    params
  })
}

export const getSecurityDashboardDrilldown = (params) => {
  return service({
    url: '/securityDashboard/drilldown',
    method: 'get',
    params
  })
}
