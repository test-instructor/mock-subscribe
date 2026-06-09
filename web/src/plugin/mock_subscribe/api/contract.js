import service from '@/utils/request'

export const getContractList = (params) => {
  return service({
    url: '/mockSubscribeContract/getContractList',
    method: 'get',
    params
  })
}

export const findContract = (params) => {
  return service({
    url: '/mockSubscribeContract/findContract',
    method: 'get',
    params
  })
}

export const updateContractStatus = (data) => {
  return service({
    url: '/mockSubscribeContract/updateContractStatus',
    method: 'put',
    data
  })
}

export const getContractRecordList = (params) => {
  return service({
    url: '/mockSubscribeContract/getContractRecordList',
    method: 'get',
    params
  })
}
