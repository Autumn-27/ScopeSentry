import { describe, expect, it } from 'vitest'
import { isArray, isBoolean, isEmpty, isFunction, isImgPath, isNumber, isString, isUrl } from './is'

describe('is helpers', () => {
  it('detects basic types', () => {
    expect(isString('a')).toBe(true)
    expect(isString(1)).toBe(false)
    expect(isNumber(1)).toBe(true)
    expect(isBoolean(false)).toBe(true)
    expect(isArray([])).toBe(true)
    expect(isFunction(() => 1)).toBe(true)
  })

  it('isEmpty handles the shapes we use', () => {
    expect(isEmpty('')).toBe(true)
    expect(isEmpty('x')).toBe(false)
    expect(isEmpty([])).toBe(true)
    expect(isEmpty([1])).toBe(false)
    expect(isEmpty({})).toBe(true)
    expect(isEmpty({ a: 1 })).toBe(false)
    expect(isEmpty(new Map())).toBe(true)
    expect(isEmpty(new Set([1]))).toBe(false)
  })

  it('isUrl accepts real urls only', () => {
    expect(isUrl('https://example.com/a?b=1')).toBe(true)
    // 裸域名不匹配（实现要求 http(s) 前缀或 www./user@ 前缀），记录现状
    expect(isUrl('example.com')).toBe(false)
    expect(isUrl('www.example.com')).toBe(true)
    expect(isUrl('not a url')).toBe(false)
  })

  it('isImgPath matches image links with an extension', () => {
    expect(isImgPath('https://cdn.example.com/a.png')).toBe(true)
    expect(isImgPath('data:image/png;base64,xxx.png')).toBe(true)
    expect(isImgPath('https://cdn.example.com/a.txt')).toBe(false)
    // 无扩展名的 data URI 不算图片路径（与实现一致）
    expect(isImgPath('data:image/png;base64,xxx')).toBe(false)
  })
})
