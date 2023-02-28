<script setup lang="ts">
import { ref } from 'vue'

const video = ref<HTMLVideoElement>()
defineExpose({
    video,
})
</script>

<template>
    <video ref="video" preload="auto" oncontextmenu="return false;">
        <source src="" />
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
