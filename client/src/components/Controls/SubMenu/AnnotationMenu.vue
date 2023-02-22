<script setup lang="ts">
  import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
  import type { RatioButtonClick } from '@/components/Controls/SubMenu/RatioButtons'
  import { onMounted, ref } from 'vue'
  import { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'

  const sub_menu = ref<SubMenu>()
  const tool_a = ref<RatioButtonClick>()
  const tool_b = ref<RatioButtonClick>()
  const tool_c = ref<RatioButtonClick>()
  const tool_d = ref<RatioButtonClick>()

  const ratio_buttons = new RatioButtons()
  onMounted(() => {
    tool_a.value = ratio_buttons.add({ toggle: () => console.log('a'), enabled: false })
    tool_b.value = ratio_buttons.add({ toggle: () => console.log('b'), enabled: false })
    tool_c.value = ratio_buttons.add({ toggle: () => console.log('c'), enabled: false })
    tool_d.value = ratio_buttons.add({ toggle: () => console.log('d'), enabled: false })
  })

  const toggle = () => sub_menu.value?.toggle()
  defineExpose({
    toggle,
    sub_menu,
  })
</script>

<template>
  <SubMenu ref="sub_menu" height="auto">
    <div id="annotation-menu">
      <img src="/icons/text.svg" class="icon" draggable="false" @click="tool_a" />
      <img src="/icons/text.svg" class="icon" draggable="false" @click="tool_b" />
      <img src="/icons/text.svg" class="icon" draggable="false" @click="tool_c" />
      <img src="/icons/text.svg" class="icon" draggable="false" @click="tool_d" />
    </div>
  </SubMenu>
</template>

<style scoped>
  #annotation-menu {
    display: flex;
    gap: 1em;

    height: 3em;
    padding: 0.5em 1em;
  }
</style>
