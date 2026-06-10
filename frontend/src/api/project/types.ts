export type ProjectData = {
  id: string
  name: string
  logo: string
  AssetCount: number
  tag: string
}

export type projectRespData = {
  code: string
  message: string
}

export type projectContent = {
  name: string
  tag: string
  target: string
  ignore: string
  logo: string
  scheduledTasks: boolean
  allNode: boolean
  node: string[]
  duplicates: string
  hour: number
  template: string
}
