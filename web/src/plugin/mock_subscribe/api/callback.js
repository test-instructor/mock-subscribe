import service from '@/utils/request'

export const getCallbackRecordList = (params) => {
  return service({
    url: '/mockSubscribeCallback/getCallbackRecordList',
    method: 'get',
    params
  })
}

export const findCallbackRecord = (params) => {
  return service({
    url: '/mockSubscribeCallback/findCallbackRecord',
    method: 'get',
    params
  })
}
