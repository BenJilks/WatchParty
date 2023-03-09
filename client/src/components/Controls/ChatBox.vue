<script setup lang="ts">
import type { SocketClient } from '@/socket_client'
import { ref } from 'vue'

interface Props {
    client_future: Promise<SocketClient>,
}

interface MessageChat {
    message: string,
}

const input_ref = ref<HTMLInputElement>()
const props = defineProps<Props>()

async function send() {
    const input = input_ref.value
    if (input === undefined)
        return

    const message = input.value.normalize().trim()
    if (message == '')
        return

    const client = await props.client_future
    client.send<MessageChat>('chat', {
        message: message,
    })

    input.value = ''
}

function on_input_key(event: KeyboardEvent) {
    if (event.key == 'Enter') {
        send()
    }
}
</script>

<template>
    <input
        ref="input_ref"
        placeholder="Chat messages..."
        @keydown="on_input_key"/>
    <img
        src="/icons/chat.svg"
        draggable="false"
        id="send"
        @click="send"
        alt="send" />
</template>

<style scoped>
#send {
    height: 100%;
    aspect-ratio: 1/1;
    padding: 0.2em;

    cursor: pointer;
    pointer-events: all;
    user-modify: none;
}
</style>
