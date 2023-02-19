<script setup lang="ts">
  import Monkey from "@/components/Monkey.vue"
  import { create_monkey, seats_in_row } from '@/monkey'
  import type { MonkeyData } from '@/monkey'

  const props = defineProps({
    row: { type: Number, required: true },
  })

  const monkeys: MonkeyData[] = [];
  const seats_in_this_row = seats_in_row(props.row)
  for (let i = 0; i < seats_in_this_row; i++) {
    if (Math.random() > 0.7) {
      const monkey = create_monkey(props.row, i)
      monkeys.push(monkey)
    }
  }
</script>

<template>
  <Monkey v-for="monkey in monkeys" :key="monkey.x_offset" :monkey="monkey" />
  <img :src="`/cinema/seats/${row + 1}.png`" class="seats-row" alt="seat" />
</template>

<style scoped>
  .seats-row {
    position: absolute;
    width: auto;
    height: 100%;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
  }
</style>
