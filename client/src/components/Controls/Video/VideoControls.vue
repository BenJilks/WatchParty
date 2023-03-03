<script setup lang="ts">
import BufferedSegment from '@/components/Controls/Video/BufferedSegment.vue'
import VolumeSlider from '@/components/Controls/Video/VolumeSlider.vue'
import type StageScreen from '@/components/Stage/StageScreen.vue'
import type { SocketClient } from '@/socket_client'
import type { Ref } from 'vue'
import { computed, reactive, ref, watch } from 'vue'

interface Props {
    screen: Ref<InstanceType<typeof StageScreen> | undefined>,
    video: Ref<HTMLVideoElement | undefined>,
    client_future: Promise<SocketClient>,
}

interface BufferedSegmentData {
    start: number,
    end: number,
}

type RequestPlayMessage = {
    playing: boolean,
    progress: number,
    video?: string,
}

let is_seeking = false
const buffered_segments = reactive<BufferedSegmentData[]>([])
const seek_bar_ref = ref<HTMLDivElement>()
const data = reactive({
    progress: 0,
    duration: 1,
    playing: false,
    muted: false,
    syncing: false,
})

const props = defineProps<Props>()
const progress_factor = computed(() =>
        data.progress / data.duration)

function update_buffered_segments(video: HTMLVideoElement) {
    buffered_segments.splice(0, buffered_segments.length)

    for (let i = 0; i < video.buffered.length; i++) {
        const start = video.buffered.start(i)
        const end = video.buffered.end(i)

        const length = end - start
        if (length / data.duration > 0.1) {
            buffered_segments.push({start: start, end: end})
        }
    }
}

function set_syncing(value: boolean) {
    data.syncing = value
    if (value) {
        props.screen.value?.set_overlay_message('Syncing Playback...')
    } else {
        props.screen.value?.set_overlay_message(undefined)
    }
}

async function set_needs_focus(error: DOMException) {
    if (error.name != 'NotAllowedError')
        return

    const give_focus = async () => {
        const video = props.video.value

        window.removeEventListener('click', give_focus)
        await video?.play()
        data.playing = true

        client.send('ready', null)
        props.screen.value?.set_overlay_message(undefined)
    }

    const client = await props.client_future
    client.send('request-play', {
        playing: data.playing,
        progress: data.progress,
    })

    window.addEventListener('click', give_focus)
    props.screen.value?.set_hide_overlay_message(false)
    props.screen.value?.set_overlay_message('Click Here to Play')
}

async function send_when_ready() {
    const video = props.video.value
    if (video?.src === undefined || video.src == '') {
        props.screen.value?.set_overlay_message('No Video Selected')
        return
    }

    async function ready_to_play() {
        const client = await props.client_future
        client.send('ready', null)
        video?.removeEventListener('canplaythrough', ready_to_play)
    }

    video?.addEventListener('canplaythrough', ready_to_play)
    if ((video?.readyState ?? 0) >= HTMLMediaElement.HAVE_FUTURE_DATA) {
        await ready_to_play()
    }
}

async function handle_request_play(message: RequestPlayMessage) {
    const video = props.video.value
    if (video === undefined) {
        return
    }

    if (!video.paused) {
        video.pause()
    }
    set_syncing(true)

    if (message.video ?? '' != '') {
        video.src = `/vids/${ message.video }`
    }
    video.currentTime = message.progress
    data.playing = message.playing
    data.progress = message.progress
    await send_when_ready()
}

async function send_request_play(message: RequestPlayMessage) {
    const client = await props.client_future
    client.send('request-play', message)
    await handle_request_play(message)
}

let buffer_updater = ref<number>()
watch(props.video, async () => {
    const video = props.video.value!!
    video.preload = 'auto'

    video.addEventListener('durationchange', () => {
        data.duration = video.duration
    })

    video.addEventListener('progress', () => {
        update_buffered_segments(video)
    })

    // NOTE: Ensure there's only ever a single interval
    clearInterval(buffer_updater.value)
    buffer_updater.value = setInterval(() =>
        update_buffered_segments(video), 200)

    video.addEventListener('timeupdate', () => {
        if (!is_seeking)
            data.progress = video.currentTime
    })

    video.addEventListener('waiting', () => {
        if (data.syncing) {
            return
        }

        send_request_play({
            playing: data.playing,
            progress: data.progress,
        })
    })

    const client = await props.client_future
    client.on<RequestPlayMessage>('request-play', handle_request_play)
    client.on('ready', () => {
        data.playing
                ? video.play().catch(set_needs_focus)
                : video.pause()
        set_syncing(false)
    })
})

async function play_pause() {
    if (data.syncing)
        return

    await send_request_play({
        playing: !data.playing,
        progress: data.progress,
    })
}

function client_seek(progress: number) {
    const video = props.video.value
    if (video === undefined)
        return

    if ('fastSeek' in video) {
        video.fastSeek(progress)
    } else {
        // NOTE: Chrome does not have a `fastSeek` function
        (video as any).currentTime = progress
    }
}

function on_seek_start() {
    if (is_seeking || data.syncing) {
        return
    }

    props.video.value?.pause()
    is_seeking = true
}

window.addEventListener('mouseup', async () => {
    if (!is_seeking) {
        return
    }

    is_seeking = false
    await send_request_play({
        playing: data.playing,
        progress: data.progress,
    })
})

window.addEventListener('mousemove', event => {
    if (!is_seeking)
        return

    const seek_bar = seek_bar_ref.value
    const bounding_rect = seek_bar?.getBoundingClientRect()
    const x_along_seek_bar = event.x - (bounding_rect?.x ?? 0) - (1.5 / 2)*16

    const progress = x_along_seek_bar / (bounding_rect?.width ?? 1)
    if (progress < 0) {
        data.progress = 0
    } else if (progress > 1) {
        data.progress = data.duration
    } else {
        data.progress = data.duration * progress
    }

    client_seek(data.progress)
})

function volume_change(volume: number) {
    const video = props.video.value
    if (video === undefined) {
        return
    }

    video.volume = volume
    video.muted = (volume == 0)
}

async function change_video(video_path: string) {
    await send_request_play({
        playing: true,
        progress: 0,
        video: video_path,
    })
}

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
        :src="`/icons/${ data.playing ? 'pause' : 'play' }.svg`"
        @click="play_pause"
        alt="play" />

    <div id="timeline" ref="seek_bar_ref" @mousedown="on_seek_start">
        <BufferedSegment
            v-for="(segment, i) in buffered_segments"
            :key="i"
            :start="segment.start / data.duration"
            :end="segment.end / data.duration" />
        <div id="done-timeline"></div>
        <div id="scrubber" @mousedown="on_seek_start"></div>
    </div>

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
</style>
