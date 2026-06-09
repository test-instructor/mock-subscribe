import service from '@/utils/request'

export const getDeductRecordList = (params) => {
  return service({
    url: '/mockSubscribeDeduct/getDeductRecordList',
    method: 'get',
    params
  })
}

export const findDeductRecord = (params) => {
  return service({
    url: '/mockSubscribeDeduct/findDeductRecord',
    method: 'get',
    params
  })
}
