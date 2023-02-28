<script setup lang="ts">
import Monkey from "@/components/ViewerMonkey.vue"
import type { MonkeyData } from '@/monkey'
import type { SocketClient } from '@/socket_client'
import { reactive } from 'vue'
import { create_monkey } from '@/monkey'

interface Props {
    row: number,
    client_future: Promise<SocketClient>,
}

const monkeys = reactive<MonkeyData[]>([])
const props = defineProps<Props>()

function add_monkey(seat: number, your_token?: string) {
    monkeys.push(create_monkey(props.row, seat, your_token))
    console.log(`Create new monkey at row ${ props.row } in seat ${ seat }`)
}

function remove_monkey(seat: number) {
    const index = monkeys.findIndex(data => data.seat == seat)
    if (index < 0)
        return

    monkeys.splice(index)
    console.log(`Removed monkey at row ${ props.row } in seat ${ seat }`)
}

function is_seat_empty(seat: number): boolean {
    return monkeys.findIndex(data => data.seat == seat) == -1
}

defineExpose({
    add_monkey,
    remove_monkey,
    is_seat_empty,
    row: props.row,
})
</script>

<template>
    <Monkey
        v-for="monkey in monkeys"
        :key="monkey.x_offset"
        :monkey="monkey"
        :row="row"
        :client_future="client_future" />

    <img
        :src="`/cinema/seats/${row + 1}.png`"
        class="seats-row"
        alt="seat"
        draggable="false" />
</template>

<style scoped>
.seats-row {
    position: absolute;
    width: auto;
    height: 100%;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    z-index: 0;
}
</style>
