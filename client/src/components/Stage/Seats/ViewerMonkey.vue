<script setup lang="ts">
import type { ClapMessage, ClapResponseMessage, ClapState, MonkeyData } from '@/monkey'
import type { SocketClient } from "@/socket_client"
import { reactive, ref } from 'vue'

const CHAT_MESSAGE_TIME_MS = 2500
const CLAP_TIME = 300

const IDLE_OFFSET = 1
const CLAP_READY_OFFSET = 3
const CLAP_OFFSET = 3

interface Props {
    monkey: MonkeyData,
    row: number,
    client_future: Promise<SocketClient>,
}

interface ChatResponseMessage {
    message: string,
    row: number,
    column: number,
}

enum Sprite {
    Idle = 'idle',
    Ready = 'clap_ready',
    Clap = 'clap',
}

interface MonkeyState {
    sprite: Sprite,
    offset: number,
}

interface ChatBubble {
    message: string,
    show: boolean,
    animation_speed: number,
}

const props = defineProps<Props>()
const monkey_image_ref = ref<HTMLInputElement | null>(null)
const monkey_state = reactive<MonkeyState>({
    sprite: Sprite.Idle,
    offset: IDLE_OFFSET,
})

const chat_message_enabled = ref(false)
const chat_message = reactive<ChatBubble>({
    message: '',
    show: false,
    animation_speed: 0.2,
})

function is_space(event: KeyboardEvent): boolean {
    return event.code == 'Space'
}

async function update_clap_state(state: ClapState) {
    const client = await props.client_future
    client.send<ClapMessage>('clap', {
        state: state,
        token: props.monkey.your_token ?? '',
    })
}

let idle_timeout: number | undefined = undefined

function clap_down() {
    monkey_state.sprite = Sprite.Ready
    monkey_state.offset = CLAP_READY_OFFSET
    update_clap_state('ready')
    clearTimeout(idle_timeout)
}

function clap_up() {
    monkey_state.sprite = Sprite.Clap
    monkey_state.offset = CLAP_OFFSET
    update_clap_state('clap')

    const monkey_image = monkey_image_ref.value!
    monkey_image.onload = () => {
        idle_timeout = setTimeout(() => {
            monkey_state.sprite = Sprite.Idle
            monkey_state.offset = IDLE_OFFSET
        }, CLAP_TIME)
        monkey_image.onload = null
    }
}

if (props.monkey.your_token != undefined) {
    window.addEventListener('keydown', event =>
            is_space(event) && clap_down())
    window.addEventListener('keyup', event =>
            is_space(event) && clap_up())
}

function show_chat_message(message: string) {
    chat_message_enabled.value = true
    setTimeout(() => {
        chat_message.message = message
        chat_message.show = true

        setTimeout(() => chat_message.show = false,
                CHAT_MESSAGE_TIME_MS)
        setTimeout(() => chat_message_enabled.value = false,
                CHAT_MESSAGE_TIME_MS + chat_message.animation_speed*1000)
    }, 10)
}

function is_me(row: number, column: number): boolean {
    return row == props.row && column == props.monkey.seat
}

props.client_future.then(client => {
    client.on<ClapResponseMessage>('clap', message => {
        if (!is_me(message.row, message.column))
            return

        switch (message.state) {
            case 'ready': {
                monkey_state.sprite = Sprite.Ready
                monkey_state.offset = CLAP_READY_OFFSET
                clearTimeout(idle_timeout)
                break
            }

            case 'clap': {
                monkey_state.sprite = Sprite.Clap
                monkey_state.offset = CLAP_OFFSET

                const monkey_image = monkey_image_ref.value!
                monkey_image.onload = () => {
                    idle_timeout = setTimeout(() => {
                        monkey_state.sprite = Sprite.Idle
                        monkey_state.offset = IDLE_OFFSET
                    }, CLAP_TIME)
                    monkey_image.onload = null
                }

                break
            }
        }
    })

    client.on<ChatResponseMessage>('chat', message => {
        if (!is_me(message.row, message.column))
            return
        show_chat_message(message.message)
    })
})
</script>

<template>
    <img
        :src="`/monkeys/${ monkey_state.sprite }.png`"
        ref="monkey_image_ref"
        class="monkey"
        alt="Moonkie"
        draggable="false" />

    <dev v-if="chat_message_enabled" class="message-box">
        <div class="message">
            <img src="/icons/chat-bottom.svg" draggable="false" />
            <text>{{ chat_message.message }}</text>
        </div>
    </dev>
</template>

<style scoped>
.monkey {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);

    bottom: v-bind('`${ props.monkey.bottom + monkey_state.offset }vh`');
    height: v-bind('`${ props.monkey.height * 1.1 }vh`');
    margin-left: v-bind('`${ props.monkey.x_offset }vh`');

    filter: brightness(0.5);
}

.message-box {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);

    bottom: v-bind('`${ props.monkey.bottom + props.monkey.height + 1.7 }vh`');
    margin-left: calc(v-bind('`${ props.monkey.x_offset }vh`'));
    margin-bottom: v-bind('chat_message.show ? "0" : "-0.5em"');

    transition:
            opacity v-bind('`${ chat_message.animation_speed }s`'),
            margin-bottom v-bind('`${ chat_message.animation_speed }s`');
    opacity: v-bind("chat_message.show ? 1 : 0");
}

.message {
    display: flex;
    position: absolute;
    padding: 0.8em;

    max-width: 20em;
    max-height: 10em;
    min-width: 5em;

    width: max-content;
    height: auto;
    left: 0;
    bottom: 0;

    border-radius: 1em;
    background-color: #FFFB;
}

.message text {
    display: inline-block;
    overflow-y: auto;

    font-size: 1.6em;
    word-wrap: break-word;
}

.message img {
    position: absolute;
    height: 1.5em;

    color: #FFFB;
    bottom: -1.5em;
}
</style>
