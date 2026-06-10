import { InputInstance } from 'element-plus'
import { Ref } from 'vue'

export type AssetData = {
  id: string
  domain: string
  host: string
  ip: string
  port: number
  service: string
  title: string
  status: number
  banner: string
  products: string[]
  time: String
  icon: string
  screenshot: string
  type: string
  statuscode: number
  ResponseBodyHash: string
  url: string
}

export type AssetScreenshot = {
  screenshot: string
}

export type AssetStatistics = {
  Port: { value: number; number: number }[]
  Service: { value: string; number: number }[]
  Product: { value: string; number: number }[]
  Icon: { value: string; number: number; icon_hash: string }[]
  Title: { value: string; number: number }[]
}

export type AssetDetail = {
  json: string
}

export type AssetChangeLogField = {
  fieldname: string
  old: string
  new: string
}
export type AssetChangeLog = {
  timestamp: string
  isExpanded: boolean
  change: AssetChangeLogField[]
}

export type SubdomainData = {
  id: string
  host: string
  type: string
  value: string[]
  ip: string[]
  time: string
}

export type URLData = {
  ID: string
  URL: string
  Source: string
  Type: string
  Input: string
  Time: string
}

export type CrawlerData = {
  ID: string
  Method: string
  URL: string
  GetParameter: string
  PostParameter: string
  Time: string
}

export type SensitiveData = {
  ID: string
  url: string
  color: string
  name: String
  time: string
  body: string
  match: string[]
}

export type SensitiveNames = {
  color: string
  name: string
  count: number
}

export type SensitiveBody = {
  body: string
}

export type DirScanData = {
  ID: string
  URL: string
  Title: string
  Status: String
  Length: string
  Time: string
}

export type PageMonitoringData = {
  ID: string
  URL: string
  OldResponseBodyMD5: string
  CurrentResponseBodyMD5: String
  Time: string
}

export type SubdomaintakerData = {
  host: string
  type: string
  value: string
  response: string
}

export type PageMResponse = {
  hash: string
  content: string
}

export type PageMHistory = {
  diff: string[]
}

export type RowState = {
  inputVisible: boolean
  inputValue: string
  inputRef: Ref<InputInstance | null>
}
