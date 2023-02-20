<script setup lang="ts">
  interface Props {
    start: number,
    end: number,
  }

  const props = defineProps<Props>()
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
    width: calc(100% * v-bind('props.end - props.start'));
    height: 100%;

    border-radius: 1em;
    background-color: var(--buffered-color);
  }
</style>
