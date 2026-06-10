import { SUCCESS_CODE } from '@/constants'
import { MockMethod } from 'vite-plugin-mock'

const timeout = 1000

export default [
  {
    url: '/api/project/data',
    method: 'get',
    timeout,
    response: () => {
      return {
        code: SUCCESS_CODE,
        data: [
          {
            ID: '12112313',
            ProjectName: 'Test',
            Logo: 'https://g.csdnimg.cn/static/logo/favicon32.ico',
            AssetCount: 50,
            TagName: 'Hackone'
          },
          {
            ID: '12112313',
            ProjectName: 'Test',
            Logo: 'https://element-plus.org/safari-pinned-tab.svg',
            AssetCount: 50,
            TagName: 'Bugcrowd'
          }
        ]
      }
    }
  }
] as MockMethod[]
