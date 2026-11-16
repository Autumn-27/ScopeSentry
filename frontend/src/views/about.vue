<script setup lang="ts">
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n()
import {
  InfoFilled,
  Platform,
  Document,
  Monitor,
  Coffee,
  Money,
  Check
} from '@element-plus/icons-vue'
import { ref } from 'vue'
// 导入图片
import alipayQR from '@/assets/qrcode/zfb.png'
import wechatQR from '@/assets/qrcode/wx.jpg'

// 创建图片数据
const qrCodes = ref([
  {
    name: t('localization.alipay'),
    img: alipayQR
  },
  {
    name: t('localization.wechat'),
    img: wechatQR
  }
])
</script>
<template>
  <div class="about-container p-6">
    <!-- 项目信息卡片 -->
    <el-card class="mb-6 transform hover:shadow-lg transition-all duration-300">
      <template #header>
        <div class="flex items-center">
          <el-icon class="mr-2"><InfoFilled /></el-icon>
          <span>{{ t('localization.projectIntro') }}</span>
        </div>
      </template>
      <div class="project-info">
        <h2 class="text-xl font-bold mb-4">{{ t('localization.websites') }}</h2>
        <ul class="link-list space-y-3 list-disc pl-6">
          <li>
            {{ t('localization.officialDocs') }}
            <a
              href="https://www.scope-sentry.top"
              target="_blank"
              class="text-blue-500 hover:text-blue-700 cursor-pointer"
            >
              https://www.scope-sentry.top
            </a>
          </li>
          <li>
            Github：
            <a
              href="https://github.com/Autumn-27/ScopeSentry"
              target="_blank"
              class="text-blue-500 hover:text-blue-700 cursor-pointer"
            >
              https://github.com/Autumn-27/ScopeSentry
            </a>
          </li>
          <li>
            {{ t('localization.scannerSource') }}
            <a
              href="https://github.com/Autumn-27/ScopeSentry-Scan"
              target="_blank"
              class="text-blue-500 hover:text-blue-700 cursor-pointer"
            >
              https://github.com/Autumn-27/ScopeSentry-Scan
            </a>
          </li>
          <li>
            {{ t('localization.pluginMarket') }}
            <a
              href="https://plugin.scope-sentry.top/"
              target="_blank"
              class="text-blue-500 hover:text-blue-700 cursor-pointer"
            >
              https://plugin.scope-sentry.top/
            </a>
          </li>
          <li>
            SecureFlow：
            <a
              href="https://plugin.scope-sentry.top/"
              target="_blank"
              class="text-blue-500 hover:text-blue-700 cursor-pointer"
            >
              https://secureflow.scope-sentry.top
            </a>
          </li>
        </ul>
      </div>
    </el-card>

    <!-- 赞助合作卡片 -->
    <el-card class="transform hover:shadow-lg transition-all duration-300">
      <template #header>
        <div class="flex items-center">
          <el-icon class="mr-2"><Money /></el-icon>
          <span>{{ t('localization.sponsorCooperate') }}</span>
        </div>
      </template>
      <div class="sponsor-content p-4">
        <h3 class="text-2xl font-bold mb-4 text-center">{{ t('localization.supportProject') }}</h3>
        <p class="text-gray-600 mb-8 text-center max-w-2xl mx-auto">
          {{ t('localization.supportDescription') }}
        </p>
        <div class="sponsor-options grid grid-cols-1 md:grid-cols-2 gap-8">
          <el-card shadow="hover" class="sponsor-card">
            <template #header>
              <h4 class="font-bold text-lg text-center">{{ t('localization.sponsor') }}</h4>
            </template>
            <div class="flex flex-col items-center p-4">
              <!-- 赞助内容 -->
              <div class="qr-codes flex flex-col md:flex-row gap-8 justify-center items-center">
                <div
                  v-for="qr in qrCodes"
                  :key="qr.name"
                  class="qr-code-item flex flex-col items-center"
                >
                  <img
                    :src="qr.img"
                    :alt="t('localization.paymentQr', { name: qr.name })"
                    class="w-64 h-64 object-cover mb-3 hover:scale-105 transition-transform duration-300"
                  />
                  <span class="text-gray-600 text-lg">{{ qr.name }}</span>
                </div>
              </div>
            </div>
          </el-card>

          <el-card shadow="hover" class="sponsor-card">
            <template #header>
              <h4 class="font-bold text-lg text-center">{{ t('localization.cooperate') }}</h4>
            </template>
            <div class="flex flex-col items-center p-4">
              <p class="text-gray-600 mb-6 text-center break-all">{{ t('localization.cooperateDescription') }}</p>
              <el-button type="primary" size="large" class="w-full md:w-auto">
                rainy-autumn@outlook.com
              </el-button>
            </div>
          </el-card>
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.about-container {
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.sponsor-card {
  background: linear-gradient(145deg, #ffffff 0%, #f8f9fa 100%);
  border: none;
  width: 100%;
  overflow: hidden;
}

.el-card {
  border-radius: 12px;
  border: none;
  overflow: hidden;
}

:deep(.el-card__header) {
  border-bottom: 2px solid #f0f0f0;
  padding: 1rem 1.5rem;
}

:deep(.el-button--large) {
  padding: 12px 30px;
  font-weight: 600;
}

.link-list a {
  text-decoration: none;
  transition: all 0.3s ease;
  cursor: pointer;
}

.link-list a:hover {
  text-decoration: underline;
  cursor: pointer;
}

:deep(.el-card__body) {
  padding: 1rem;
}

.qr-code-item img {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.qr-codes {
  width: 100%;
  padding: 1rem;
}

@media (max-width: 768px) {
  .qr-code-item img {
    width: 200px; /* 移动端稍微小一点 */
    height: 200px;
  }
}
</style>
