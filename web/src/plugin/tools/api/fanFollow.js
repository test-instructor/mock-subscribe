import service from '@/utils/request'

export const createFanFollow = (data) => {
  return service({
    url: '/toolsFanFollow/createFanFollow',
    method: 'post',
    data
  })
}

export const getFanFollowList = (params) => {
  return service({
    url: '/toolsFanFollow/getFanFollowList',
    method: 'get',
    params
  })
}
