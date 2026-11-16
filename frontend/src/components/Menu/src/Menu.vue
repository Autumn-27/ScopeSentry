<script lang="tsx">
import { computed, defineComponent, unref, PropType } from 'vue'
import { ElMenu, ElScrollbar } from 'element-plus'
import { useAppStore } from '@/store/modules/app'
import { usePermissionStore } from '@/store/modules/permission'
import { useRenderMenuItem } from './components/useRenderMenuItem'
import { useRouter } from 'vue-router'
import { isUrl } from '@/utils/is'
import { useDesign } from '@/hooks/web/useDesign'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('menu')

export default defineComponent({
  name: 'Menu',
  props: {
    menuSelect: {
      type: Function as PropType<(index: string) => void>,
      default: undefined
    }
  },
  setup(props) {
    const appStore = useAppStore()

    const layout = computed(() => appStore.getLayout)
    const { push, currentRoute } = useRouter()

    const permissionStore = usePermissionStore()

    const menuMode = computed((): 'vertical' | 'horizontal' => {
      // 竖
      const vertical: LayoutType[] = ['classic', 'topLeft', 'cutMenu']

      if (vertical.includes(unref(layout))) {
        return 'vertical'
      } else {
        return 'horizontal'
      }
    })

    const routers = computed(() =>
      unref(layout) === 'cutMenu' ? permissionStore.getMenuTabRouters : permissionStore.getRouters
    )

    const collapse = computed(() => appStore.getCollapse)

    const uniqueOpened = computed(() => appStore.getUniqueOpened)

    const activeMenu = computed(() => {
      const { meta, path } = unref(currentRoute)
      // if set path, the sidebar will highlight the path you set
      if (meta.activeMenu) {
        return meta.activeMenu as string
      }
      return path
    })

    const menuSelect = (index: string) => {
      if (props.menuSelect) {
        props.menuSelect(index)
      }
      // 自定义事件
      if (isUrl(index)) {
        window.open(index)
      } else {
        push(index)
      }
    }

    const renderMenuWrap = () => {
      if (unref(layout) === 'top') {
        return renderMenu()
      } else {
        return <ElScrollbar>{renderMenu()}</ElScrollbar>
      }
    }

    const renderMenu = () => {
      return (
        <ElMenu
          defaultActive={unref(activeMenu)}
          mode={unref(menuMode)}
          collapse={
            unref(layout) === 'top' || unref(layout) === 'cutMenu' ? false : unref(collapse)
          }
          uniqueOpened={unref(layout) === 'top' ? false : unref(uniqueOpened)}
          textColor="var(--left-menu-text-color)"
          activeTextColor="var(--left-menu-text-active-color)"
          onSelect={menuSelect}
        >
          {{
            default: () => {
              const { renderMenuItem } = useRenderMenuItem(unref(menuMode))
              return renderMenuItem(unref(routers))
            }
          }}
        </ElMenu>
      )
    }

    return () => (
      <div
        id={prefixCls}
        class={[
          `${prefixCls} ${prefixCls}__${unref(menuMode)}`,
          'h-[100%] overflow-hidden flex-col bg-[var(--left-menu-bg-color)]',
          {
            'w-[var(--left-menu-min-width)]': unref(collapse) && unref(layout) !== 'cutMenu',
            'w-[var(--left-menu-max-width)]': !unref(collapse) && unref(layout) !== 'cutMenu'
          }
        ]}
      >
        {renderMenuWrap()}
      </div>
    )
  }
})
</script>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-menu';

.@{prefix-cls} {
  position: relative;
  padding: 8px;
  overflow: hidden auto;
  background-color: var(--app-surface);
  transition: width var(--transition-time-02);

  :deep(.@{elNamespace}-menu) {
    width: 100%;
    background: transparent;
    border-right: none;

    .@{elNamespace}-menu-item,
    .@{elNamespace}-sub-menu__title {
      height: 40px;
      margin: 2px 0;
      padding-right: 12px;
      line-height: 40px;
      color: var(--el-text-color-regular);
      border-radius: var(--app-radius-sm);
      transition:
        background-color 0.16s ease,
        color 0.16s ease;

      &:hover {
        color: var(--el-text-color-primary) !important;
        background-color: var(--el-menu-hover-bg-color) !important;
      }
    }

    .@{elNamespace}-menu-item.is-active {
      font-weight: 600;
      color: var(--app-primary) !important;
      background-color: var(--app-primary-soft) !important;

      &:hover {
        background-color: var(--app-primary-soft) !important;
      }
    }

    .is-active > .@{elNamespace}-sub-menu__title {
      font-weight: 600;
      color: var(--el-text-color-primary) !important;
    }

    .@{elNamespace}-sub-menu .@{elNamespace}-menu {
      background: transparent;
    }

    .@{elNamespace}-menu-item .@{elNamespace}-icon,
    .@{elNamespace}-sub-menu__title .@{elNamespace}-icon {
      margin-right: 10px;
      font-size: 17px;
    }
  }

  /* 折叠态：图标居中，去掉左右内边距 */
  :deep(.@{elNamespace}-menu--collapse) {
    width: var(--left-menu-min-width);

    .@{elNamespace}-menu-item,
    .@{elNamespace}-sub-menu__title {
      justify-content: center;
      padding: 0 !important;

      .@{elNamespace}-icon {
        margin-right: 0;
      }
    }
  }

  /* 水平菜单（top 布局） */
  &__horizontal {
    height: var(--top-tool-height) !important;
    padding: 0 8px;
    background: transparent;

    :deep(.@{elNamespace}-menu--horizontal) {
      height: var(--top-tool-height);
      background: transparent;
      border-bottom: none;

      .@{elNamespace}-menu-item,
      .@{elNamespace}-sub-menu__title {
        height: var(--top-tool-height);
        margin: 0 2px;
        line-height: var(--top-tool-height);
        border-bottom: 2px solid transparent;
        border-radius: 0;
        background: transparent !important;

        &:hover {
          background: transparent !important;
        }
      }

      .@{elNamespace}-menu-item.is-active {
        color: var(--app-primary) !important;
        background: transparent !important;
        border-bottom-color: var(--app-primary);

        &::after {
          display: none !important;
        }
      }
    }
  }
}
</style>

<style lang="less">
@prefix-cls: ~'@{namespace}-menu-popper';

/* 折叠后弹出的子菜单（挂在 body 上，不能用 scoped） */
.@{prefix-cls}--vertical,
.@{prefix-cls}--horizontal {
  .el-sub-menu__title,
  .el-menu-item {
    &:hover {
      color: var(--el-text-color-primary) !important;
      background-color: var(--el-menu-hover-bg-color) !important;
    }
  }

  .el-menu-item.is-active {
    color: var(--app-primary) !important;
    background-color: var(--app-primary-soft) !important;
  }
}
</style>
