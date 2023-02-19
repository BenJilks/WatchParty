<script setup lang="ts">
import type { Ref } from 'vue'
import SeatRow from '@/components/SeatRow.vue'
import { onMounted, ref } from 'vue'
import { SocketClient } from '@/socket_client'
import { seats_in_row } from '@/monkey'

const ROW_COUNT = 7
const rowsRef = ref<SeatRow[]>([])

type Seat = {
  row: number,
  column: number
}

type MessageStateUpdate = {
  seats_not_free: Seat[],
  your_token: string,
  your_seat: Seat,
}

function update_seat(update: MessageStateUpdate, row: SeatRow, seat: number) {
  const should_have_monkey = update.seats_not_free.findIndex(x =>
      x.row == row.row() && x.column == seat) != -1
  const has_monkey = !row.is_seat_empty(seat)

  const is_your_monkey = update.your_seat.row == row.row()
      && update.your_seat.column == seat
  const your_token = is_your_monkey ? update.your_token : undefined

  if (should_have_monkey && !has_monkey) {
    row.add_monkey(seat, your_token)
  } else if (!should_have_monkey && has_monkey) {
    row.remove_monkey(seat)
  }
}

function scan_seats_for_changes(update: MessageStateUpdate) {
  for (let rowIndex = 0; rowIndex < ROW_COUNT; rowIndex++) {
    const row: SeatRow = rowsRef.value[rowIndex]
    const seats_for_row = seats_in_row(rowIndex)

    for (let seat = 0; seat < seats_for_row; seat++) {
      update_seat(update, row, seat)
    }
  }
}

function on_client_connected(client: SocketClient) {
  if (rowsRef.value.length == 0) {
    console.log('Waiting for rows to be loaded')
    onMounted(() => on_client_connected(client))
    return
  }

  console.log('Listening for state changes')
  client.on('update-state', (update: MessageStateUpdate) => {
    scan_seats_for_changes(update)
  })
}

interface Props {
  client_future: Promise<SocketClient>,
}
const props = defineProps<Props>()

defineExpose({
  on_client_connected,
})
</script>

<template>
  <div class="seats-background"></div>
  <div class="front-board"></div>
  <div class="seats">
    <SeatRow v-for = "id in ROW_COUNT"
             :row = "id - 1"
             :key = "id - 1" ref="rowsRef"
             :client_future = "client_future" />
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
