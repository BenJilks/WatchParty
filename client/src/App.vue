<script setup lang="ts">
import Seats from '@/components/Seats.vue'
import Stage from '@/components/Stage.vue'
import { open_socket_client, SocketClient } from '@/socket_client'
import {onMounted, ref} from "vue"

const client_ref = ref<SocketClient>()
const seats_ref = ref<Seats>()

function onSocketConnected(client: SocketClient) {
  client_ref.value = client

  const seats = seats_ref.value
  if (seats == undefined) {
    onMounted(() => onSocketConnected(client))
  } else {
    seats.on_client_connected(client)
  }
}

if (client_ref.value === undefined) {
  open_socket_client().then(new_client => {
    console.log('Connected to web socket server')
    onSocketConnected(new_client)
  })
  // .catch(error => {
  //   console.log(`Failed to connect to web socket server: ${error}`)
  // })
}
</script>

<template>
  <div class="background">
    <!--
    <video class="screen">
      <source src="/test.mp4" type="video/mp4">
    </video>
    !-->
    <img src="/test.png" class="screen" alt="" />
    <Stage />
    <Seats ref="seats_ref" />
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

    background-color: #bbbf;
    box-shadow: 0 0 5vh 1vh #fff6;
    border-radius: 0.5vh;
  }
</style>
