<script setup lang="ts">
import Drawing from "@/components/Controls/Drawing/Drawing.vue";
import { computed, ref } from 'vue'
import VideoPlayer from "@/components/Controls/Video/VideoPlayer.vue";

const screen_ref = ref<HTMLDivElement>()
const video_player_ref = ref<VideoPlayer>()
const drawing_ref = ref<Drawing>()

const synchronising = ref(true)
const needs_focus = ref(false)

const tools = computed(() => drawing_ref?.value.tools)
const video_ref = computed(() => video_player_ref?.value.video_ref)
defineExpose({
  set_synchronising: (value: boolean) => synchronising.value = value,
  set_needs_focus: (value: boolean) => needs_focus.value = value,
  video_ref,
  tools,
})
</script>

<template>
  <div class="screen" ref="screen_ref">
    <VideoPlayer ref="video_player_ref" />
    <Drawing ref="drawing_ref" />
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
    z-index: 3;

    pointer-events: none;
    user-select: none;
  }
</style>
