<script setup lang="ts">
import Drawing from '@/components/Controls/Drawing/DrawingCanvas.vue'
import VideoPlayer from '@/components/Controls/Video/VideoPlayer.vue'
import { computed, ref } from 'vue'

const screen = ref<HTMLDivElement>()
const video_player = ref<typeof VideoPlayer>()
const drawing = ref<typeof Drawing>()

const synchronising = ref(true)
const needs_focus = ref(false)

defineExpose({
    set_synchronising: (value: boolean) => synchronising.value = value,
    set_needs_focus: (value: boolean) => needs_focus.value = value,
    video: computed(() => video_player.value?.video),
    tools: computed(() => drawing.value?.tools),
})
</script>

<template>
    <div class="screen" ref="screen">
        <VideoPlayer ref="video_player" />
        <Drawing ref="drawing" />
        <div v-if="synchronising && !needs_focus" class="overlay">
            Synchronising Viewers...
        </div>
        <div v-if="needs_focus" class="overlay">
            Click Anywhere to Play
        </div>
    </div>
</template>

<style scoped>
.screen {
    display: flex;
    justify-content: center;
    align-items: center;

    position: absolute;
    bottom: calc(var(--seat-height) + 7vh);
    left: 50%;

    width: calc(100% - 25vw);
    height: calc(100% - var(--seat-height) - 9vh);
    transform: translateX(-50%);

    background-color: #0a0a0aff;
    box-shadow: 0 0 5vh 2vh #6666;
    border-radius: 0.5vh;

    overflow: hidden;
}

.screen * {
    text-wrap: none;
    word-break: keep-all;
}

.screen .overlay {
    padding: 0.5em 1em;
    border-radius: 0.5em;
    width: auto;
    height: auto;

    font-family: 'Trebuchet MS', sans-serif;
    font-size: 3em;
    font-weight: bold;

    color: white;
    background-color: #000a;
    z-index: 3;

    pointer-events: none;
    user-select: none;
}
</style>
