<script setup lang="ts">
import type { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'
import type { Ref } from 'vue'
import { onMounted, ref } from 'vue'

interface Props {
    height: string,
    ratio_buttons: Ref<RatioButtons<any>>,
    icon: string,
}

const show = ref(false)
const disabled = ref(true)

const button = ref<HTMLImageElement>()
const props = defineProps<Props>()

const tool = ref({
    toggle: toggle,
    enabled: false,
})

function toggle() {
    if (!show.value) {
        disabled.value = false
        tool.value.enabled = true
        setTimeout(() => show.value = true, 0)
    } else {
        show.value = false
        tool.value.enabled = false
        setTimeout(() => disabled.value = true, 200)
    }
}

onMounted(() => {
    button.value?.addEventListener('click',
        props.ratio_buttons.value.add(tool.value))
})

defineExpose({
    button,
})
</script>

<template>
    <img
        ref="button"
        class="icon"
        id="menu-button"
        :src='`/icons/${ icon }`'
        alt="" />

    <div v-if="!disabled" class="panel" id="menu">
        <slot></slot>
    </div>
</template>

<style scoped>
#menu {
    position: absolute;
    left: 0;
    width: 100%;
    height: v-bind('height');
    bottom: calc(var(--controls-height) + var(--floating-y));

    padding: 0.2em;
    align-items: start;

    transition:
        opacity 0.2s,
        transform 0.2s;
    opacity: v-bind('show ? 1 : 0');
    transform: translateY(v-bind('`${ show ? 0 : 4 }em`'));
}
</style>
