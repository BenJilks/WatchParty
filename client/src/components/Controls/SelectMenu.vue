<script setup lang="ts">
  import { ref } from 'vue'
  import VideoItem from "@/components/Controls/VideoItem.vue";

  interface Emits {
    (e: 'selected', video: string): void,
  }

  const emit = defineEmits<Emits>()
  const show = ref(false)
  const disabled = ref(true)

  function toggle() {
    if (show.value) {
      show.value = false
      setTimeout(() => disabled.value = true, 200)
    } else {
      disabled.value = false
      setTimeout(() => show.value = true, 0)
    }
  }

  function selected() {
    emit('selected', 'test')
  }

  defineExpose({
    toggle,
  })
</script>

<template>
  <div v-if="!disabled" class="panel" id="menu">
    <div id="content-list">
      <VideoItem
          v-for="i in 20"
          title="Whats Your Problem - FULL VIDEO ARCHIVE-MlrfQJJbfsQ.mp4"
          @selected="selected"
          :key="i" />
    </div>
  </div>
</template>

<style scoped>
  #menu {
    position: absolute;
    bottom: calc(var(--controls-height) + var(--floating-y));

    width: 100%;
    height: 50vh;
    padding: 2em;

    transition:
        opacity 0.2s,
        transform 0.2s;
    opacity: v-bind('show ? 1 : 0');
    transform: translateY(v-bind('`${ show ? 0 : 4 }em`'));
  }

  #content-list {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    border-radius: 1em;

    display: grid;
    grid-auto-flow: row;
    grid-template-columns: repeat(auto-fill, minmax(15em, 1fr));
    gap: 0.5em;
  }
</style>
