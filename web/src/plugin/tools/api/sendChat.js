import service from '@/utils/request'

export const createSendChatTask = (data) => {
  return service({
    url: '/toolsSendChat/createSendChatTask',
    method: 'post',
    data
  })
}

export const getSendChatTaskList = (params) => {
  return service({
    url: '/toolsSendChat/getSendChatTaskList',
    method: 'get',
    params
  })
}

export const stopSendChatTask = (params) => {
  return service({
    url: '/toolsSendChat/stopSendChatTask',
    method: 'put',
    params
  })
}
