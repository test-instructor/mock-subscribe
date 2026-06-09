import service from '@/utils/request'

export const createMerchant = (data) => {
  return service({
    url: '/mockSubscribeMerchant/createMerchant',
    method: 'post',
    data
  })
}

export const updateMerchant = (data) => {
  return service({
    url: '/mockSubscribeMerchant/updateMerchant',
    method: 'put',
    data
  })
}

export const deleteMerchant = (params) => {
  return service({
    url: '/mockSubscribeMerchant/deleteMerchant',
    method: 'delete',
    params
  })
}

export const findMerchant = (params) => {
  return service({
    url: '/mockSubscribeMerchant/findMerchant',
    method: 'get',
    params
  })
}

export const getMerchantList = (params) => {
  return service({
    url: '/mockSubscribeMerchant/getMerchantList',
    method: 'get',
    params
  })
}
