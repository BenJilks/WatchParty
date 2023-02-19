<script setup lang="ts">
import SeatRow from '@/components/SeatRow.vue'
import { onMounted, ref } from "vue";

const ROW_COUNT = 6

const rowsRef = ref<SeatRow[]>([])
onMounted(() => {
  const rowIndex = Math.trunc(Math.random() * ROW_COUNT)
  const row = rowsRef.value[rowIndex]
  row.add_monkey()
})

</script>

<template>
  <div class="seats-background"></div>
  <div class="front-board"></div>
  <div class="seats">
    <SeatRow v-for="id in ROW_COUNT" :row="id" :key="id" ref="rowsRef" />
  </div>
</template>

<style scoped>
  .front-board {
    position: absolute;
    user-select: none;
    pointer-events: none;

    width: 100%;
    height: var(--front-board-height);

    bottom: calc(var(--seat-height) - var(--front-board-height));
    background: linear-gradient(0deg,
    #0B0B0BFF 0%,
    #181616FF 100%);
  }

  .seats-background {
    position: absolute;
    user-select: none;
    pointer-events: none;

    width: 100%;
    height: var(--seat-height);
    bottom: 0;

    background-color: black;
  }

  .seats {
    position: absolute;
    width: 100%;
    height: var(--seat-height);
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
  }
</style>
