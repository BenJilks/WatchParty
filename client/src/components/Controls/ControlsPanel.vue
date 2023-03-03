<script setup lang="ts">
import VideoMenu from '@/components/Controls/SubMenu/VideoMenu.vue'
import AnnotationMenu from '@/components/Controls/SubMenu/AnnotationMenu.vue'
import VideoControls from '@/components/Controls/Video/VideoControls.vue'
import ChatBox from '@/components/Controls/ChatBox.vue'
import type StageScreen from '@/components/Stage/StageScreen.vue'
import type { Ref } from 'vue'
import type { SocketClient } from '@/socket_client'
import { computed, onMounted, ref } from 'vue'
import { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'

interface Props {
    screen: Ref<InstanceType<typeof StageScreen> | undefined>,
    client_future: Promise<SocketClient>,
}

const ratio_buttons = ref(new RatioButtons())
const controls = ref<HTMLDivElement>()
const controls_indicator = ref<HTMLDivElement>()

const video_controls = ref<typeof VideoControls>()
const props = defineProps<Props>()

function change_video(video_file: string) {
    console.log(`Selected video: '${ video_file }'`)
    video_controls.value?.change_video(video_file)
}

async function video_selected(video_file: string) {
    change_video(video_file)
    ratio_buttons.value.close_current()
}

const saved_overlay = ref<string>()
function set_controls_visible(visible: boolean) {
    if (controls.value === undefined || controls_indicator.value === undefined) {
        return
    }

    if (!visible && video_controls.value?.volume_slider_open())
        return
    if (!visible && video_controls.value?.get_is_seeking())
        return
    if (!visible && ratio_buttons.value.is_any_selected())
        return

    const translation = visible
            ? '0'
            : `calc(var(--controls-height) + var(--floating-y))`
    controls.value.style.transform = `translateY(${ translation })`
    controls_indicator.value.style.opacity = `${ visible ? 0 : 0.4 }`
}

onMounted(async () => {
    window.addEventListener('mousemove', event => {
        const height = (controls.value?.getBoundingClientRect()?.height ?? 0) * 6
        const show_controls = (event.screenY >= window.innerHeight - height)
        set_controls_visible(show_controls)
    })

    window.addEventListener('click', () => {
        props.screen.value?.set_hide_overlay_message(
            ratio_buttons.value.is_any_selected())
    })
})
</script>

<template>
    <div ref="controls_indicator" class="controls-indicator"></div>
    <div ref="controls" class="controls">
        <div class="panel">
            <ChatBox :client_future="client_future" />
        </div>

        <div id="video-panel" class="panel">
            <VideoControls
                ref="video_controls"
                :screen="computed(() => screen.value)"
                :video="computed(() => screen.value?.video)"
                :client_future="client_future" />
        </div>

        <div class="panel">
            <AnnotationMenu
                :tools="computed(() => screen.value?.tools)"
                :ratio_buttons="computed(() => ratio_buttons)" />
            <VideoMenu
                :client_future="client_future"
                :ratio_buttons="computed(() => ratio_buttons)"
                @selected="video_selected" />
        </div>
    </div>
</template>

<style scoped>
* {
    --floating-x: 3em;
    --floating-y: 2em;
    --controls-height: 3em;
}

.controls {
    position: fixed;
    display: flex;
    align-items: center;
    gap: 2em;

    height: var(--controls-height);
    width: calc(100% - var(--floating-x)*2);
    bottom: var(--floating-y);
    left: var(--floating-x);

    transition: transform 0.2s;
    transform: translateY(calc(var(--controls-height) + var(--floating-y)));
}

.controls-indicator {
    position: fixed;
    height: 2em;
    bottom: -1em;
    left: 10em;
    right: 10em;

    border-radius: 1em;
    background-color: #fff;

    transition: opacity 0.2s;
    opacity: 0.3;
}

#video-panel {
    width: 100%;
}
</style>
