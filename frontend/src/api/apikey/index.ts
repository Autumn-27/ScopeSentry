import request from '@/axios'
import type { ApiKeyData, CreateApiKeyResponse } from './types'
import { commonRespData } from '../scommon/types'

interface ApiKeyListResponse {
  list: ApiKeyData[]
}

export const getApiKeyListApi = (): Promise<IResponse<ApiKeyListResponse>> => {
  return request.get({ url: '/api/apikey/list' })
}

export const createApiKeyApi = (name: string): Promise<IResponse<CreateApiKeyResponse>> => {
  return request.post({ url: '/api/apikey/create', data: { name } })
}

export const deleteApiKeyApi = (id: string): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/apikey/delete', data: { id } })
}
