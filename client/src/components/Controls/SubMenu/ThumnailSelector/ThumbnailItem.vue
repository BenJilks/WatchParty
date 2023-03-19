<script setup lang="ts">
export interface ItemData {
    name: string,
    item_file: string,
    thumbnail_file: string,
}

interface Props {
    item: ItemData,
}

interface Emits {
    (e: 'selected', item_file: string, name: string): void,
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

function selected() {
    emit('selected', props.item.item_file, props.item.name)
}
</script>

<template>
    <div class="item" @click="selected">
        <div id="content">
            <img
                :src="`/thumbnails/${ props.item.thumbnail_file }`"
                :alt="props.item.name"
                draggable="false" />
            <text>{{ props.item.name }}</text>
        </div>
    </div>
</template>

<style scoped>
.item {
    display: flex;
    align-content: center;
    justify-content: center;

    width: 100%;
    padding: 1em;
    height: var(--thumbnail-size);
}

.item:hover {
    border-radius: 0.7em;
    background-color: white;
    cursor: pointer;
}

.item #content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5em;

    width: var(--thumbnail-size);
    height: 100%;
    overflow: hidden;
}

.item img {
    width: 100%;
    height: auto;
    aspect-ratio: 5/3;

    object-fit: cover;
    border-radius: 0.5em;
}

.item text {
    display: inline-block;
    width: 100%;

    font-family: "Roboto", "Arial", sans-serif;
    font-weight: 500;
}
</style>
