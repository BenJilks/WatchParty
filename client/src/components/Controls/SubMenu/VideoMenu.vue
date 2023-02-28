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

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const sub_menu = ref<typeof SubMenu>()
const video_list = reactive<VideoData[]>([])

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

request_video_list()
</script>

<template>
    <SubMenu
        ref="sub_menu"
        height="50vh"
        :ratio_buttons="computed(() => ratio_buttons.value)"
        icon="up.svg">

        <div id="content-list">
            <VideoItem
                v-for="video in video_list"
                :video="video"
                :key="video.name"
                @selected="selected" />
        </div>
    </SubMenu>
</template>

<style scoped>
#content-list {
    width: 100%;
    height: 100%;

    overflow-y: auto;
    border-radius: 1em;
    padding: 2em;

    display: grid;
    grid-auto-flow: row;
    grid-template-columns: repeat(auto-fill, minmax(15em, 1fr));
    grid-auto-rows: min-content;
    gap: 0.5em;
}
</style>
