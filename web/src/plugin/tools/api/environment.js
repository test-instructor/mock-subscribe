import service from '@/utils/request'

export const createEnvironment = (data) => {
  return service({
    url: '/toolsEnvironment/createEnvironment',
    method: 'post',
    data
  })
}

export const updateEnvironment = (data) => {
  return service({
    url: '/toolsEnvironment/updateEnvironment',
    method: 'put',
    data
  })
}

export const deleteEnvironment = (params) => {
  return service({
    url: '/toolsEnvironment/deleteEnvironment',
    method: 'delete',
    params
  })
}

export const findEnvironment = (params) => {
  return service({
    url: '/toolsEnvironment/findEnvironment',
    method: 'get',
    params
  })
}

export const getEnvironmentList = (params) => {
  return service({
    url: '/toolsEnvironment/getEnvironmentList',
    method: 'get',
    params
  })
}
