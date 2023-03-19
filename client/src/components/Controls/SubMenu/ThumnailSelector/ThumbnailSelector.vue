<script setup lang="ts">
import ThumbnailItem from '@/components/Controls/SubMenu/ThumnailSelector/ThumbnailItem.vue'
import type { ItemData } from '@/components/Controls/SubMenu/ThumnailSelector/ThumbnailItem.vue'

interface Props {
    item_list: ItemData[],
}

interface Emits {
    (e: 'selected', item_file: string, name: string): void,
}

defineProps<Props>()
const emit = defineEmits<Emits>()

function selected(item_file: string, name: string) {
    emit('selected', item_file, name)
}
</script>

<template>
    <div id="thumbnail-list">
        <ThumbnailItem
            v-for="item in item_list"
            :item="item"
            :key="item.name"
            @selected="selected" />
    </div>
</template>

<style scoped>
* {
    --thumbnail-size: 13em;
}

#thumbnail-list {
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

#thumbnail-list::-webkit-scrollbar {
    display: none;
}
</style>
