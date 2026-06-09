import service from '@/utils/request'

export const getCallbackRecordList = (params) => {
  return service({
    url: '/papay/callback-records',
    method: 'get',
    params
  })
}

export const findCallbackRecord = (params) => {
  return service({
    url: '/papay/callback-record',
    method: 'get',
    params
  })
}
