import service from '@/utils/request'

export const getDirFilePipelineProfile = () => {
  return service({
    url: '/dirPipeline/profile',
    method: 'get'
  })
}

export const runDirFilePipelineOnce = (data) => {
  return service({
    url: '/dirPipeline/runOnce',
    method: 'post',
    data
  })
}
