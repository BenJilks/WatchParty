<script setup lang="ts">
import type {MonkeyActionMessage, MonkeyActionResponseMessage, MonkeyData} from '@/monkey'
import {MonkeyAction} from '@/monkey'
import type {SocketClient} from '@/socket_client'
import {reactive, ref} from 'vue'

const CHAT_MESSAGE_TIME_MS = 2500
const CLAP_TIME = 300

const IDLE_OFFSET = 1
const CLAP_READY_OFFSET = 3
const CLAP_OFFSET = 3
const LEAN_ANGLE_DEGREES = 15

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

enum LeanState {
    None,
    Left,
    Right,
}

interface MonkeyState {
    sprite: Sprite,
    offset: number,
    lean: number,
    lean_state: LeanState,
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
    lean: 0,
    lean_state: LeanState.None,
})

const chat_message_enabled = ref(false)
const chat_message = reactive<ChatBubble>({
    message: '',
    show: false,
    animation_speed: 0.2,
})

async function send_action(action: MonkeyAction) {
    const client = await props.client_future
    client.send<MonkeyActionMessage>('monkey-action', {
        action: action,
        token: props.monkey.your_token ?? '',
    })
}

let idle_timeout: number | undefined = undefined

function clap_down() {
    monkey_state.sprite = Sprite.Ready
    monkey_state.offset = CLAP_READY_OFFSET
    send_action(MonkeyAction.Ready)
    clearTimeout(idle_timeout)
}

function clap_up() {
    monkey_state.sprite = Sprite.Clap
    monkey_state.offset = CLAP_OFFSET
    send_action(MonkeyAction.Clap)

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
    window.addEventListener('keydown', event => {
        switch (event.code) {
            case 'Space':
                clap_down()
                break

            case 'ArrowLeft':
                monkey_state.lean = -LEAN_ANGLE_DEGREES
                monkey_state.lean_state = LeanState.Left
                break

            case 'ArrowRight':
                monkey_state.lean = LEAN_ANGLE_DEGREES
                monkey_state.lean_state = LeanState.Right
                break
        }
    })
    window.addEventListener('keyup', event => {
        switch (event.code) {
            case 'Space':
                clap_up()
                break

            case 'ArrowLeft':
                if (monkey_state.lean_state == LeanState.Left) {
                    monkey_state.lean = 0
                    monkey_state.lean_state = LeanState.None
                }
                break

            case 'ArrowRight':
                if (monkey_state.lean_state == LeanState.Right) {
                    monkey_state.lean = 0
                    monkey_state.lean_state = LeanState.None
                }
                break
        }
    })
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
    client.on<MonkeyActionResponseMessage>('monkey-action', message => {
        if (!is_me(message.row, message.column))
            return

        switch (message.action) {
            case MonkeyAction.Ready: {
                monkey_state.sprite = Sprite.Ready
                monkey_state.offset = CLAP_READY_OFFSET
                clearTimeout(idle_timeout)
                break
            }

            case MonkeyAction.Clap: {
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
    transform:
        translateX(calc(-50% + v-bind('`${ 7 * monkey_state.lean / LEAN_ANGLE_DEGREES }%`')))
        rotate(v-bind('`${ monkey_state.lean }deg`'));

    bottom: v-bind('`${ props.monkey.bottom + monkey_state.offset }vh`');
    height: v-bind('`${ props.monkey.height * 1.1 }vh`');
    margin-left: v-bind('`${ props.monkey.x_offset }vh`');

    filter: brightness(0.5);
    transition: transform 0.2s;
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
