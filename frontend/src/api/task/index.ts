import request from '@/axios'
import type {
  TaskData,
  taskRespData,
  TaskDetail,
  TaskProgessInfo,
  ScheduledTaskData,
  TemplateData,
  TemplateDetail
} from './types'
import type { commonRespData } from '../scommon/types'

interface TaskDataResponse {
  list: TaskData[]
  total: number
}
export const getTaskDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<TaskDataResponse>> => {
  return request.post({ url: '/api/task/', data: { search, pageIndex, pageSize } })
}

export const stopTaskApi = (ids: string[]): Promise<IResponse<TaskDataResponse>> => {
  return request.post({ url: '/api/task/stop', data: { ids } })
}

export const starTaskApi = (ids: string[]): Promise<IResponse<TaskDataResponse>> => {
  return request.post({ url: '/api/task/start', data: { ids } })
}

export const addTaskApi = (
  name: string,
  target: string,
  ignore: string,
  node: string[],
  allNode: boolean,
  duplicates: string,
  scheduledTasks: boolean,
  hour: number,
  template: string,
  targetTp: string,
  search: string,
  filter: Record<string, any>,
  targetNumber: number,
  targetIds: string[],
  project: string[],
  targetSource: string,
  day: number,
  minute: number,
  week: number,
  bindProject: string | null,
  cycleType: string
): Promise<IResponse<taskRespData>> => {
  return request.post({
    url: '/api/task/add',
    data: {
      name,
      target,
      ignore,
      node,
      allNode,
      duplicates,
      scheduledTasks,
      hour,
      template,
      targetTp,
      search,
      filter,
      targetNumber,
      targetIds,
      project,
      targetSource,
      day,
      minute,
      week,
      bindProject,
      cycleType
    }
  })
}

export const addScheduledTaskApi = (
  name: string,
  target: string,
  ignore: string,
  node: string[],
  allNode: boolean,
  duplicates: string,
  scheduledTasks: boolean,
  hour: number,
  template: string,
  targetTp: string,
  search: string,
  filter: Record<string, any>,
  targetNumber: number,
  targetIds: string[],
  project: string[],
  targetSource: string,
  day: number,
  minute: number,
  week: number,
  bindProject: string | null,
  cycleType: string
): Promise<IResponse<taskRespData>> => {
  return request.post({
    url: '/api/task/scheduled/add',
    data: {
      name,
      target,
      ignore,
      node,
      allNode,
      duplicates,
      scheduledTasks,
      hour,
      template,
      targetTp,
      search,
      filter,
      targetNumber,
      targetIds,
      project,
      targetSource,
      day,
      minute,
      week,
      bindProject,
      cycleType
    }
  })
}

export const updateScheduleApi = (
  id: string,
  name: string,
  target: string,
  ignore: string,
  node: string[],
  allNode: boolean,
  duplicates: string,
  scheduledTasks: boolean,
  hour: number,
  template: string,
  targetTp: string,
  search: string,
  filter: Record<string, any>,
  targetNumber: number,
  targetIds: string[],
  project: string[],
  targetSource: string,
  day: number,
  minute: number,
  week: number,
  bindProject: string | null,
  cycleType: string
): Promise<IResponse<taskRespData>> => {
  return request.post({
    url: '/api/task/scheduled/update',
    data: {
      id,
      name,
      target,
      ignore,
      node,
      allNode,
      duplicates,
      scheduledTasks,
      hour,
      template,
      targetTp,
      search,
      filter,
      targetNumber,
      targetIds,
      project,
      targetSource,
      day,
      minute,
      week,
      bindProject,
      cycleType
    }
  })
}

export const getTaskDetailApi = (id: string): Promise<IResponse<TaskDetail>> => {
  return request.post({ url: '/api/task/detail', data: { id } })
}

export const getScheduleDetailApi = (id: string): Promise<IResponse<TaskDetail>> => {
  return request.post({ url: '/api/task/scheduled/detail', data: { id } })
}

export const deleteTaskApi = (ids: string[], delA: boolean): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/delete', data: { ids, delA } })
}

export const retestTaskApi = (id: string): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/retest', data: { id } })
}

interface TaskProgessInforesp {
  list: TaskProgessInfo[]
  total: number
}
export const getTaskProgressApi = (
  id: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<TaskProgessInforesp>> => {
  return request.post({ url: '/api/task/progress/info', data: { id, pageIndex, pageSize } })
}

interface ScheduledTaskDataResponse {
  list: ScheduledTaskData[]
  total: number
}

export const getScheduledTaskDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<ScheduledTaskDataResponse>> => {
  return request.post({ url: '/api/task/scheduled', data: { search, pageIndex, pageSize } })
}

export const taskRunApi = (id: string): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/scheduled/run', data: { id } })
}

export const scheduledDeleteTaskApi = (ids: string[]): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/scheduled/delete', data: { ids } })
}

export const getScheduledTaskPageMonitDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<ScheduledTaskDataResponse>> => {
  return request.post({
    url: '/api/task/scheduled/pagemonit/data',
    data: { search, pageIndex, pageSize }
  })
}

export const deleteScheduledTaskPageMonitApi = (
  ids: string[]
): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/scheduled/pagemonit/delete', data: { ids } })
}

export const updateScheduledTaskPageMonitApi = (
  hour: number,
  node: string[],
  allNode: boolean,
  scheduledTasks: boolean
): Promise<IResponse<taskRespData>> => {
  return request.post({
    url: '/api/task/scheduled/pagemonit/update',
    data: {
      hour,
      node,
      allNode,
      scheduledTasks
    }
  })
}

export const addScheduledTaskPageMonitApi = (url: string): Promise<IResponse<taskRespData>> => {
  return request.post({
    url: '/api/task/scheduled/pagemonit/add',
    data: {
      url
    }
  })
}

interface TaskDataResponse {
  list: TaskData[]
  total: number
}
export const getTemplateDataApi = (
  search: string,
  pageIndex: number,
  pageSize: number
): Promise<IResponse<TaskDataResponse>> => {
  return request.post({ url: '/api/task/template', data: { search, pageIndex, pageSize } })
}

export const getTemplateDetailApi = (id: string): Promise<IResponse<TemplateDetail>> => {
  return request.post({ url: '/api/task/template/detail', data: { id } })
}

export const saveTemplateDetailApi = (
  result: Record<string, any>,
  id: string
): Promise<IResponse<TemplateDetail>> => {
  return request.post({ url: '/api/task/template/save', data: { result, id } })
}

export const deleteTemplateDetailApi = (ids: string[]): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/template/delete', data: { ids } })
}

export const syancProjectApi = (
  ids: string[],
  option: string,
  project: string,
  tag: string,
  name: string
): Promise<IResponse<commonRespData>> => {
  return request.post({ url: '/api/task/sync', data: { ids, option, project, tag, name } })
}

interface TaskNameData {
  id: string
  name: string
}

export const getTaskNamesApi = (): Promise<IResponse<TaskNameData[]>> => {
  return request.get({ url: '/api/task/names' })
}
