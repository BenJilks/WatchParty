<script setup lang="ts">
  import BufferedSegment from '@/components/Controls/BufferedSegment.vue'
  import type { Ref } from 'vue'
  import { computed, reactive, ref, watch } from 'vue'
  import { SocketClient } from '@/socket_client'

  interface Props {
    video: Ref<HTMLVideoElement | null>,
    client_future: Promise<SocketClient>,
  }

  interface BufferedSegmentData {
    start: number,
    end: number,
  }

  type VideoStateMessage = {
    playing: boolean,
    progress: number,
  }

  let is_seeking = false
  const seek_bar_ref = ref<HTMLDivElement | null>(null)
  const data = reactive({
    progress: 0,
    duration: 1,
    play_pause_icon: 'play',
    volume_icon: 'volume',
  })

  const buffered_segments = reactive<BufferedSegmentData[]>([])
  function process_time_ranges(time_ranges: TimeRanges) {
    let start = 0, end = 0
    for (let i = 0; i < time_ranges.length; i++) {
      start = Math.min(start, time_ranges.start(i))
      end = Math.max(end, time_ranges.end(i))
    }

    buffered_segments.push({
      start: start,
      end: end,
    })
  }

  function update_buffered_segments(video: HTMLVideoElement) {
    buffered_segments.splice(0, buffered_segments.length)
    process_time_ranges(video.buffered)
  }

  const props = defineProps<Props>()
  const progress_factor = computed(() =>
      data.progress / data.duration)

  function update_video_event_listeners(video: HTMLVideoElement) {
    video.preload = 'auto'
    video.addEventListener('durationchange', () =>
        data.duration = video.duration)
    video.addEventListener('progress', () =>
        update_buffered_segments(video))
    video.addEventListener('timeupdate', () => {
      if (!is_seeking)
        data.progress = video.currentTime
    })
  }

  watch(props.video, async () => {
    const video = props.video.value!
    update_video_event_listeners(video)

    const client = await props.client_future
    client.on<VideoStateMessage>('video', message => {
      const video = props.video.value!
      data.progress = message.progress
      video.currentTime = data.progress

      message.playing
          ? video.play()
          : video.pause()

      data.play_pause_icon = message.playing
          ? 'pause'
          : 'play'
    })
  })

  async function play(video: HTMLVideoElement) {
    video.play().catch(() => {
      video.oncanplay = () =>
          video.play()
    })
    data.play_pause_icon = 'pause'
  }

  function pause(video: HTMLVideoElement) {
    video.pause()
    data.play_pause_icon = 'play'
  }

  async function play_pause() {
    const video = props.video.value
    if (video == null)
      return

    video.paused
        ? await play(video)
        : pause(video)
    data.progress = video.currentTime

    const client = await props.client_future
    client.send<VideoStateMessage>('video', {
      playing: !video.paused,
      progress: video.currentTime,
    })
  }

  function toggle_mute() {
    const video = props.video.value
    if (video == null)
      return

    video.muted = !video.muted
    data.volume_icon = video.muted
        ? 'mute' : 'volume'
  }

  function seek(progress: number) {
    const video = props.video.value
    if (video == null)
      return

    if ('fastSeek' in video) {
      video.fastSeek(progress)
    } else {
      // NOTE: Chrome does not have a `fastSeek` function
      (video as any).currentTime = progress
    }
  }

  function on_seek_start() {
    if (is_seeking)
      return
    props.video.value?.pause()
    is_seeking = true
  }

  async function send_video_update() {
    const video = props.video.value
    if (video == null)
      return

    const client = await props.client_future
    client.send<VideoStateMessage>('video', {
      playing: !video.paused,
      progress: data.progress,
    })
  }

  window.addEventListener('mouseup', async () => {
    if (!is_seeking)
      return
    is_seeking = false

    const playing = (data.play_pause_icon == 'pause')
    if (playing) {
      await props.video.value?.play()
    }

    await send_video_update()
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

    seek(data.progress)
  })

  function change_video(video_file: string) {
    const video = props.video.value
    if (video == null)
      return

    video.src = `/vids/${ video_file }`
    data.progress = 0
    data.duration = video.duration
    data.play_pause_icon = video.paused ? 'play' : 'pause'
    data.volume_icon = video.muted ? 'mute' : 'volume'
    buffered_segments.splice(0, buffered_segments.length)
  }

  defineExpose({
    change_video,
  })
</script>

<template>
  <img
      class="icon"
      draggable="false"
      :src="`/icons/${data.play_pause_icon}.svg`"
      @click="play_pause" />

  <div id="timeline" ref="seek_bar_ref" @mousedown="on_seek_start">
    <BufferedSegment
        v-for="(segment, i) in buffered_segments"
        :key="i"
        :start="segment.start / data.duration"
        :end="segment.end / data.duration" />
    <div id="done-timeline"></div>
    <div id="scrubber" @mousedown="on_seek_start"></div>
  </div>

  <img
      class="icon"
      draggable="false"
      :src="`/icons/${data.volume_icon}.svg`"
      @click="toggle_mute" />
</template>

<style scoped>
  * {
    --scrubber-size: 1.5em;
    --unplayed-color: #4445;
    --played-color: #333F;
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
