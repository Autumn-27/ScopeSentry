import { SUCCESS_CODE } from '@/constants'
import { MockMethod } from 'vite-plugin-mock'

const timeout = 1000

export default [
  {
    url: '/api/node/status',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        },
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        },
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        },
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        },
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        },
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        },
        {
          ID: 1,
          NodeName: 'test2',
          TaskCount: '12',
          NodeStatus: '1',
          NodeUsageLoad: 20,
          NodeUsageMemory: 75,
          NodeUsageCpu: 90,
          CreateTime: '2023-12-12'
        }
      ]
      return {
        code: SUCCESS_CODE,
        data: {
          list: mockData,
          total: 400
        }
      }
    }
  }
] as MockMethod[]
