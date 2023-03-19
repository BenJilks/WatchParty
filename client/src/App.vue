<script setup lang="ts">
import SeatingArea from '@/components/Stage/Seats/SeatingArea.vue'
import StageBackground from '@/components/Stage/StageBackground.vue'
import ControlsPanel from '@/components/Controls/ControlsPanel.vue'
import StageScreen from '@/components/Stage/Screen/StageScreen.vue'
import type { SocketClient } from "@/socket_client";
import { inject, ref } from 'vue'

const screen = ref<InstanceType<typeof StageScreen>>()
const client_future = inject<Promise<SocketClient>>('client_future')!!
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
