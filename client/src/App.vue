<script setup lang="ts">
import Seats from '@/components/Seats.vue'
import Stage from '@/components/Stage.vue'
import Controls from '@/components/Controls/Controls.vue'
import Screen from '@/components/Screen.vue'
import { open_socket_client } from '@/socket_client'
import { computed, onMounted, reactive, ref } from 'vue'

const screen_ref = ref<Screen | null>(null)
const seats_ref = ref<Seats>()
const client_future = reactive(open_socket_client())

onMounted(async () => {
  try {
    const client = await client_future
    console.log('Connected to web socket server')
    seats_ref.value.on_client_connected(client)
  } catch (error) {
    console.error(`Failed to connect to web socket server: ${error}`)
  }
})
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
