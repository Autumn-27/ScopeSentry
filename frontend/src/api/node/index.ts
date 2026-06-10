import request from '@/axios'
import { NodeData, nodeRespData, nodeLogRespData, pluginInfoData } from './types'
import { commonRespData } from '../scommon/types'

interface NodeDataResponse {
  list: NodeData[]
}

interface NodeDataOnlineResponse {
  list: string[]
}

export const getNodeDataApi = (): Promise<IResponse<NodeDataResponse>> => {
  return request.get({ url: '/api/node' })
}

export const restartNodeApi = (name: string): Promise<IResponse<commonRespData>> => {
  return request.post({
    url: '/api/node/restart',
    data: {
      name
    }
  })
}

export const getNodeDataOnlineApi = (): Promise<IResponse<NodeDataOnlineResponse>> => {
  return request.get({ url: '/api/node/online' })
}

export const updateNodeConfigDataApi = (
  oldName: string,
  name: string,
  ModulesConfig: string,
  state: boolean
): Promise<IResponse<NodeDataResponse>> => {
  return request.post({
    url: '/api/node/config/update',
    data: {
      oldName,
      name,
      ModulesConfig,
      state
    }
  })
}

export const deleteNodeApi = (names: string[]): Promise<IResponse<nodeRespData>> => {
  return request.post({ url: '/api/node/delete', data: { names } })
}

export const getNodeLogApi = (name: string): Promise<IResponse<nodeLogRespData>> => {
  return request.post({ url: '/api/node/log', data: { name } })
}

interface pluginInfoDataresp {
  list: pluginInfoData[]
}
export const getPluginInfoApi = (name: string): Promise<IResponse<pluginInfoDataresp>> => {
  return request.post({ url: '/api/node/plugin', data: { name } })
}
