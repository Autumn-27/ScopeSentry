export type DashboardTotalTypes = {
  asetCount: number
  subdomainCount: number
  sensitiveCount: number
  urlCount: number
  vulnerabilityCount: number
}

export type UserAccessSource = {
  value: number
  name: string
}

export type WeeklyUserActivity = {
  value: number
  name: string
}

export type MonthlySales = {
  name: string
  estimate: number
  actual: number
}

export type VersionData = {
  name: string
  cversion: string
  lversion: string
  msg: string
}
