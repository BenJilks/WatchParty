<script setup lang="ts">
  import VideoControls from '@/components/Controls/VideoControls.vue'
  import Chat from '@/components/Controls/Chat.vue'
  import type { Ref } from 'vue'
  import { SocketClient } from '@/socket_client'

  interface Props {
    video: Ref<HTMLVideoElement | null>,
    client_future: Promise<SocketClient>,
  }

  defineProps<Props>()
</script>

<template>
  <div class="panel">
    <div class="sub-panel">
      <Chat :client_future="client_future" />
    </div>

    <div id="video-panel" class="sub-panel">
      <VideoControls :video="video" :client_future="client_future" />
    </div>
  </div>
</template>

<style scoped>
  * {
    --floating-x: 3em;
    --floating-y: 2em;
  }

  .panel {
    position: fixed;
    display: flex;
    align-items: center;
    gap: 2em;

    height: 3em;
    width: calc(100% - var(--floating-x)*2);
    bottom: var(--floating-y);
    left: var(--floating-x);
  }

  .sub-panel {
    display: flex;
    align-items: center;
    gap: 1em;
    height: 100%;

    padding: 0.5em 1em;
    border-radius: 0.5em;
    background-color: #FFF9;
  }

  #video-panel {
    width: 100%;
  }
</style>
