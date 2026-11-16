// 获取传入的值的类型
const getValueType = (value: any) => {
  const type = Object.prototype.toString.call(value)
  return type.slice(8, -1)
}

export const useStorage = (type: 'sessionStorage' | 'localStorage' = 'sessionStorage') => {
  const setStorage = (key: string, value: any) => {
    const valueType = getValueType(value)
    window[type].setItem(key, JSON.stringify({ type: valueType, value }))
  }

  const getStorage = (key: string) => {
    // localStorage 可能被禁用（隐私模式）或写入损坏，读不出来按“没有”处理，
    // 不能让根组件 setup 抛错导致白屏
    try {
      const value = window[type].getItem(key)
      if (!value) return value
      const parsed = JSON.parse(value)
      if (parsed && typeof parsed === 'object' && 'value' in parsed) {
        return parsed.value
      }
      // 兼容被手工写成裸 JSON 的场景
      return parsed
    } catch (e) {
      return null
    }
  }

  const removeStorage = (key: string) => {
    window[type].removeItem(key)
  }

  const clear = (excludes?: string[]) => {
    // 获取排除项
    const keys = Object.keys(window[type])
    const defaultExcludes = ['dynamicRouter', 'serverDynamicRouter']
    const excludesArr = excludes ? [...excludes, ...defaultExcludes] : defaultExcludes
    const excludesKeys = excludesArr ? keys.filter((key) => !excludesArr.includes(key)) : keys
    // 排除项不清除
    excludesKeys.forEach((key) => {
      window[type].removeItem(key)
    })
    // window[type].clear()
  }

  return {
    setStorage,
    getStorage,
    removeStorage,
    clear
  }
}
