import request from '@/axios'
import type { dictRespData, fileData, portDictData } from './types'
import { commonRespData } from '../scommon/types'

interface portDictDataResponse {
  list: portDictData[]
  total: number
}
export const getPortDictDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<portDictDataResponse>> => {
  return request.post({ url: '/api/dictionary/port/data', data: { search, pageIndex, pageSize } })
}

export const deletePortDictDataApi = (ids: string[]): Promise<IResponse<dictRespData>> => {
  return request.post({ url: '/api/dictionary/port/delete', data: { ids } })
}

export const upgradePortDictDataApi = (
  id: string,
  name: string,
  value: string
): Promise<IResponse<dictRespData>> => {
  return request.post({ url: '/api/dictionary/port/upgrade', data: { id, name, value } })
}
export const addPortDictDataApi = (
  name: string,
  value: String
): Promise<IResponse<dictRespData>> => {
  return request.post({ url: '/api/dictionary/port/add', data: { name, value } })
}

interface fileDataResponse {
  list: fileData[]
}

export const getManagetListApi = (): Promise<IResponse<fileDataResponse>> => {
  return request.get({ url: '/api/dictionary/manage/list' })
}

export const createDictApi = (formData: FormData): Promise<IResponse<commonRespData>> => {
  // 发送 POST 请求
  return request.post({
    url: '/api/dictionary/manage/create',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    } as any
  })
}

export const downloadDictApi = (id: string) => {
  return request.get({ url: '/api/dictionary/manage/download?id=' + id, responseType: 'blob' })
}

export const deleteDictApi = (ids: string[]): Promise<IResponse<dictRespData>> => {
  return request.post({ url: '/api/dictionary/manage/delete', data: { ids } })
}
