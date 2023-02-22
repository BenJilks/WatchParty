<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import TextTool from '@/components/Controls/Screen/TextTool'
import LineTool from '@/components/Controls/Screen/LineTool'
import AnnotationTool from '@/components/Controls/Screen/AnnotationTool'
import type {CursorProperty} from "csstype";

const screen_ref = ref<HTMLDivElement>()
const video_ref = ref<HTMLVideoElement>()
const synchronising = ref(true)
const needs_focus = ref(false)

const tools = ref<{[name: string]: AnnotationTool}>({})
onMounted(() => {
  tools.value = {
    text: new TextTool(screen_ref.value!),
    line: new LineTool(screen_ref.value!),
  }
})

function current_tool(): AnnotationTool | undefined {
  for (const name in tools.value) {
    const tool = tools.value[name]
    if (tool.enabled)
      return tool
  }

  return undefined
}

const current_cursor = computed(() =>
    current_tool()?.cursor ?? 'default')

function screen_click(event: MouseEvent) {
  const tool = current_tool()
  if (tool === undefined)
    return

  const screen_rect = video_ref.value?.getBoundingClientRect()
  const screen_x = screen_rect?.x ?? 0
  const screen_y = screen_rect?.y ?? 0
  tool.on_click(event.clientX - screen_x, event.clientY - screen_y)
}

defineExpose({
  set_synchronising: (value: boolean) => synchronising.value = value,
  set_needs_focus: (value: boolean) => needs_focus.value = value,
  video_ref,
  tools,
})
</script>

<template>
  <div class="screen" ref="screen_ref" @click="screen_click">
    <video ref="video_ref" preload="auto" oncontextmenu="return false;">
      <source :src="'/vids/[Rhythm Heaven] - Fan Club (Perfect) (English)-DNbvktlB0gU.mp4'" type="video/mp4">
    </video>
    <div v-if="synchronising && !needs_focus" class="overlay">
      Synchronising Viewers...
    </div>
    <div v-if="needs_focus" class="overlay">
      Click Anywhere to Play
    </div>
  </div>
</template>

<style scoped>
  .screen {
    display: flex;
    justify-content: center;
    align-items: center;

    position: absolute;
    bottom: calc(var(--seat-height) + 7vh);
    left: 50%;

    width: calc(100% - 25vw);
    height: calc(100% - var(--seat-height) - 9vh);
    transform: translateX(-50%);

    background-color: #0a0a0aff;
    box-shadow: 0 0 5vh 2vh #6666;
    border-radius: 0.5vh;

    cursor: v-bind('current_cursor');
    overflow: hidden;
  }

  .screen * {
    text-wrap: none;
    word-break: keep-all;
  }

  .screen video {
    position: absolute;
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .screen .overlay {
    padding: 0.5em 1em;
    border-radius: 0.5em;
    width: auto;
    height: auto;

    font-family: 'Trebuchet MS', sans-serif;
    font-size: 3em;
    font-weight: bold;

    color: white;
    background-color: #000a;
    z-index: 1;

    pointer-events: none;
    user-select: none;
  }
</style>
