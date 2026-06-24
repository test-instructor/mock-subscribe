import service from '@/utils/request'

export const createUserRelation = (data) => {
  return service({
    url: '/toolsUserRelation/createUserRelation',
    method: 'post',
    data
  })
}

export const findUserRelation = (params) => {
  return service({
    url: '/toolsUserRelation/findUserRelation',
    method: 'get',
    params
  })
}

export const getUserRelationList = (params) => {
  return service({
    url: '/toolsUserRelation/getUserRelationList',
    method: 'get',
    params
  })
}

export const deleteUserRelation = (params) => {
  return service({
    url: '/toolsUserRelation/deleteUserRelation',
    method: 'delete',
    params
  })
}

export const getUserIdsByEnvironment = (params) => {
  return service({
    url: '/toolsUserRelation/getUserIdsByEnvironment',
    method: 'get',
    params
  })
}
