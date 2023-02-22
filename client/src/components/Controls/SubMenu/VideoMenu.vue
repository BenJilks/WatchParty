<script setup lang="ts">
  import VideoItem from '@/components/Controls/Video/VideoItem.vue'
  import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
  import type { VideoData } from '@/components/Controls/Video/VideoItem.vue'
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
  const sub_menu = ref<SubMenu>()
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

  function selected(video_file: string) {
    emit('selected', video_file)
  }

  const toggle = () => sub_menu.value?.toggle()
  defineExpose({
    toggle,
    sub_menu,
  })

  request_video_list()
</script>

<template>
  <SubMenu ref="sub_menu" height="50vh">
    <div id="content-list">
      <VideoItem
          v-for="video in video_list"
          :video="video"
          :key="video.name"
          @selected="selected" />
    </div>
  </SubMenu>
</template>

<style scoped>
  #content-list {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    border-radius: 1em;
    margin: 2em;

    display: grid;
    grid-auto-flow: row;
    grid-template-columns: repeat(auto-fill, minmax(15em, 1fr));
    grid-auto-rows: min-content;
    gap: 0.5em;
  }
</style>
