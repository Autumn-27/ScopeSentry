import { SUCCESS_CODE } from '@/constants'
import { MockMethod } from 'vite-plugin-mock'

const timeout = 1000

export default [
  {
    url: '/api/task/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 50,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 60,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 70,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 80,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 100,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
        },
        {
          TaskName: 'test',
          TaskCount: 500,
          TaskProgress: 30,
          CreateTime: '2023-12-22',
          EndTime: '2023-12-12'
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
