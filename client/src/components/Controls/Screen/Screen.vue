<script setup lang="ts">
import {computed, ref} from 'vue'
import TextTool from '@/components/Controls/Screen/TextTool'
import LineTool from '@/components/Controls/Screen/LineTool'
import AnnotationTool from '@/components/Controls/Screen/AnnotationTool'
import type {CursorProperty} from "csstype";

const video_ref = ref<HTMLVideoElement | null>(null)
const synchronising = ref(true)
const needs_focus = ref(false)

const tools = ref<{[name: string]: AnnotationTool}>({
  text: new TextTool(),
  line: new LineTool(),
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

defineExpose({
  set_synchronising: (value: boolean) => synchronising.value = value,
  set_needs_focus: (value: boolean) => needs_focus.value = value,
  video_ref,
  tools,
})
</script>

<template>
  <div class="screen">
    <video ref="video_ref">
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
