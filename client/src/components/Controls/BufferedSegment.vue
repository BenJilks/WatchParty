<script setup lang="ts">
  interface Props {
    start: number,
    end: number,
  }

  const props = defineProps<Props>()
  const length = props.end - props.start
</script>

<template>
  <div class="segment"></div>
</template>

<style scoped>
  * {
    --buffered-color: #0005;
  }

  .segment {
    position: absolute;
    left: calc(100% * v-bind('props.start'));
    width: calc(100% * v-bind('length'));
    height: 100%;

    border-radius: 1em;
    background-color: var(--buffered-color);
  }
</style>
