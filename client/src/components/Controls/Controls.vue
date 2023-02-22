<script setup lang="ts">
  import VideoControls from '@/components/Controls/Video/VideoControls.vue'
  import Chat from '@/components/Controls/Chat.vue'
  import VideoMenu from '@/components/Controls/SubMenu/VideoMenu.vue'
  import Screen from '@/components/Controls/Screen/Screen.vue'
  import AnnotationMenu from '@/components/Controls/SubMenu/AnnotationMenu.vue'
  import SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'
  import type { Ref } from 'vue'
  import type { RatioButtonClick } from '@/components/Controls/SubMenu/RatioButtons'
  import { SocketClient } from '@/socket_client'
  import { computed, onMounted, ref } from 'vue'
  import { RatioButtons } from '@/components/Controls/SubMenu/RatioButtons'

  interface Props {
    screen_ref: Ref<Screen | null>,
    client_future: Promise<SocketClient>,
  }

  const video_controls = ref<VideoControls>()
  const select_menu_open_ref = ref(false)
  defineProps<Props>()

  const video_menu = ref<VideoMenu>()
  const annotations_menu = ref<AnnotationMenu>()
  const toggle_video_menu = ref<RatioButtonClick>()
  const toggle_annotation_menu = ref<RatioButtonClick>()
  const ratio_buttons = new RatioButtons<SubMenu>()

  function change_video(video_file: string) {
    console.log(`Selected video: '${ video_file }'`)
    video_controls.value?.change_video(video_file)
  }

  async function video_selected(video_file: string) {
    change_video(video_file)
    ratio_buttons.close_current()
  }

  const controls_ref = ref<HTMLDivElement | null>(null)
  const controls_indicator_ref = ref<HTMLDivElement | null>(null)
  function set_controls_visible(visible: boolean) {
    if (!visible && video_controls.value?.volume_slider_open())
      return
    if (!visible && video_controls.value?.get_is_seeking())
      return
    if (!visible && ratio_buttons.is_any_selected())
      return

    const controls = controls_ref.value!
    const translation = visible
        ? '0'
        : `calc(var(--controls-height) + var(--floating-y))`
    controls.style.transform = `translateY(${ translation })`

    const controls_indicator = controls_indicator_ref.value!
    controls_indicator.style.opacity = `${ visible ? 0 : 0.4 }`
  }

  onMounted(async () => {
    toggle_video_menu.value = ratio_buttons.add(video_menu.value!)
    toggle_annotation_menu.value = ratio_buttons.add(annotations_menu.value!)

    window.addEventListener('mousemove', event => {
      const controls = controls_ref.value!
      const height = controls.getBoundingClientRect().height * 6
      const show_controls = (event.screenY >= window.innerHeight - height)
      set_controls_visible(show_controls)
    })
  })
</script>

<template>
  <div ref="controls_indicator_ref" class="controls-indicator"></div>
  <div ref="controls_ref" class="controls">
    <div class="panel">
      <Chat :client_future="client_future" />
    </div>

    <div id="video-panel" class="panel">
      <VideoControls
          ref="video_controls"
          :screen_ref="computed(() => screen_ref.value)"
          :video_ref="computed(() => screen_ref.value?.video_ref)"
          :client_future="client_future" />
    </div>

    <div class="panel">
      <img
          class="icon"
          id="annotate-button"
          src="/icons/edit.png"
          @click="toggle_annotation_menu"
          alt="" />
      <img
          class="icon"
          id="menu-button"
          src="/icons/up.svg"
          @click="toggle_video_menu"
          alt="" />
    </div>

    <VideoMenu
        ref="video_menu"
        :client_future="client_future"
        @selected="video_selected" />
    <AnnotationMenu
        :tools="computed(() => screen_ref.value.tools)"
        ref="annotations_menu" />
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

  #menu-button {
    transition: transform 0.2s;
    transform: rotate(v-bind('`${ select_menu_open_ref ? 180 : 0 }deg`'))
  }
</style>
