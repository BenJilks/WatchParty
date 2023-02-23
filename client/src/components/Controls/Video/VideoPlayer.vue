<script setup lang="ts">
import { ref } from 'vue'

const video_ref = ref<HTMLVideoElement>()
defineExpose({
  video_ref,
})
</script>

<template>
  <video ref="video_ref" preload="auto" oncontextmenu="return false;">
    <source :src="'/vids/[Rhythm Heaven] - Fan Club (Perfect) (English)-DNbvktlB0gU.mp4'" type="video/mp4">
  </video>
</template>

<style scoped>
  video {
    position: absolute;
    width: 100%;
    height: 100%;
    object-fit: contain;
    z-index: 1;
  }
</style>
