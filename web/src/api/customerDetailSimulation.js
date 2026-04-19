import service from '@/utils/request'

export const getCustomerDetailSimulation = () => {
  return service({
    url: '/customerDetailSimulation/detail',
    method: 'get'
  })
}
