<script setup lang="ts">
  import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
  import AnnotationTool from '@/components/Controls/Screen/AnnotationTool'
  import type { RatioButtonClick } from '@/components/Controls/SubMenu/RatioButtons'
  import type { Ref } from 'vue'
  import { computed, onMounted, ref } from 'vue'
  import { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'

  interface Props {
    tools: Ref<{[name: string]: AnnotationTool}>,
  }

  interface ToolButton {
    callback: RatioButtonClick,
    tool: AnnotationTool,
  }

  const sub_menu = ref<SubMenu>()
  const props = defineProps<Props>()
  const tool_callbacks = ref<ToolButton[]>()

  const ratio_buttons = new RatioButtons()
  onMounted(() => {
    tool_callbacks.value = Object.keys(props.tools.value)
        .map(key => props.tools.value[key])
        .map(tool => ({
          callback: ratio_buttons.add(tool),
          tool: tool,
        }))
  })

  function toggle() {
    ratio_buttons.close_current()
    sub_menu.value?.toggle()
  }

  const enabled = computed(() => sub_menu.value?.enabled)
  defineExpose({
    toggle,
    enabled,
  })
</script>

<template>
  <SubMenu ref="sub_menu" height="auto">
    <div id="annotation-menu">
      <img
          v-for="(button, i) in tool_callbacks"
          :key="i"
          :src="`/icons/${ button.tool.icon }`"
          class="icon"
          draggable="false"
          @click="button.callback" />
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
