<script setup lang="ts">
export interface ImageData {
    file: string,
    name: string,
}

interface Props {
    image?: ImageData,
}

defineProps<Props>()
</script>

<template>
    <img
        :class="image === undefined ? 'hidden' : ''"
        :src="`/images/${ image?.file }`"
        :alt="image?.name"
        oncontextmenu="return false;" />
</template>

<style scoped>
img {
    position: absolute;
    width: 100%;
    height: 100%;
    object-fit: contain;
    z-index: 2;
}

.hidden {
    display: none;
}
</style>
