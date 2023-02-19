<script setup lang="ts">
  import type { MonkeyData } from '@/monkey';
  import { reactive } from 'vue';

  interface Props {
    monkey: MonkeyData,
  }

  const props = defineProps<Props>()
  const sprite = reactive({ name: "idle" })

  function startClap() {
    sprite.name = 'clap_ready'
    setTimeout(clap, 100)
  }

  function clap() {
    sprite.name = 'clap'
    setTimeout(endClap, 100)
  }

  function endClap() {
    sprite.name = 'idle'
    setTimeout(startClap, Math.random() * 1000)
  }

  setTimeout(startClap, Math.random() * 1000)
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
