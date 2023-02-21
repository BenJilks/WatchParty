<script setup lang="ts">
  interface Props {
    title: string,
  }

  interface Emits {
    (e: 'selected'): void,
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  function selected() {
    emit('selected')
  }
</script>

<template>
  <div class="item" @click="selected">
    <div id="content">
      <img src="/thumbnails/test.png" draggable="false" alt="" />
      <text>{{ props.title }}</text>
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
