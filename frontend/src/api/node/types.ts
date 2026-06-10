export type NodeData = {
  name: string
  running: number
  finished: number
  state: number
  cpuNum: number
  memNum: number
  updateTime: string
  maxTaskNum: string
  urlThread: string
  urlMaxNum: string
}

export type nodeRespData = {
  code: string
  message: string
}

export type nodeLogRespData = {
  code: string
  logs: string
}

export type pluginInfoData = {
  name: string
  install: number
  check: number
}
