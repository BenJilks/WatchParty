<script setup lang="ts">
  import BufferedSegment from '@/components/Controls/BufferedSegment.vue'
  import type { Ref } from 'vue'
  import { computed, reactive, ref, watch } from 'vue'

  interface Props {
    video: Ref<HTMLVideoElement | null>,
  }

  interface BufferedSegmentData {
    start: number,
    end: number,
  }

  let is_seeking = false
  const seek_bar_ref = ref<HTMLDivElement | null>(null)
  const data = reactive({
    progress: 0,
    duration: 0,
    play_pause_icon: 'play',
  })

  const buffered_segments = reactive<BufferedSegmentData[]>([])
  function process_time_ranges(time_ranges: TimeRanges) {
    let start = 0, end = 0
    for (let i = 0; i < time_ranges.length; i++) {
      start = Math.min(start, time_ranges.start(i))
      end = Math.max(end, time_ranges.end(i))
    }

    buffered_segments.push({
      start: start / data.duration,
      end: end / data.duration,
    })
  }

  function update_buffered_segments(video: HTMLVideoElement) {
    buffered_segments.splice(0, buffered_segments.length)
    process_time_ranges(video.buffered)
  }

  const props = defineProps<Props>()
  const progress_factor = computed(() =>
      data.progress / data.duration)

  watch(props.video, () => {
    const video = props.video.value!
    video.preload = 'auto'
    video.addEventListener('durationchange', () =>
        data.duration = video.duration)
    video.addEventListener('progress', () =>
        update_buffered_segments(video))
    video.addEventListener('timeupdate', () => {
      if (!is_seeking)
        data.progress = video.currentTime
    })
  })

  function play(video: HTMLVideoElement) {
    video.play()
    data.play_pause_icon = 'pause'
  }

  function pause(video: HTMLVideoElement) {
    video.pause()
    data.play_pause_icon = 'play'
  }

  function play_pause() {
    const video = props.video.value
    if (video != null) {
      video.paused
          ? play(video)
          : pause(video)
    }
  }

  function on_seek_start() {
    if (is_seeking)
      return
    props.video.value?.pause()
    is_seeking = true
  }

  window.addEventListener('mouseup', () => {
    if (!is_seeking)
      return
    if (data.play_pause_icon == 'pause')
      props.video.value?.play()
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
    if (progress < 0) {
      data.progress = 0
    } else if (progress > 1) {
      data.progress = data.duration
    } else {
      data.progress = data.duration * progress
    }

    if ('fastSeek' in video) {
      video.fastSeek(data.progress)
    } else {
      // NOTE: Chrome does not have a `fastSeek` function
      (video as any).currentTime = data.progress
    }
  })
</script>

<template>
  <img
      :src="`/icons/${data.play_pause_icon}.svg`"
      draggable="false"
      id="play"
      @click="play_pause" />

  <div id="timeline" ref="seek_bar_ref" @mousedown="on_seek_start">
    <BufferedSegment
        v-for="(segment, i) in buffered_segments"
        :key="i"
        :start="segment.start"
        :end="segment.end" />
    <div id="done-timeline"></div>
    <div id="scrubber" @mousedown="on_seek_start"></div>
  </div>
</template>

<style scoped>
  * {
    --scrubber-size: 1.5em;
    --unplayed-color: #4445;
    --played-color: #333F;
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
    pointer-events: all;
    user-modify: none;
  }
</style>
