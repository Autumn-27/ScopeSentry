import request from '@/axios'

interface projectInfoResponse {
  name: string
  tag: string
  scheduledTasks: boolean
  hour: number
  AssetCount: number
  root_domains: string[]
  next_time: string
}
export const getProjectInfoApi = (id: string): Promise<IResponse<projectInfoResponse>> => {
  return request.post({ url: '/api/project_aggregation/project/info', data: { id } })
}

interface projectInfoResponse {
  subdomainCount: number
  vulCount: number
}
export const getProjectAssetCountApi = (id: string): Promise<IResponse<projectInfoResponse>> => {
  return request.post({ url: '/api/project_aggregation/project/asset/count', data: { id } })
}

interface projectVulLevelInfoResponse {
  _id: string
  count: number
}
export const getProjectVulLevelCountApi = (
  id: string
): Promise<IResponse<projectVulLevelInfoResponse[]>> => {
  return request.post({ url: '/api/project_aggregation/project/vul/statistics', data: { id } })
}

interface projectVulLevelInfoResponse {
  _id: string
  count: number
}

type VulData = {
  url: string
  vulname: string
  level: string
  time: string
  matched: string
}

interface VulDataResponse {
  list: VulData[]
}
export const getProjectVulDataApi = (id: string): Promise<IResponse<VulDataResponse>> => {
  return request.post({ url: '/api/project_aggregation/project/vul/data', data: { id } })
}

export type SubdomainData = {
  id: string
  host: string
  type: string
  value: string[]
  ip: string[]
  time: string
}

interface SubdomainDataResponse {
  list: SubdomainData[]
  total: number
}
export const getProjectSubdomainDataApi = (
  search: string,
  filter: Record<string, any>,
  fq: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/project_aggregation/project/subdomain/data',
    data: { search, filter, fq }
  })
}

export const getProjectPortDataApi = (
  search: string,
  filter: Record<string, any>,
  fq: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/project_aggregation/project/port/data',
    data: { search, filter, fq }
  })
}

export const getProjectServiceDataApi = (
  search: string,
  filter: Record<string, any>,
  fq: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/project_aggregation/project/service/data',
    data: { search, filter, fq }
  })
}
