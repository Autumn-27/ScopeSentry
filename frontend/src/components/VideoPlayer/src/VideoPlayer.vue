<script setup lang="ts">
// xgplayer 只在真正需要播放时加载（约 75KB gzip 前 + 播放器样式）
import type PlayerType from 'xgplayer'
import { ref, unref, onMounted, watch, onBeforeUnmount, nextTick } from 'vue'

const props = defineProps({
  url: {
    type: String,
    default: '',
    required: true
  },
  poster: {
    type: String,
    default: ''
  }
})

const playerRef = ref<InstanceType<typeof PlayerType>>()

const videoEl = ref<HTMLDivElement>()

let PlayerCtor: typeof PlayerType | null = null

const initPlayer = async () => {
  if (!unref(videoEl)) return
  if (!PlayerCtor) {
    const [mod] = await Promise.all([import('xgplayer'), import('xgplayer/dist/index.min.css')])
    PlayerCtor = mod.default
  }
  playerRef.value = new PlayerCtor({
    autoplay: false,
    ...props,
    el: unref(videoEl)
  }) as InstanceType<typeof PlayerType>
}

onMounted(() => {
  void initPlayer()
})

watch(
  () => props,
  async (newProps) => {
    await nextTick()
    if (newProps) {
      unref(playerRef)?.setConfig(newProps)
    }
  },
  {
    deep: true
  }
)

onBeforeUnmount(() => {
  unref(playerRef)?.destroy()
})

defineExpose({
  playerExpose: () => unref(playerRef)
})
</script>

<template>
  <div ref="videoEl"></div>
</template>
