<script setup lang="ts">
import BufferedSegment from '@/components/Controls/Video/BufferedSegment.vue'
import VolumeSlider from '@/components/Controls/Video/VolumeSlider.vue'
import type StageScreen from '@/components/Stage/StageScreen.vue'
import type { SocketClient } from '@/socket_client'
import type { BufferedSegmentData, PlaybackData } from '@/components/Controls/Video/VideoSync'
import type { Ref } from 'vue'
import SyncedVideo from '@/components/Controls/Video/VideoSync'
import { computed, reactive, ref, watch } from 'vue'

interface Props {
    screen: Ref<InstanceType<typeof StageScreen> | undefined>,
    video: Ref<HTMLVideoElement | undefined>,
    client_future: Promise<SocketClient>,
}

let is_seeking = false
const buffered_segments = reactive<BufferedSegmentData[]>([])
const seek_bar_ref = ref<HTMLDivElement>()
const time = ref<HTMLElement>()
const playback_data = reactive<PlaybackData>({
    progress: 0,
    duration: 0.1,
    playing: false,
    syncing: true,
})

const props = defineProps<Props>()
const progress_factor = computed(() =>
        playback_data.progress / playback_data.duration)

function update_buffer_segments() {
    buffered_segments.splice(0)
    buffered_segments.push(...synced_video?.get_updated_buffered_segments() ?? [])
}

let synced_video: SyncedVideo | null = null
let buffer_updater = ref<number>()

watch(props.video, async () => {
    const client = await props.client_future
    const video = props.video.value!!
    const screen = props.screen.value!!
    synced_video = new SyncedVideo(client, video, screen, ref(playback_data))

    // NOTE: Ensure there's only ever a single interval
    clearInterval(buffer_updater.value)
    buffer_updater.value = setInterval(update_buffer_segments, 200)
    video.addEventListener('progress', update_buffer_segments)
})

async function play_pause() {
    if (playback_data.syncing)
        return
    await synced_video?.send_toggle_play_pause()
}

function on_seek_start() {
    if (is_seeking || playback_data.syncing) {
        return
    }

    props.video.value?.pause()
    is_seeking = true
}

function volume_change(volume: number) {
    synced_video?.change_volume(volume)
}

window.addEventListener('mouseup', async () => {
    if (!is_seeking) {
        return
    }

    is_seeking = false
    await synced_video?.force_sync()
})

window.addEventListener('mousemove', event => {
    if (!is_seeking)
        return

    const seek_bar = seek_bar_ref.value
    const bounding_rect = seek_bar?.getBoundingClientRect()
    const x_along_seek_bar = event.x - (bounding_rect?.x ?? 0) - (1.5 / 2)*16

    const progress = x_along_seek_bar / (bounding_rect?.width ?? 1)
    synced_video?.seek(progress)
})

async function change_video(video_path: string) {
    await synced_video?.send_request_play({
        playing: true,
        progress: 0,
        video: video_path,
    })
}

function format_time(time_in_ms: number): string {
    const seconds = Math.floor(time_in_ms) % 60
    const minutes = Math.floor(time_in_ms / 60) % 60
    const hours = Math.floor(time_in_ms / (60 * 60)) % 60
    const format_component = (component: number) =>
        component >= 10 ? component : `0${ component }`

    const result = `${format_component(minutes)}:${format_component(seconds)}`
    if (hours > 0) {
        return `${format_component(hours)}:${result}`
    } else {
        return result
    }
}

const time_text = computed(() =>
    `${ format_time(playback_data.progress) } / ${ format_time(playback_data.duration) }`)

const volume = ref<typeof VolumeSlider>()
const volume_slider_open = () => volume.value?.show_slider
const get_is_seeking = () => is_seeking

defineExpose({
    change_video,
    volume_slider_open,
    get_is_seeking,
})
</script>

<template>
    <img
        class="icon"
        draggable="false"
        :src="`/icons/${ playback_data.playing ? 'pause' : 'play' }.svg`"
        @click="play_pause"
        alt="play" />

    <div id="timeline" ref="seek_bar_ref" @mousedown="on_seek_start">
        <BufferedSegment
            v-for="(segment, i) in buffered_segments"
            :key="i"
            :start="segment.start"
            :end="segment.end" />
        <div id="done-timeline"></div>
        <div id="scrubber" @mousedown="on_seek_start"></div>
    </div>

    <text id="time">{{ time_text }}</text>
    <VolumeSlider ref="volume" @volume_change="volume_change" />
</template>

<style scoped>
* {
    --scrubber-size: 1.5em;
    --unplayed-color: #4445;
    --played-color: #333F;
}

#timeline {
    position: relative;
    width: 100%;
    height: 0.5em;

    border-radius: 1em;
    background-color: var(--unplayed-color);
}

#done-timeline {
    position: absolute;
    width: calc(100% * v-bind('progress_factor'));
    height: 100%;

    border-radius: 1em;
    background-color: var(--played-color);
}

#scrubber {
    position: absolute;
    width: var(--scrubber-size);
    height: var(--scrubber-size);

    top: 0.25em;
    transform: translateY(-50%);
    left: calc(100% * v-bind('progress_factor'));

    border-radius: 50%;
    background-color: black;

    cursor: pointer;
    pointer-events: all;
    user-modify: none;
}

#time {
    display: inline-block;
    width: auto;
    word-break: keep-all;
    white-space: nowrap;
}
</style>
