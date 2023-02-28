<script setup lang="ts">
import { computed, ref } from 'vue'

interface Emits {
    (e: 'volume_change', volume: number): void
}
const emit = defineEmits<Emits>()

const volume = ref(1)
const over_slider = ref(false)
const volume_bar_ref = ref<HTMLDivElement | null>(null)
const is_dragging = ref(false)

function toggle_mute() {
    volume.value = volume.value > 0 ? 0 : 1
    emit('volume_change', volume.value)
}

function mouse_enter() {
    over_slider.value = true
}

function mouse_leave() {
    over_slider.value = false
}

window.addEventListener('mousemove', (event: MouseEvent) => {
    if (volume_bar_ref.value == null || event.buttons != 1)
        return
    is_dragging.value = true

    const volume_bar = volume_bar_ref.value
    const rect = volume_bar?.getBoundingClientRect()
    const y = (event.clientY - rect.y) / rect.height

    volume.value = Math.min(Math.max(1 - y, 0), 1)
    emit('volume_change', volume.value)
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
    <div class="icon" @mouseenter="mouse_enter">
        <div id="slider" @mouseleave="mouse_leave">
            <div v-if="show_slider" id="bar" ref="volume_bar_ref">
                <div id="marker"></div>
                <div id="done"></div>
            </div>
        </div>

        <img
            id="volume"
            draggable="false"
            :src="`/icons/${ volume === 0 ? 'mute' : 'volume' }.svg`"
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
}

.icon img {
    position: relative;
    width: 100%;
    padding: 0 0.2em;

    object-fit: contain;
    cursor: pointer;
    pointer-events: all;
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
    pointer-events: all;
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
    bottom: v-bind('`${ 100*volume }%`');
    transform: translateY(50%);

    border-radius: 50%;
    background-color: black;
}

#done {
    position: absolute;
    width: 100%;
    height: v-bind('`${ 100*volume }%`');
    bottom: 0;

    border-radius: 4em;
    background-color: black;
}
</style>
