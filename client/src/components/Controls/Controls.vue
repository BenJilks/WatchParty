<script setup lang="ts">
  import VideoControls from '@/components/Controls/VideoControls.vue'
  import Chat from '@/components/Controls/Chat.vue'
  import SelectMenu from '@/components/Controls/SelectMenu.vue'
  import type { Ref } from 'vue'
  import { SocketClient } from '@/socket_client'
  import {onMounted, ref} from 'vue'

  interface Props {
    video: Ref<HTMLVideoElement | null>,
    client_future: Promise<SocketClient>,
  }

  type ChangeVideoMessage = {
    video_file: string,
  }

  const select_menu_ref = ref<SelectMenu | null>(null)
  const video_controls_ref = ref<VideoControls | null>(null)
  const select_menu_open_ref = ref(false)
  const props = defineProps<Props>()

  function toggle_select_menu() {
    const select_menu = select_menu_ref.value
    if (select_menu != null) {
      select_menu_open_ref.value = !select_menu_open_ref.value
      select_menu?.toggle()
    }
  }

  function change_video(video_file: string) {
    console.log(`Selected video: '${ video_file }'`)
    const video_controls = video_controls_ref.value
    if (video_controls != null)
      video_controls.change_video(video_file)
  }

  async function video_selected(video_file: string) {
    change_video(video_file)
    toggle_select_menu()

    const client = await props.client_future
    client.send<ChangeVideoMessage>('video-change', {
      video_file: video_file,
    })
  }

  onMounted(async () => {
    const client = await props.client_future
    client.on<ChangeVideoMessage>('video-change', message => {
      change_video(message.video_file)
    })
  })
</script>

<template>
  <div class="controls">
    <div class="panel">
      <Chat :client_future="client_future" />
    </div>

    <div id="video-panel" class="panel">
      <VideoControls
          ref="video_controls_ref"
          :video="video"
          :client_future="client_future" />
    </div>

    <div class="panel">
      <img
          class="icon"
          id="menu-button"
          src="/icons/up.svg"
          @click="toggle_select_menu"
          alt="" />
    </div>

    <SelectMenu
        ref="select_menu_ref"
        :client_future="client_future"
        @selected="video_selected" />
  </div>
</template>

<style scoped>
  * {
    --floating-x: 3em;
    --floating-y: 2em;

    --controls-height: 3em;
    --panel-color: #FFFb;
  }

  .controls {
    position: fixed;
    display: flex;
    align-items: center;
    gap: 2em;

    height: var(--controls-height);
    width: calc(100% - var(--floating-x)*2);
    bottom: var(--floating-y);
    left: var(--floating-x);
  }

  #video-panel {
    width: 100%;
  }

  #menu-button {
    transition: transform 0.2s;
    transform: rotate(v-bind('`${ select_menu_open_ref ? 180 : 0 }deg`'))
  }
</style>
