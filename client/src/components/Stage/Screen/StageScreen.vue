<script setup lang="ts">
import VideoPlayer from '@/components/Stage/Screen/VideoPlayer.vue'
import ImageViewer from '@/components/Stage/Screen/ImageViewer.vue'
import type { ImageData } from '@/components/Stage/Screen/ImageViewer.vue'
import { computed, ref } from 'vue'

const screen = ref<HTMLDivElement>()
const video_player = ref<typeof VideoPlayer>()
const displayed_image = ref<ImageData | undefined>(undefined)

const overlay_message = ref<string | undefined>()
const hide_overlay_message = ref(false)
const show_overlay_message = computed(() =>
    !hide_overlay_message.value && overlay_message.value !== undefined)

defineExpose({
    set_overlay_message: (value: string | undefined) => overlay_message.value = value,
    set_hide_overlay_message: (value: boolean) => hide_overlay_message.value = value,
    set_displayed_image: (image?: ImageData) => displayed_image.value = image,
    video: computed(() => video_player.value?.video),
})
</script>

<template>
    <div class="screen" ref="screen">
        <ImageViewer :image="displayed_image" />
        <VideoPlayer ref="video_player" />

        <div id="shadow-overlay"></div>
        <div class="overlay">
            {{ overlay_message }}
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
    border-radius: 1vh;

    overflow: hidden;
    box-shadow: 0 0 5vh 2vh #6666;
}

.screen #shadow-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;

    box-shadow: inset 0 0 0.5em 0.25em #000;
    border-radius: 1vh;
    border: 0.2em solid #000;

    pointer-events: none;
    z-index: 3;
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

    color: #dddf;
    background-color: #222a;

    pointer-events: none;
    user-select: none;

    transition: opacity 0.3s;
    opacity: v-bind('show_overlay_message ? 1 : 0');
    z-index: 1;
}
</style>
