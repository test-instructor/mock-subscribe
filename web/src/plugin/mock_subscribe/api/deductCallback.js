import service from '@/utils/request'

export const getDeductCallbackRecordList = (params) => {
  return service({
    url: '/mockSubscribeDeductCallback/getDeductCallbackRecordList',
    method: 'get',
    params
  })
}

export const findDeductCallbackRecord = (params) => {
  return service({
    url: '/mockSubscribeDeductCallback/findDeductCallbackRecord',
    method: 'get',
    params
  })
}
