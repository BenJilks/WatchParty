<script setup lang="ts">
import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
import type AnnotationTool from '@/components/Controls/Drawing/AnnotationTool'
import type { RatioButtonClick } from '@/components/Controls/SubMenu/RatioButtons'
import type { Ref } from 'vue'
import { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'
import { computed, ref, watch } from 'vue'

interface Props {
    tools: Ref<{[name: string]: AnnotationTool}>,
    ratio_buttons: Ref<RatioButtons<any>>,
}

interface ToolButton {
    callback: RatioButtonClick,
    tool: AnnotationTool,
}

const sub_menu = ref<typeof SubMenu>()
const props = defineProps<Props>()
const tool_callbacks = ref<ToolButton[]>()

const ratio_buttons = new RatioButtons()
watch(props.tools, () => {
    tool_callbacks.value = Object.keys(props.tools.value)
        .map(key => props.tools.value[key])
        .map(tool => ({
            callback: ratio_buttons.add(tool),
            tool: tool,
        }))
})

watch(sub_menu, () => {
    const button = sub_menu.value?.button
    watch(button, () => {
        button.value?.addEventListener('click', () =>
            ratio_buttons.close_current())
    })
})
</script>

<template>
    <SubMenu
        ref="sub_menu"
        height="auto"
        :ratio_buttons="computed(() => props.ratio_buttons.value)"
        icon="edit.png">

        <div id="annotation-menu">
            <img
                v-for="(button, i) in tool_callbacks"
                :key="i"
                :src="`/icons/${ button.tool.icon }`"
                class="icon"
                draggable="false"
                @click="button.callback"
                :alt="button.tool.name" />
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
