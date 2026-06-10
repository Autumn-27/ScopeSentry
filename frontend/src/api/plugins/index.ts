import request from '@/axios'
import axios from 'axios'
import type { LogRespData, pluginData, RemotePluginResponse } from './types'
import { commonRespData } from '../scommon/types'

interface pluginDataResponse {
  list: pluginData[]
  total: number
}

interface PluginDeleteItem {
  hash: string
  module: string
}

export const getPluginDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  type?: string
): Promise<IResponse<pluginDataResponse>> => {
  return request.post({ url: '/api/plugin', data: { search, pageIndex, pageSize, type } })
}

export const getPluginDetailApi = (id: string): Promise<IResponse<pluginData>> => {
  return request.post({ url: '/api/plugin/detail', data: { id } })
}

export const savePluginDataApi = (
  id: string,
  name: string,
  version: string,
  module: string,
  parameter: string,
  help: string,
  introduction: string,
  source: string,
  key: string,
  parameterList?: string,
  type?: string,
  hash?: string
): Promise<IResponse<commonRespData>> => {
  return request.post({
    url: '/api/plugin/save',
    data: {
      id,
      name,
      version,
      module,
      parameter,
      help,
      introduction,
      source,
      key,
      parameterList,
      type,
      hash
    }
  })
}

export const deletePluginDataApi = (
  items: PluginDeleteItem[]
): Promise<IResponse<commonRespData>> => {
  return request.post({
    url: '/api/plugin/delete',
    data: { data: items }
  })
}

export const checkKeyApi = (key: string): Promise<IResponse<LogRespData>> => {
  return request.post({ url: '/api/plugin/key/check', data: { key } })
}

export const getPluginLogApi = (
  module: string,
  hash: string,
  type?: string
): Promise<IResponse<LogRespData>> => {
  return request.post({ url: '/api/plugin/log', data: { module, hash, type } })
}

export const cleanPluginLogApi = (
  module: string,
  hash: string,
  type?: string
): Promise<IResponse<LogRespData>> => {
  return request.post({ url: '/api/plugin/log/clean', data: { module, hash, type } })
}

export const cleanAllPluginLogApi = (): Promise<IResponse<LogRespData>> => {
  return request.post({ url: '/api/plugin/log/clean/all' })
}

export const getPluginDataByModuleApi = (
  module: string
): Promise<IResponse<pluginDataResponse>> => {
  return request.post({ url: '/api/plugin/module', data: { module } })
}

export const reInstallPluginApi = (
  node: string,
  hash: string,
  module: string
): Promise<IResponse<pluginDataResponse>> => {
  return request.post({ url: '/api/plugin/reinstall', data: { node, hash, module } })
}

export const reCheckPluginApi = (
  node: string,
  hash: string,
  module: string
): Promise<IResponse<pluginDataResponse>> => {
  return request.post({ url: '/api/plugin/recheck', data: { node, hash, module } })
}

export const uninstallPluginApi = (
  node: string,
  hash: string,
  module: string
): Promise<IResponse<pluginDataResponse>> => {
  return request.post({ url: '/api/plugin/uninstall', data: { node, hash, module } })
}

// 获取本地已安装的插件列表
export const getLocalPluginListApi = (): Promise<IResponse<pluginDataResponse>> => {
  return request.post({
    url: '/api/plugin',
    data: { search: '', pageIndex: 1, pageSize: 1000, type: 'all' }
  })
}

// 获取远程插件市场数据
export const getRemotePluginMarketApi = (): Promise<any> => {
  // 使用原生 axios 直接调用远程 API，避免添加 Authorization header
  return axios
    .post(
      'https://api.scope-sentry.top/api/common/plugin/search',
      {
        name: '',
        priceStatus: 2,
        module: 'All',
        tags: '',
        page: 1,
        size: 200
      },
      {
        headers: {
          'Content-Type': 'application/json'
        }
      }
    )
    .then((res) => res.data)
}

// 获取插件导出数据
export const getPluginExportDataApi = (hash: string, token?: string): Promise<any> => {
  // 使用原生 axios 直接调用远程 API
  let url = `https://api.scope-sentry.top/api/common/plugin/export-data/${hash}`
  if (token) {
    url += `?token=${token}`
  }
  return axios
    .get(url, {
      headers: {
        'Content-Type': 'application/json'
      }
    })
    .then((res) => res.data)
}

// 导入插件
export const importPluginApi = (
  json: string,
  source: string,
  isSystem: boolean,
  key: string
): Promise<IResponse<commonRespData>> => {
  return request.post({
    url: '/api/plugin/import/data',
    data: {
      json,
      source,
      isSystem,
      key
    }
  })
}

// 更新插件状态
export const updatePluginStatusApi = (
  id: string,
  status: boolean
): Promise<IResponse<commonRespData>> => {
  return request.post({
    url: '/api/plugin/status',
    data: {
      id,
      status
    }
  })
}

// 运行插件一次
export const runPluginOnceApi = (hash: string): Promise<IResponse<commonRespData>> => {
  return request.post({
    url: '/api/plugin/run',
    data: {
      hash
    }
  })
}
