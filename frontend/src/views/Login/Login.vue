<script setup lang="ts">
import { computed } from 'vue'
import { LoginForm } from './components'
import { ThemeSwitch } from '@/components/ThemeSwitch'
import { LocaleDropdown } from '@/components/LocaleDropdown'
import { useI18n } from '@/hooks/web/useI18n'
import { underlineToHump } from '@/utils'
import { useAppStore } from '@/store/modules/app'
import { useDesign } from '@/hooks/web/useDesign'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('login')

const appStore = useAppStore()

const { t } = useI18n()

const title = computed(() => underlineToHump(appStore.getTitle))

const features = computed(() => [
  t('login.feature1'),
  t('login.feature2'),
  t('login.feature3')
])

const year = new Date().getFullYear()
</script>

<template>
  <div :class="prefixCls">
    <!-- 品牌侧：深色渐变 + 网格 + 光晕，深浅主题下都保持同一调性 -->
    <section :class="`${prefixCls}__brand`">
      <div :class="`${prefixCls}__grid`"></div>
      <div :class="`${prefixCls}__glow`"></div>

      <div :class="`${prefixCls}__brand-inner`">
        <div :class="`${prefixCls}__brand-logo`">
          <img src="@/assets/imgs/logo.png" alt="logo" />
          <span>{{ title }}</span>
        </div>

        <h1 :class="`${prefixCls}__headline`">{{ t('login.welcome') }}</h1>
        <p :class="`${prefixCls}__tagline`">{{ t('login.tagline') }}</p>

        <ul :class="`${prefixCls}__features`">
          <li v-for="item in features" :key="item">
            <Icon icon="ant-design:check-circle-filled" :size="16" />
            <span>{{ item }}</span>
          </li>
        </ul>
      </div>

      <div :class="`${prefixCls}__brand-foot`">© {{ year }} Scope Sentry</div>
    </section>

    <!-- 表单侧 -->
    <section :class="`${prefixCls}__panel`">
      <div :class="`${prefixCls}__toolbar`">
        <ThemeSwitch />
        <LocaleDropdown />
      </div>

      <div :class="`${prefixCls}__card`">
        <div :class="`${prefixCls}__card-head`">
          <h2>{{ t('login.login') }}</h2>
          <p>{{ t('login.subtitle') }}</p>
        </div>
        <LoginForm />
      </div>

      <div :class="`${prefixCls}__panel-foot`">{{ title }}</div>
    </section>
  </div>
</template>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-login';

.@{prefix-cls} {
  display: grid;
  grid-template-columns: minmax(440px, 46%) 1fr;
  min-height: 100vh;
  overflow: auto;
  background-color: var(--el-bg-color-page);

  /* ---------------- 品牌侧 ---------------- */
  &__brand {
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 48px 56px;
    overflow: hidden;
    color: #eef2ff;
    background:
      radial-gradient(1200px 600px at -10% -20%, rgb(47 111 235 / 55%), transparent 60%),
      linear-gradient(160deg, #0b1220 0%, #101a2e 45%, #0a0f18 100%);
  }

  &__grid {
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgb(255 255 255 / 6%) 1px, transparent 1px),
      linear-gradient(90deg, rgb(255 255 255 / 6%) 1px, transparent 1px);
    background-size: 44px 44px;
    mask-image: radial-gradient(circle at 30% 20%, #000 0%, transparent 72%);
  }

  &__glow {
    position: absolute;
    top: 8%;
    right: -160px;
    width: 420px;
    height: 420px;
    background: radial-gradient(circle, rgb(47 111 235 / 45%) 0%, transparent 65%);
    filter: blur(10px);
  }

  &__brand-inner {
    position: relative;
    z-index: 1;
    max-width: 520px;
    margin-top: 12vh;
  }

  &__brand-logo {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 56px;

    img {
      width: 40px;
      height: 40px;
      border-radius: 10px;
      box-shadow: 0 6px 18px rgb(0 0 0 / 35%);
    }

    span {
      font-size: 17px;
      font-weight: 600;
      letter-spacing: 0.02em;
    }
  }

  &__headline {
    margin: 0 0 14px;
    font-size: 34px;
    font-weight: 700;
    line-height: 1.25;
    letter-spacing: -0.02em;
    color: #f8fafc;
  }

  &__tagline {
    max-width: 460px;
    margin: 0 0 36px;
    font-size: 15px;
    line-height: 1.7;
    color: rgb(226 232 240 / 72%);
  }

  &__features {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 0;
    margin: 0;
    list-style: none;

    li {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 13.5px;
      color: rgb(226 232 240 / 82%);
    }

    :deep(.iconify) {
      color: #6d9af1;
    }
  }

  &__brand-foot {
    position: relative;
    z-index: 1;
    font-size: 12px;
    color: rgb(226 232 240 / 45%);
  }

  /* ---------------- 表单侧 ---------------- */
  &__panel {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px 72px;
    background:
      radial-gradient(760px 420px at 80% 0%, rgb(47 111 235 / 8%), transparent 60%),
      var(--el-bg-color-page);
  }

  &__toolbar {
    position: absolute;
    top: 20px;
    right: 24px;
    display: flex;
    gap: 8px;
    align-items: center;
  }

  &__card {
    width: 100%;
    max-width: 420px;
    padding: 32px 32px 28px;
    background: var(--app-surface);
    border: 1px solid var(--app-border);
    border-radius: 16px;
    box-shadow: var(--app-shadow-md);
  }

  &__card-head {
    margin-bottom: 24px;

    h2 {
      margin: 0 0 8px;
      font-size: 22px;
      font-weight: 650;
      letter-spacing: -0.02em;
    }

    p {
      margin: 0;
      font-size: 13px;
      color: var(--el-text-color-secondary);
    }
  }

  /* 登录表单细节：与卡片保持同一套留白和圆角，避免出现"框中框" */
  &__card :deep(.login-form) {
    .el-form-item {
      margin-bottom: 20px;
    }

    .el-form-item--label-top .el-form-item__label {
      margin-bottom: 8px;
    }

    .el-input__wrapper {
      min-height: 42px;
      padding: 1px 12px;
      border-radius: 10px;
    }

    .el-checkbox {
      height: 24px;
      margin-right: 0;
    }

    .el-checkbox__label {
      padding-left: 8px;
      font-size: 13px;
      color: var(--el-text-color-secondary);
    }

    .el-button {
      height: 44px;
      font-size: 14px;
      font-weight: 600;
      border-radius: 10px;
    }

    /* Chrome 自动填充在深色下会强行变白，这里压回主题色 */
    input:-webkit-autofill,
    input:-webkit-autofill:hover,
    input:-webkit-autofill:focus {
      -webkit-text-fill-color: var(--el-text-color-primary);
      caret-color: var(--el-text-color-primary);
      box-shadow: 0 0 0 1000px var(--app-surface) inset;
      transition: background-color 9999s ease-out 0s;
    }
  }

  &__panel-foot {
    position: absolute;
    right: 0;
    bottom: 24px;
    left: 0;
    font-size: 12px;
    color: var(--el-text-color-placeholder);
    text-align: center;
  }

  /* ---------------- 响应式 ---------------- */
  @media (max-width: 1080px) {
    grid-template-columns: 1fr;

    &__brand {
      display: none;
    }

    &__panel {
      padding-top: 96px;
    }
  }
}
</style>
