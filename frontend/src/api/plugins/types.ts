export type pluginData = {
  id: string
  module: string
  version: string
  name: string
  hash: string
  parameter: string
  help: string
  introduction: string
  isSystem: boolean
  source: string
  parameterList?: string
}

export type LogRespData = {
  code: string
  data: string
}

export type RemotePluginData = {
  id: number
  name: string
  module: string
  priceStatus: number
  price: number | null
  hash: string
  introduction: string
  version: string
  createTime: string
  username: string
  isInstalled: boolean
  needUpdate: boolean
  isSystem: boolean
  type?: string // 'server' | 'scan' | '' (空字符串视为 'scan')
}

export interface RemotePluginResponse {
  data: RemotePluginData[]
}
