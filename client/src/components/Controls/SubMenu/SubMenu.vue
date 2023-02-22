<script setup lang="ts">
import {computed, ref} from 'vue'

interface Props {
  height: string,
}

defineProps<Props>()
const show = ref(false)
const disabled = ref(true)
const enabled = ref(false)

function toggle() {
  if (!show.value) {
    disabled.value = false
    enabled.value = true
    setTimeout(() => show.value = true, 0)
  } else {
    show.value = false
    enabled.value = false
    setTimeout(() => disabled.value = true, 200)
  }
}

defineExpose({
  toggle,
  enabled,
})
</script>

<template>
  <div v-if="!disabled" class="panel" id="menu">
    <slot></slot>
  </div>
</template>

<style scoped>
#menu {
  position: absolute;
  bottom: calc(var(--controls-height) + var(--floating-y));
  align-items: start;

  width: 100%;
  height: v-bind('height');
  padding: 0.2em;

  transition:
      opacity 0.2s,
      transform 0.2s;
  opacity: v-bind('show ? 1 : 0');
  transform: translateY(v-bind('`${ show ? 0 : 4 }em`'));
}
</style>
