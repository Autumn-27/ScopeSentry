import { SUCCESS_CODE } from '@/constants'
import { MockMethod } from 'vite-plugin-mock'

const timeout = 1000

export default [
  {
    url: '/api/asset/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.net',
          IP: '192.168.1.2',
          Port: 8080,
          Title: 'Example Title 2',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY29jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Service: 'http',
          Status: 302,
          Time: '2023-01-02T14:30:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.net',
          IP: '192.168.1.2',
          Port: 8080,
          Title: 'Example Title 2',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY29jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Service: 'http',
          Status: 404,
          Time: '2023-01-02T14:30:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        },
        {
          ID: '2232265dwadsdwadwa',
          Domain: 'example.com',
          IP: '333.333.333.333',
          Port: 80,
          Title: 'Example Title 1',
          Banner: `HTTP/1.1 200 OK\nAccess-Control-Allow-Origin: *\nContent-Type: application/javascript\nCache-Control: max-age=31536000,immutable\nEtag: W/"2744-RiWvKoVG6OPYG/IKC4v+QzY2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx9jo"\nDate: Sun, 24 Dec 2023 13:20:01 GMT\nConnection: keep-alive\nKeep-Alive: timeout=5\nContent-Length: 10052`,
          Products: ['Apache', 'PHP', 'Test'],
          Status: 200,
          Service: 'http',
          Time: '2023-01-01T12:00:00Z'
        }
      ]

      const total = 400
      return {
        code: SUCCESS_CODE,
        data: {
          list: mockData,
          total: total
        }
      }
    }
  },
  {
    url: '/api/asset/statistics',
    method: 'post',
    timeout,
    response: () => {
      return {
        code: SUCCESS_CODE,
        data: {
          Port: [
            { value: 8080, number: 10 },
            { value: 80, number: 5 }
          ],
          Service: [
            { value: 'HTTP', number: 15 },
            { value: 'HTTPS', number: 8 }
          ],
          Product: [
            { value: 'Web Server', number: 20 },
            { value: 'Load Balancer', number: 5 }
          ],
          Icon: [
            {
              value:
                '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',
              icon_hash: '11561481615',
              number: 25
            },
            {
              value:
                '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',
              icon_hash: '11561481615',
              number: 25
            }
          ]
        }
      }
    }
  },
  {
    url: '/api/asset/detail',
    method: 'post',
    timeout,
    response: () => {
      return {
        code: SUCCESS_CODE,
        data: {
          ID: '12315611949816515',
          Domain: 'google.com',
          IP: '127.0.0.1',
          URL: 'http://google.com/test',
          Port: 80,
          Title: 'google',
          Status: 200,
          FaviconHash: '11111111',
          Jarm: '11111111',
          Time: '2023-12-24T18:06:49+08:00',
          Products: ['a', 'b', 'c'],
          Service: 'http',
          RawHeaders:
            'HTTP/1.1 200 OK\r\nConnection: close\r\nTransfer-Encoding: chunked\r\nContent-Type: text/html\r\nDate: Sun, 24 Dec 2023 10:06:48 GMT\r\nEtag: W/"64f6a44a-1bcad"\r\nLast-Modified: Tue, 05 Sep 2023 03:45:14 GMT\r\nServer: nginx\r\nStrict-Transport-Security: max-age=31536000\r\nVary: Accept-Encoding\r\n\r\n',
          TLSData:
            '{"host":"xxxxxxxxxxxxxxxx","port":"443","probe_status":true,"tls_version":"tls13","cipher":"TLS_AES_256_GCM_SHA384","mismatched":true,"not_before":"2023-02-16T00:00:00Z","not_after":"2024-02-16T23:59:59Z","subject_dn":"CN=xxxxxxxx","subject_cn":"xxxxxxxxxxxx","subject_an":["xxxxxxxxxxxxx","xxxxxxxxxxxxxxxxxx"],"issuer_dn":"CN=xxxxxxxxxxxxxxxxxxxxxxxxx\\, Inc., C=CN","issuer_cn":"TrustAsia RSA DV TLS CA G2","issuer_org":["TrustAsia Technologies, Inc."],"fingerprint_hash":{"md5":"xxxxxxxxxxxxxxxxxxxxxxxxxx","sha1":"xxxxxxxxxxxxxxxxxxxxxxxxxxx","sha256":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},"tls_connection":"ctls","sni":"xxxxxxxxxxxxxxxxx"}',
          Hash: '{"body_md5":"xxxxxx","body_mmh3":"-xxxxx","body_sha256":"xxxxxxx","body_simhash":"xxxxxxxx","header_md5":"xxxxxxxxxxx","header_mmh3":"xxxxxxxxxxxxxxx","header_sha256":"xxxxxxxxxxxxxxxx","header_simhash":"xxxxxxxxxxxxxxxxxxxxxx"}',
          ResponseBody:
            'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
        }
      }
    }
  },
  {
    url: '/api/subdomain/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123', '123.1230123'],
          IP: ['123.123.123', '123.1230123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123', '123.1230123'],
          IP: ['123.123.123', '123.1230123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123', '123.1230123'],
          IP: ['123.123.123', '123.1230123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123', '123.1230123'],
          IP: ['123.123.123', '123.1230123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123', '123.1230123'],
          IP: ['123.123.123', '123.1230123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123'],
          IP: ['123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123'],
          IP: ['123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123'],
          IP: ['123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'A',
          RecordValue: ['123.123.123'],
          IP: ['123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'CNAME',
          RecordValue: ['xxxxxxxxxxxxxxxxx'],
          IP: ['123.123.123', '123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'CNAME',
          RecordValue: ['xxxxxxxxxxxxxxxxx'],
          IP: ['123.123.123', '123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'CNAME',
          RecordValue: ['xxxxxxxxxxxxxxxxx'],
          IP: ['123.123.123', '123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          SubDomain: 'xxxxx.com',
          RecordType: 'CNAME',
          RecordValue: ['xxxxxxxxxxxxxxxxx'],
          IP: ['123.123.123', '123.123.123'],
          Time: '2023-12-24T18:06:49+08:00'
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
  },
  {
    url: '/api/url/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'https://xxxxxxxxxxxx/lhG6/uoxo7agf',
          Source: 'body',
          Type: 'form',
          Input: 'https://xxxxxxxxxxxx/',
          Time: '2023-12-24T18:06:49+08:00'
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
  },
  {
    url: '/api/crawler/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'POST',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Method: 'GET',
          URL: 'https:/xxxxxxxxxxxxx/123',
          GetParameter: 'a=123&b=45&c=dwasdwa',
          PostParameter: 'a=123&b=45&c=dwasdwa',
          Time: '2023-12-24T18:06:49+08:00'
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
  },
  {
    url: '/api/sensitive/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          Info: 'apikey=xxxxx',
          Name: 'APIKEY',
          URL: 'http://xxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
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
  },
  {
    url: '/api/dirscan/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          Title: 'xxxx',
          Status: 200,
          Length: 123456,
          Time: '2023-12-24T18:06:49+08:00'
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
  },
  {
    url: '/api/pagemonitoring/data',
    method: 'post',
    timeout,
    response: () => {
      const mockData = [
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
        },
        {
          ID: '12315611949816515',
          URL: 'http://xxxxxxxxxxxxx',
          OldResponseBodyMD5: 'xxxx',
          CurrentResponseBodyMD5: 'xxxxxxxxxxxxxxxxxxxxx',
          Time: '2023-12-24T18:06:49+08:00'
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
