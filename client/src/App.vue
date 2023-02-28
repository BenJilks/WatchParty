<script setup lang="ts">
import SeatingArea from '@/components/Stage/Seats/SeatingArea.vue'
import StageBackground from '@/components/Stage/StageBackground.vue'
import ControlsPanel from '@/components/Controls/ControlsPanel.vue'
import StageScreen from '@/components/Stage/StageScreen.vue'
import type { SocketClient } from "@/socket_client";
import { inject, ref } from 'vue'

const screen = ref<InstanceType<typeof StageScreen>>()
const client_future = inject<Promise<SocketClient>>('client_future')!!

function zoom() {
    let zoom = 1

    window.onwheel = event => {
        if (!event.shiftKey && !event.ctrlKey)
            return

        zoom -= event.deltaY * (1.0 / window.innerHeight) * 0.4
        const scales: { [name: string]: number } = {
            '--seat-height': 20 * zoom,
            '--floor-height': 23 * zoom,
            '--curtain-offset': -(2.0 / 9.0 * 50 * (1.0 / zoom)) + 3,
            '--front-board-height': 6 * zoom,
        }

        for (const prop in scales) {
            document.body.style.setProperty(prop, `${scales[prop]}vh`)
        }
    }
}

zoom()
</script>

<template>
    <div class="background">
        <StageScreen ref="screen" />
        <StageBackground />
        <SeatingArea :client_future="client_future" />
        <ControlsPanel
            :screen="ref(screen)"
            :client_future="client_future" />
    </div>
</template>

<style scoped>
    .background {
        position: fixed;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;

        background-color: #050505;
    }
</style>
