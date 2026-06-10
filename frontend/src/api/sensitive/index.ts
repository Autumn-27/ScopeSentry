import request from '@/axios'
import type { sensitiveData, sensitiveRespData } from './types'

interface sensitiveDataResponse {
  list: sensitiveData[]
  total: number
}
export const getSensitiveDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<sensitiveDataResponse>> => {
  return request.post({ url: '/api/sensitive/data', data: { search, pageIndex, pageSize } })
}

export const addSensitiveDataApi = (
  name: string,
  regular: string,
  color: string,
  state: boolean
): Promise<IResponse<sensitiveRespData>> => {
  return request.post({ url: '/api/sensitive/add', data: { name, regular, color, state } })
}

export const updateSensitiveDataApi = (
  id: string,
  name: string,
  regular: string,
  color: string,
  state: boolean
): Promise<IResponse<sensitiveRespData>> => {
  return request.post({ url: '/api/sensitive/update', data: { id, name, regular, color, state } })
}

export const deleteSensitiveDataApi = (ids: string[]): Promise<IResponse<sensitiveRespData>> => {
  return request.post({ url: '/api/sensitive/delete', data: { ids } })
}

export const updateStateSensitiveDataApi = (
  ids: string[],
  state: boolean
): Promise<IResponse<sensitiveRespData>> => {
  return request.post({ url: '/api/sensitive/update/state', data: { ids, state } })
}
