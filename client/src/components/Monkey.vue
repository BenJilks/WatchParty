<script setup lang="ts">
import type {ClapMessage, ClapResponseMessage, MonkeyData} from '@/monkey'
  import { SocketClient } from "@/socket_client"
  import { reactive } from 'vue'

  interface Props {
    monkey: MonkeyData,
    row: number,
    client_future: Promise<SocketClient>,
  }

  enum Sprite {
    Idle = 'idle',
    Ready = 'clap_ready',
    Clap = 'clap',
  }

  const props = defineProps<Props>()
  const sprite = reactive({ name: Sprite.Idle })

  function isSpace(event: Event): boolean {
    let keyboard_event = event as KeyboardEvent
    let key = keyboard_event.key
    return key === ' '
  }

  async function changeSpriteState(new_sprite: Sprite) {
    sprite.name = new_sprite

    const client = await props.client_future
    client.send<ClapMessage>('clap', {
      sprite: new_sprite,
      token: props.monkey.your_token ?? '',
    })
  }

  function clapDown() {
    changeSpriteState(Sprite.Ready)
  }

  function clapUp() {
    changeSpriteState(Sprite.Clap)
    setTimeout(() => {
      changeSpriteState(Sprite.Idle)
    }, 100)
  }

  if (props.monkey.your_token != undefined) {
    window.addEventListener('keydown', event =>
        isSpace(event) && clapDown())
    window.addEventListener('keyup', event =>
        isSpace(event) && clapUp())
  }

  props.client_future.then(client => {
    client.on<ClapResponseMessage>('clap', message => {
      if (message.row != props.row || message.column != props.monkey.seat) {
        return
      }
      sprite.name = message.sprite as Sprite
    })
  })
</script>

<template>
  <img :src="`/monkeys/${sprite.name}.png`" class="monkey" alt="Moonkie">
</template>

<style scoped>
  .monkey {
    position: absolute;
    left: calc(50%);
    transform: translateX(-50%);

    bottom: v-bind('`${ props.monkey.bottom + 1 }vh`');
    height: v-bind('`${ props.monkey.height }vh`');
    margin-left: v-bind('`${ props.monkey.x_offset }vh`');

    filter: brightness(0.5);
  }
</style>
