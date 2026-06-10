import request from '@/axios'
import type {
  DashboardTotalTypes,
  UserAccessSource,
  WeeklyUserActivity,
  MonthlySales,
  VersionData
} from './types'

interface DashboardTotalTypesResponse {
  data: DashboardTotalTypes
}

export const getAssetStatisticsApi = (): Promise<IResponse<DashboardTotalTypesResponse[]>> => {
  return request.get({ url: '/api/assets/statistics' })
}

export const getUserAccessSourceApi = (): Promise<IResponse<UserAccessSource[]>> => {
  return request.get({ url: '/mock/analysis/userAccessSource' })
}

export const getWeeklyUserActivityApi = (): Promise<IResponse<WeeklyUserActivity[]>> => {
  return request.get({ url: '/mock/analysis/weeklyUserActivity' })
}

export const getMonthlySalesApi = (): Promise<IResponse<MonthlySales[]>> => {
  return request.get({ url: '/mock/analysis/monthlySales' })
}

interface VersionDataResponse {
  list: VersionData[]
}
export const getVersionDataApi = (): Promise<IResponse<VersionDataResponse>> => {
  return request.get({ url: '/api/system/version' })
}

export const UPDATEsYSTEMApi = (
  server: string,
  scan: string,
  key: string
): Promise<IResponse<VersionDataResponse>> => {
  return request.post({ url: '/api/system/update', data: { server, scan, key } })
}
