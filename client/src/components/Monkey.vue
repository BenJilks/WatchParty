<script setup lang="ts">
  import type { MonkeyData } from '@/monkey'
  import type { Ref } from 'vue'
  import { SocketClient } from "@/socket_client"
  import { reactive } from 'vue'

  interface Props {
    monkey: MonkeyData,
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
    // TODO: Send clap message to server
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

  if (props.monkey.your_token !== undefined) {
    window.addEventListener('keydown', event =>
        isSpace(event) && clapDown())
    window.addEventListener('keyup', event =>
        isSpace(event) && clapUp())
  }
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
