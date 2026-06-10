import request from '@/axios'
import {
  AssetData,
  AssetStatistics,
  AssetDetail,
  SubdomainData,
  URLData,
  CrawlerData,
  SensitiveData,
  DirScanData,
  PageMonitoringData,
  SubdomaintakerData,
  SensitiveBody,
  PageMResponse,
  PageMHistory,
  SensitiveNames,
  AssetChangeLog,
  AssetScreenshot
} from './types'
import { commonRespData } from '../scommon/types'

interface AssetDataResponse {
  list: AssetData[]
  total: number
}
export const getAssetApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<AssetDataResponse>> => {
  return request.post({ url: '/api/assets/asset', data: { search, pageIndex, pageSize, filter } })
}

export const getAssetCardApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<AssetDataResponse>> => {
  return request.post({
    url: '/api/assets/asset/card',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getAssetScreenshotApi = (id: string): Promise<IResponse<AssetScreenshot>> => {
  return request.post({ url: '/api/asset/screenshot', data: { id } })
}

export const getAssetStatisticsPortApi = (
  search: string,
  filter: Record<string, any>
): Promise<IResponse<AssetStatistics>> => {
  return request.post({ url: '/api/assets/statistics/port', data: { search, filter } })
}

export const getAssetStatisticsTypeApi = (
  search: string,
  filter: Record<string, any>
): Promise<IResponse<AssetStatistics>> => {
  return request.post({ url: '/api/assets/statistics/service', data: { search, filter } })
}

export const getAssetStatisticsiconApi = (
  search: string,
  filter: Record<string, any>,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<AssetStatistics>> => {
  return request.post({
    url: '/api/assets/statistics/icon',
    data: { search, filter, pageIndex, pageSize }
  })
}

export const getAssetStatisticsappApi = (
  search: string,
  filter: Record<string, any>
): Promise<IResponse<AssetStatistics>> => {
  return request.post({ url: '/api/assets/statistics/app', data: { search, filter } })
}

export const getAssetStatisticsTitleApi = (
  search: string,
  filter: Record<string, any>
): Promise<IResponse<AssetStatistics>> => {
  return request.post({ url: '/api/asset/statistics/title', data: { search, filter } })
}

export const getAssetDetailApi = (id: string): Promise<IResponse<AssetDetail>> => {
  return request.post({ url: '/api/assets/asset/detail', data: { id } })
}

export const getAssetChangeLogApi = (id: string): Promise<IResponse<AssetChangeLog[]>> => {
  return request.post({ url: '/api/assets/asset/changelog', data: { id } })
}

export const updateStatusApi = (
  id: string,
  tp: string,
  status: number
): Promise<IResponse<AssetChangeLog[]>> => {
  return request.post({ url: '/api/assets/common/update_status', data: { id, tp, status } })
}

interface SubdomainDataResponse {
  list: SubdomainData[]
  total: number
}

export const getSubdomainApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/assets/subdomain',
    data: { search, pageIndex, pageSize, filter }
  })
}

interface URLDataResponse {
  list: URLData[]
  total: number
}

export const getURLApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>,
  sort: Record<string, any>
): Promise<IResponse<URLDataResponse>> => {
  return request.post({
    url: '/api/assets/url',
    data: { search, pageIndex, pageSize, filter, sort }
  })
}

interface CrawlerDataResponse {
  list: CrawlerData[]
  total: number
}

export const getCrawlerApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<CrawlerDataResponse>> => {
  return request.post({ url: '/api/assets/crawler', data: { search, pageIndex, pageSize, filter } })
}

interface SensitiveDataResponse {
  list: SensitiveData[]
  total: number
}

interface SensitiveDataNumberResponse {
  all: number
  total: number
}

export const getSensitiveResultApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SensitiveDataResponse>> => {
  return request.post({
    url: '/api/assets/sensitive',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getSensitiveResultNumberApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SensitiveDataNumberResponse>> => {
  return request.post({
    url: '/api/assets/sensitive/number',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getSensitiveResultBodyApi = (id: string): Promise<IResponse<SensitiveBody>> => {
  return request.post({ url: '/api/assets/sensitive/body', data: { id } })
}

interface SensitiveNamesResponse {
  list: SensitiveNames[]
}

export const getSensitiveNamesApi = (
  search: string,
  filter: Record<string, any>
): Promise<IResponse<SensitiveNamesResponse>> => {
  return request.post({
    url: '/api/assets/sensitive/names',
    data: { search, filter }
  })
}

interface SensitiveInfoResponse {
  list: string[]
}

export const getSensitiveInfoApi = (
  sid: string,
  search: string,
  filter: Record<string, any>
): Promise<IResponse<SensitiveInfoResponse>> => {
  return request.post({
    url: '/api/assets/sensitive/info',
    data: { sid, search, filter }
  })
}

interface DirScanDataResponse {
  list: DirScanData[]
  total: number
}
export const getDirScanApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>,
  sort: Record<string, any>
): Promise<IResponse<DirScanDataResponse>> => {
  return request.post({
    url: '/api/assets/dirscan',
    data: { search, pageIndex, pageSize, filter, sort }
  })
}

interface PageMonitoringDataResponse {
  list: PageMonitoringData[]
  total: number
}
export const getPageMonitoringApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<PageMonitoringDataResponse>> => {
  return request.post({
    url: '/api/assets/page-monitoring',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getPageMonitoringResponseApi = (
  id: string,
  flag: string
): Promise<IResponse<PageMResponse>> => {
  return request.post({
    url: '/api/page/monitoring/response',
    data: { id, flag }
  })
}

export const getPageMonitoringDiffApi = (id: string): Promise<IResponse<PageMHistory>> => {
  return request.post({
    url: '/api/assets/page-monitoring/diff',
    data: { id }
  })
}

interface SubdomaintakerDataResponse {
  list: SubdomaintakerData[]
  total: number
}
export const getSubdomaintakerApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SubdomaintakerDataResponse>> => {
  return request.post({
    url: '/api/assets/subdomain/taker',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const delDataApi = (ids: string[], index: string): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/assets/common/delete', data: { ids, index } })
}

export const addTagApi = (
  id: string,
  tp: string,
  tag: string
): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/assets/common/add_tag', data: { id, tp, tag } })
}

export const deleteTagApi = (
  id: string,
  tp: string,
  tag: string
): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/assets/common/delete_tag', data: { id, tp, tag } })
}

interface TotalDataResponse {
  total: number
}
export const totalDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>,
  index: string
): Promise<IResponse<TotalDataResponse>> => {
  return request.post({
    url: '/api/assets/common/total',
    data: { search, pageIndex, pageSize, filter, index }
  })
}

export const getRootDomainApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/assets/root_domain',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getAppApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/assets/app',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getMpApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/assets/mp',
    data: { search, pageIndex, pageSize, filter }
  })
}

export const getIPAssetApi = (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
): Promise<IResponse<SubdomainDataResponse>> => {
  return request.post({
    url: '/api/assets/ip',
    data: { search, pageIndex, pageSize, filter }
  })
}
