<script setup lang="ts">
import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
import ThumbnailSelector from '@/components/Controls/SubMenu/ThumnailSelector/ThumbnailSelector.vue'
import type { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'
import type { ItemData } from '@/components/Controls/SubMenu/ThumnailSelector/ThumbnailItem.vue'
import type { SocketClient } from '@/socket_client'
import type { Ref } from 'vue'
import {computed, onMounted, reactive} from 'vue'

interface Props {
    ratio_buttons: Ref<RatioButtons<any>>,
    client_future: Promise<SocketClient>
}

interface ImageListMessage {
    images: ItemData[],
}

interface Emits {
    (e: 'selected', picture: string, name: string): void,
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const image_list = reactive<ItemData[]>([])

function selected(picture_file: string, name: string) {
    emit('selected', picture_file, name)
}

async function request_image_list() {
    if (image_list.length > 0)
        return

    const client = await props.client_future
    client.on<ImageListMessage>('image-list', message => {
        image_list.splice(0, image_list.length)
        image_list.push(...message.images)
    })
    client.send('image-list', null)
}

onMounted(() => {
    request_image_list()
})
</script>

<template>
    <SubMenu
        ref="sub_menu"
        height="70vh"
        :ratio_buttons="computed(() => ratio_buttons.value)"
        icon="gallery.svg">

        <div id="content">
            <ThumbnailSelector
                :item_list="image_list"
                @selected="selected" />
        </div>
    </SubMenu>
</template>

<style scoped>
#content {
    flex-direction: column;
    gap: 0.5em;

    width: 100%;
    height: 100%;
    padding: 2em 1em;
    overflow-y: hidden;
}
</style>
