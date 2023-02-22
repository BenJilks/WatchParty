<script setup lang="ts">
  export type VideoData = {
    name: string,
    video_file: string,
    thumbnail_file: string,
  }

  interface Props {
    video: VideoData,
  }

  interface Emits {
    (e: 'selected', video_file: string): void,
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  function selected() {
    emit('selected', props.video.video_file)
  }
</script>

<template>
  <div class="item" @click="selected">
    <div id="content">
      <img :src="`/thumbnails/${ props.video.thumbnail_file }`" draggable="false" alt="" />
      <text>{{ props.video.name }}</text>
    </div>
  </div>
</template>

<style scoped>
  .item {
    display: flex;
    align-content: center;
    justify-content: center;

    width: 100%;
    padding: 1em;
    height: 15em;
  }

  .item:hover {
    border-radius: 0.7em;
    background-color: white;
    cursor: pointer;
  }

  .item #content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5em;

    width: 14em;
    height: 100%;
    overflow: hidden;
  }

  .item img {
    width: 100%;
    aspect-ratio: 3/5;

    height: 10em;
    object-fit: cover;
    border-radius: 0.5em;
  }

  .item text {
    display: inline-block;
    width: 100%;

    font-family: "Roboto", "Arial", sans-serif;
    font-weight: 500;
  }
</style>
