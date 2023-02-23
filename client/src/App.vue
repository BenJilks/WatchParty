<script setup lang="ts">
import Seats from '@/components/Seats.vue'
import Stage from '@/components/Stage.vue'
import Controls from '@/components/Controls/Controls.vue'
import Screen from '@/components/Controls/Screen.vue'
import { computed, inject, onMounted, ref } from 'vue'
import {SocketClient} from "@/socket_client";

const screen_ref = ref<Screen | null>(null)
const seats_ref = ref<Seats>()
const client_future = inject<Promise<SocketClient>>('client_future')

onMounted(async () => {
  try {
    const client = await client_future
    console.log('Connected to web socket server')
    seats_ref.value.on_client_connected(client)
  } catch (error) {
    console.error(`Failed to connect to web socket server: ${error}`)
  }
})

function zoom() {
  let zoom = 1

  window.onwheel = event => {
    if (!event.shiftKey && !event.ctrlKey)
      return

    zoom -= event.deltaY * (1.0 / window.innerHeight) * 0.4
    const scales: { [name: string]: number } = {
      '--seat-height': 20 * zoom,
      '--floor-height': 23 * zoom,
      '--curtain-offset': -(2.0 / 9.0 * 50 * (1.0 / zoom)) + 3,
      '--front-board-height': 6 * zoom,
    }

    for (const prop in scales) {
      document.body.style.setProperty(prop, `${scales[prop]}vh`)
    }
  }
}

zoom()
</script>

<template>
  <div class="background">
    <Screen ref="screen_ref" />
    <Stage />
    <Seats
        ref="seats_ref"
        :client_future="client_future" />
    <Controls
        :screen_ref="computed(() => screen_ref)"
        :client_future="client_future" />
  </div>
</template>

<style scoped>
  .background {
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;

    background-color: #050505;
  }
</style>
