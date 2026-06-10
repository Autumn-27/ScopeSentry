export type TaskData = {
  id: string
  name: string
  taskNum: string
  progress: string
  creatTime: string
  endTime: string
}

export type taskRespData = {
  code: number
  message: string
}

export type TaskDetail = {
  name: string
  target: string
  ignore: string
  node: string[]
  allNode: boolean
  scheduledTasks: boolean
  hour: number
  duplicates: string
  template: string
  search: string
  filter: Record<string, any>
  targetNumber: number
  targetIds: string[]
  project: string[]
  targetSource: string
  day: number
  minute: number
  week: number
  bindProject: string | null
  cycleType: string
}

export type TaskProgessInfo = {
  subdomain: string[]
  subdomainTakeover: string[]
  portScan: string[]
  assetMapping: string[]
  urlScan: string[]
  crawler: string[]
  dirScan: string[]
  vulnerability: string[]
  all: string[]
}

export type ScheduledTaskData = {
  ID: string
  name: string
  taskType: string
  lastTime: string
  nextTime: string
  state: string
}

export type ScanTemplateData = {
  ID: string
  name: string
  taskNum: string
  progress: string
  creatTime: string
  endTime: string
}

export type TemplateData = {
  ID: string
  name: string
}

export type TemplateDetail = {
  id: string
  name: string
  TargetHandler: string[]
  Parameters: {
    TargetHandler: Record<string, string>
    SubdomainScan: Record<string, string>
    SubdomainSecurity: Record<string, string>
    PortScanPreparation: Record<string, string>
    PortScan: Record<string, string>
    PortFingerprint: Record<string, string>
    AssetMapping: Record<string, string>
    AssetHandle: Record<string, string>
    URLScan: Record<string, string>
    WebCrawler: Record<string, string>
    URLSecurity: Record<string, string>
    DirScan: Record<string, string>
    VulnerabilityScan: Record<string, string>
  }
  SubdomainScan: string[]
  SubdomainSecurity: string[]
  PortScanPreparation: string[]
  PortScan: string[]
  PortFingerprint: string[]
  AssetMapping: string[]
  AssetHandle: string[]
  URLScan: string[]
  WebCrawler: string[]
  URLSecurity: string[]
  DirScan: string[]
  VulnerabilityScan: string[]
  vullist: string[]
}
