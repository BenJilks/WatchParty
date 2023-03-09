<script setup lang="ts">
import { getCookieNumber, setCookie } from '@/cookie'
import { computed, ref } from 'vue'

interface Emits {
    (e: 'volume_change', volume: number): void
}
const emit = defineEmits<Emits>()

const volume = ref(getCookieNumber("volume") ?? 0.5)
const muted = ref(false)
const over_slider = ref(false)
const volume_bar_ref = ref<HTMLDivElement | null>(null)
const is_dragging = ref(false)

function toggle_mute() {
    muted.value = !muted.value
    emit('volume_change', muted.value ? 0 : volume.value)
}

function mouse_enter() {
    over_slider.value = true
}

function mouse_leave() {
    over_slider.value = false
}

function set_volume(value: number) {
    volume.value = Math.min(Math.max(value, 0), 1)
    setCookie('volume', volume.value.toString())
    emit('volume_change', volume.value)
}

function mouse_wheel(event: WheelEvent) {
    if (muted.value && event.deltaY >= 0) {
        return
    }

    muted.value = false
    set_volume(volume.value - event.deltaY * 0.0005)
}

window.addEventListener('mousemove', (event: MouseEvent) => {
    if (volume_bar_ref.value == null || event.buttons != 1) {
        return
    }
    is_dragging.value = true
    muted.value = false

    const volume_bar = volume_bar_ref.value
    const rect = volume_bar.getBoundingClientRect()
    const y = (event.clientY - rect.y) / rect.height
    set_volume(1 - y)
})

window.addEventListener('mouseup', () => {
    is_dragging.value = false
})

const show_slider = computed(() => over_slider.value || is_dragging.value)
defineExpose({
    show_slider,
})
</script>

<template>
    <div
        class="icon"
        @mouseenter="mouse_enter"
        @mouseleave="mouse_leave"
        @wheel="mouse_wheel">

        <div id="slider" @mouseleave="mouse_leave">
            <div v-if="show_slider" id="bar" ref="volume_bar_ref">
                <div id="marker"></div>
                <div id="done"></div>
            </div>
        </div>

        <img
            id="volume"
            draggable="false"
            :src="`/icons/${ (muted || volume === 0) ? 'mute' : 'volume' }.svg`"
            @mouseenter="mouse_enter"
            @click="toggle_mute"
            alt="volume" />
    </div>
</template>

<style scoped>
.icon {
    position: relative;
    cursor: default;
    padding: 0 0.6em;
    width: 4em;

    pointer-events: all;
}

.icon img {
    position: relative;
    width: 100%;
    padding: 0 0.2em;

    object-fit: contain;
    cursor: pointer;
    pointer-events: all;
    z-index: 101;
}

#slider {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 0.2em;
    padding: 0.2em 0;

    position: absolute;
    width: 100%;
    height: auto;
    min-height: 100%;
    left: 0;
    bottom: 0;

    border-radius: 1em;
    background-color: #fffe;

    cursor: pointer;
    z-index: 100;
}

#volume {
    height: 100%;
    aspect-ratio: 1/1;
}

#bar {
    position: relative;
    width: 0.6em;
    height: 8em;
    margin-top: 1em;
    margin-bottom: 100%;

    border-radius: 4em;
    background-color: #0006;
}

#marker {
    position: absolute;
    width: 200%;
    aspect-ratio: 1/1;

    left: -50%;
    bottom: v-bind('`${ muted ? 0 : 100*volume }%`');
    transform: translateY(50%);

    border-radius: 50%;
    background-color: black;
}

#done {
    position: absolute;
    width: 100%;
    height: v-bind('`${ muted ? 0 : 100*volume }%`');
    bottom: 0;

    border-radius: 4em;
    background-color: black;
}
</style>
