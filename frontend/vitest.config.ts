import { defineConfig } from 'vitest/config'
import { resolve } from 'node:path'

// 前端单测配置：jsdom 提供 localStorage/window，别名与 vite.config.ts 保持一致
export default defineConfig({
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  test: {
    environment: 'jsdom',
    include: ['src/**/*.spec.ts'],
    setupFiles: ['./src/test/setup.ts'],
    clearMocks: true
  }
})
