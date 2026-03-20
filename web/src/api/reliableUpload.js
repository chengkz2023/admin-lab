import service from '@/utils/request'

export const getReliableUploadProfile = () => {
  return service({
    url: '/reliableUpload/profile',
    method: 'get'
  })
}
