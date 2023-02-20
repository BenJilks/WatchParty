<script setup lang="ts">
import Seats from '@/components/Seats.vue'
import Stage from '@/components/Stage.vue'
import { open_socket_client } from '@/socket_client'
import {computed, onMounted, reactive, ref} from "vue"
import Controls from "@/components/Controls/Controls.vue";

const seats_ref = ref<Seats>()
const video_ref = ref<HTMLVideoElement | null>(null)

const client_future = reactive(open_socket_client())
onMounted(async () => {
  try {
    const client = await client_future
    console.log('Connected to web socket server')
    seats_ref.value.on_client_connected(client)
  } catch (error) {
    console.log(`Failed to connect to web socket server: ${error}`)
  }
})
</script>

<template>
  <div class="background">
    <video ref="video_ref" class="screen" muted>
      <source src="/test.mp4" type="video/mp4">
    </video>

    <Stage />
    <Seats ref="seats_ref" :client_future="client_future" />
    <Controls :video="computed(() => video_ref)" />
  </div>
</template>

<style scoped>
  .background {
    position: fixed;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;

    background-color: #050505;
  }

  .screen {
    position: absolute;
    bottom: calc(var(--seat-height) + 7vh);

    object-fit: contain;
    width: calc(100% - 25vw);
    height: calc(100% - var(--seat-height) - 9vh);
    left: 50%;
    transform: translateX(-50%);

    background-color: #0a0a0aff;
    box-shadow: 0 0 5vh 1vh #fff6;
    border-radius: 0.5vh;
  }
</style>
