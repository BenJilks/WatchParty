<script setup lang="ts">
  import type { Ref } from 'vue'
  import {computed, reactive, ref, watch} from 'vue'

  interface Props {
    video: Ref<HTMLVideoElement | null>,
  }

  const seek_bar_ref = ref<HTMLDivElement | null>(null)
  const data = reactive({
    progress: 0,
    duration: 0,
    play_pause_icon: 'play',
  })

  const progress_factor = computed(() => {
      return data.progress / data.duration
  })

  const props = defineProps<Props>()
  watch(props.video, () => {
    const video = props.video.value
    video?.addEventListener('timeupdate', () => {
      data.duration = video.duration
      data.progress = video.currentTime
    })
  })

  function play_pause() {
    const video = props.video.value
    if (video == null)
      return

    if (video.paused) {
      video.play()
      data.play_pause_icon = 'pause'
    } else {
      video.pause()
      data.play_pause_icon = 'play'
    }
  }

  let is_seeking = false

  function on_seek_start() {
    is_seeking = true
  }

  window.addEventListener('mouseup', () => {
    is_seeking = false
  })

  window.addEventListener('mousemove', event => {
    if (!is_seeking)
      return

    const video = props.video.value
    if (video == null)
      return

    const seek_bar = seek_bar_ref.value
    const bounding_rect = seek_bar?.getBoundingClientRect()
    const x_along_seek_bar = event.x - (bounding_rect?.x ?? 0) - (1.5 / 2)*16
    const progress = x_along_seek_bar / (bounding_rect?.width ?? 1)
    data.progress = data.duration * progress
    video.fastSeek(data.progress)
  })
</script>

<template>
  <div class="panel">
    <img
        :src="`/icons/${data.play_pause_icon}.svg`"
        draggable="false"
        id="play"
        @click="play_pause" />

    <div id="timeline" ref="seek_bar_ref" @mousedown="on_seek_start">
      <div id="scrubber" @mousedown="on_seek_start"></div>
      <div id="done-timeline"></div>
    </div>
  </div>
</template>

<style scoped>
  * {
    --floating-x: 3em;
    --floating-y: 2em;

    --scrubber-size: 1.5em;
    --unplayed-color: #4445;
    --played-color: #444F;
  }

  .panel {
    position: fixed;
    display: flex;
    align-items: center;
    gap: 1.5em;

    height: 4em;
    width: calc(100% - var(--floating-x)*2);
    bottom: var(--floating-y);
    left: var(--floating-x);

    padding: 1em 2em;
    border-radius: 0.5em;
    background-color: #FFF9;
  }

  #play {
    height: 100%;
    aspect-ratio: 1/1;
    cursor: pointer;
    pointer-events: all;
    user-modify: none;
  }

  #timeline {
    position: relative;
    width: 100%;
    height: 0.5em;

    border-radius: 1em;
    background-color: var(--unplayed-color);
  }

  #done-timeline {
    position: absolute;
    width: calc(100% * v-bind('progress_factor'));
    height: 100%;

    border-radius: 1em;
    background-color: var(--played-color);
  }

  #scrubber {
    position: absolute;
    width: var(--scrubber-size);
    height: var(--scrubber-size);

    top: 0.25em;
    transform: translateY(-50%);
    left: calc(100% * v-bind('progress_factor'));

    border-radius: 50%;
    background-color: black;

    cursor: pointer;
    pointer-events: painted;
  }
</style>
