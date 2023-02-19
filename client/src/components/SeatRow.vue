<script setup lang="ts">
  import Monkey from "@/components/Monkey.vue"
  import type { MonkeyData } from '@/monkey'
  import { reactive } from 'vue'
  import { create_monkey, seats_in_row } from '@/monkey'

  const monkeys = reactive<MonkeyData[]>([])
  const props = defineProps({
    row: { type: Number, required: true },
  })

  function add_monkey() {
    const max_seats = seats_in_row(props.row)
    const seat = Math.trunc(Math.random() * max_seats)
    monkeys.push(create_monkey(props.row, seat))
    console.log(`Create new monkey at row ${ props.row } in seat ${ seat }`)
  }

  defineExpose({
    add_monkey,
  })
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
