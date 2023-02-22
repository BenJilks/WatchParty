<script setup lang="ts">
  import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
  import {ref} from "vue";

  const sub_menu = ref<SubMenu>()

  const toggle = () => sub_menu.value?.toggle()
  defineExpose({
    toggle,
    sub_menu,
  })
</script>

<template>
  <SubMenu ref="sub_menu" height="auto">
    test
  </SubMenu>
</template>

<style scoped>
</style>
