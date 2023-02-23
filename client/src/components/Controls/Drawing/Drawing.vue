<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import TextTool from '@/components/Controls/Drawing/TextTool'
import LineTool from '@/components/Controls/Drawing/LineTool'
import AnnotationTool from '@/components/Controls/Drawing/AnnotationTool'
import type { Position } from '@/components/Controls/Drawing/AnnotationTool'

const canvas_ref = ref<HTMLCanvasElement>()
const tools = ref<{[name: string]: AnnotationTool}>({})

onMounted(() => {
  const canvas = canvas_ref.value
  if (!canvas)
    return

  const context = canvas.getContext('2d')!
  canvas.width = 800
  canvas.height = 600

  tools.value = {
    text: new TextTool(context),
    line: new LineTool(context),
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

function canvas_position_from_event(event: MouseEvent): Position {
  const canvas = canvas_ref.value!
  const rect = canvas.getBoundingClientRect()
  return {
    x: (event.clientX - rect.x) * (canvas.width / rect.width),
    y: (event.clientY - rect.y) * (canvas.height / rect.height),
  }
}

function screen_mouse_down(event: MouseEvent) {
  const position = canvas_position_from_event(event)
  current_tool()?.on_mouse_down(position)
}

function screen_mouse_up(event: MouseEvent) {
  const position = canvas_position_from_event(event)
  current_tool()?.on_mouse_up(position)
}

function screen_mouse_move(event: MouseEvent) {
  const position = canvas_position_from_event(event)
  current_tool()?.on_mouse_move(position)
}

function screen_mouse_event(event: Event, callback: (tool: AnnotationTool, position: Position) => undefined): () => undefined {
  return () => {
    const tool = current_tool()
    const position = canvas_position_from_event(event as MouseEvent)
    if (tool !== undefined) {
      callback(tool, position)
    }

    return undefined
  }
}

defineExpose({
  tools,
})
</script>

<template>
  <canvas ref="canvas_ref"
          @mousedown="screen_mouse_event($event, (tool, pos) => { tool.on_mouse_down(pos) })"
          @mouseup="screen_mouse_event($event, (tool, pos) => { tool.on_mouse_up(pos) })"
          @mousemove="screen_mouse_event($event, (tool, pos) => { tool.on_mouse_move(pos) })" />
</template>

<style scoped>
  canvas {
    position: absolute;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    z-index: 2;
  }
</style>