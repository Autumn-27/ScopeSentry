import request from '@/axios'
import type { vulResultData } from './types'

interface vulResultDataResponse {
  list: vulResultData[]
  total: number
}

interface vulDetailResponse {
  res: string
  req: string
}
export const getVulResultDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<vulResultDataResponse>> => {
  return request.post({
    url: '/api/assets/vulnerability',
    data: { search, pageIndex, pageSize, filter }
  })
}

// eslint-disable-next-line prettier/prettier
export const getVulDetailApi = (
  hash: string
): Promise<IResponse<vulDetailResponse>> => {
  return request.post({
    url: '/api/assets/vulnerability/detail',
    data: { hash }
  })
}
