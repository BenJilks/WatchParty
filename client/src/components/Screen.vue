<script setup lang="ts">
import { ref } from 'vue'

const video_ref = ref<HTMLVideoElement | null>(null)
const synchronising = ref(true)
const needs_focus = ref(false)

defineExpose({
  set_synchronising: (value: boolean) => synchronising.value = value,
  set_needs_focus: (value: boolean) => needs_focus.value = value,
  video_ref,
})
</script>

<template>
  <div class="screen">
    <video ref="video_ref">
      <source :src="'/vids/[Rhythm Heaven] - Fan Club (Perfect) (English)-DNbvktlB0gU.mp4'" type="video/mp4">
    </video>
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
  }

  .screen video {
    position: absolute;
    width: 100%;
    height: 100%;
    object-fit: contain;
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
    z-index: 1;

    pointer-events: none;
  }
</style>
