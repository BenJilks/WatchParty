<script setup lang="ts">
import VideoItem from '@/components/Controls/Video/VideoItem.vue'
import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
import type { VideoData } from '@/components/Controls/Video/VideoItem.vue'
import type { SocketClient } from '@/socket_client'
import type { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'
import type { Ref } from 'vue'
import {computed, reactive, ref} from 'vue'

interface Props {
    client_future: Promise<SocketClient>,
    ratio_buttons: Ref<RatioButtons<any>>,
}

interface Emits {
    (e: 'selected', video: string): void,
}

type VideoListMessage = {
    videos: VideoData[],
}

type AddYouTubeVideoMessage = {
    url: string,
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const sub_menu = ref<typeof SubMenu>()
const video_list = reactive<VideoData[]>([])

const youtube_url = ref<HTMLInputElement>()
const can_submit_youtube = ref(false)

async function request_video_list() {
    if (request_video_list.length > 0)
        return

    const client = await props.client_future
    client.on<VideoListMessage>('video-list', message => {
        video_list.splice(0, video_list.length)
        video_list.push(...message.videos)
    })
    client.send('video-list', null)
}

function selected(video_file: string) {
    emit('selected', video_file)
}

function is_valid_youtube_url(url: string): boolean {
    const valid_starts = [
        'http://youtube.com/watch',
        'http://youtube.com/watch',
        'https://www.youtube.com/watch',
        'https://www.youtube.com/watch',
    ]

    for (const start of valid_starts) {
        if (url.startsWith(start)) {
            return true
        }
    }

    return false
}

function video_entered() {
    const url = youtube_url.value?.value ?? ''
    const is_valid = is_valid_youtube_url(url)
    if (!is_valid) {
        can_submit_youtube.value = false
        return
    }

    can_submit_youtube.value = true
}

async function submit_youtube() {
    if (!(can_submit_youtube.value ?? false)) {
        return
    }

    const input = youtube_url.value
    if (input === undefined) {
        return
    }

    const url = input.value ?? ''
    const client = await props.client_future
    client.send<AddYouTubeVideoMessage>('add-youtube-video', {
        url: url,
    })

    input.value = ''
    can_submit_youtube.value = false
}

request_video_list()

for (let i = 0; i < 100; i++) {
    video_list.push({
        name: 'Testing',
        video_file: 'Stolen Chat (El no Sabe).webm',
        thumbnail_file: 'Stolen Chat (El no Sabe).jpg'
    })
}
</script>

<template>
    <SubMenu
        ref="sub_menu"
        height="50vh"
        :ratio_buttons="computed(() => ratio_buttons.value)"
        icon="up.svg">

        <div id="content">
            <div id="add-video">
                <img class="icon" src="/icons/youtube.svg" alt="YouTube">
                <input
                    placeholder="http://youtube.com/watch?id=..."
                    ref="youtube_url"
                    @change="video_entered"
                    @keyup="video_entered" />
                <img
                    id="add"
                    alt="Add"
                    :class="`icon ${ can_submit_youtube ? 'enabled' : 'disabled' }`"
                    :src="`/icons/add${ can_submit_youtube ? '' : '_disabled' }.svg`"
                    @click="submit_youtube">
            </div>
            <div id="video-list">
                <VideoItem
                    v-for="video in video_list"
                    :video="video"
                    :key="video.name"
                    @selected="selected" />
            </div>
        </div>
    </SubMenu>
</template>

<style scoped>
* {
    --thumbnail-size: 13em;
}

#content {
    display: flex;
    flex-direction: column;
    gap: 0.5em;

    width: 100%;
    height: 100%;
    padding: 0.5em;
    overflow-y: hidden;
}

#add-video {
    display: flex;
    gap: 1em;
    align-items: center;
    justify-content: end;

    height: 2.5em;
    padding: 0 1.5em;
}

#add-video input {
    width: 30em;
}

#add-video #add {
    height: 80%;
}

#add-video .enabled {
    cursor: pointer;
    pointer-events: all;
}

#add-video .disabled {
    cursor: default;
    pointer-events: none;
}

#video-list {
    width: auto;
    height: 100%;
    margin: 0 0.5em 0.5em 0.5em;

    overflow-y: auto;
    border-radius: 1em;

    display: grid;
    grid-auto-flow: row;
    grid-template-columns: repeat(auto-fill, minmax(var(--thumbnail-size), 1fr));
    grid-auto-rows: min-content;
    gap: 0.5em;

    box-shadow: inset 0 0 1em #0006;
}

#video-list::-webkit-scrollbar {
    display: none;
}
</style>
