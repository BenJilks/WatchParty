<script setup lang="ts">
  import VideoItem from '@/components/Controls/VideoItem.vue'
  import type { VideoData } from '@/components/Controls/VideoItem.vue'
  import { reactive, ref } from 'vue'
  import { SocketClient } from '@/socket_client'

  interface Props {
    client_future: Promise<SocketClient>,
  }

  interface Emits {
    (e: 'selected', video: string): void,
  }

  type VideoListMessage = {
    videos: VideoData[],
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const show = ref(false)
  const disabled = ref(true)
  const video_list = reactive<VideoData[]>([])

  async function request_video_list() {
    if (request_video_list.length > 0)
      return

    const client = await props.client_future
    client.on<VideoListMessage>('video-list', message => {
      video_list.splice(0, video_list.length)
      video_list.push(...message.videos)
    })
    client.send('video-list', null)
  }

  function toggle() {
    if (!show.value) {
      request_video_list();
      disabled.value = false
      setTimeout(() => show.value = true, 0)
    } else {
      show.value = false
      setTimeout(() => disabled.value = true, 200)
    }
  }

  function selected(video_file: string) {
    emit('selected', video_file)
  }

  defineExpose({
    toggle,
  })
</script>

<template>
  <div v-if="!disabled" class="panel" id="menu">
    <div id="content-list">
      <VideoItem
          v-for="video in video_list"
          :video="video"
          :key="video.name"
          @selected="selected" />
    </div>
  </div>
</template>

<style scoped>
  #menu {
    position: absolute;
    bottom: calc(var(--controls-height) + var(--floating-y));
    align-items: start;

    width: 100%;
    height: 50vh;
    padding: 2em;

    transition:
        opacity 0.2s,
        transform 0.2s;
    opacity: v-bind('show ? 1 : 0');
    transform: translateY(v-bind('`${ show ? 0 : 4 }em`'));
  }

  #content-list {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    border-radius: 1em;

    display: grid;
    grid-auto-flow: row;
    grid-template-columns: repeat(auto-fill, minmax(15em, 1fr));
    grid-auto-rows: min-content;
    gap: 0.5em;
  }
</style>
