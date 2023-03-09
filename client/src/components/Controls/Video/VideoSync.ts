import type StageScreen from '@/components/Stage/StageScreen.vue'
import type { SocketClient } from '@/socket_client'
import type {Ref} from "vue";

export interface PlaybackData {
    progress: number,
    duration: number,
    playing: boolean,
    syncing: boolean,
}

export interface BufferedSegmentData {
    start: number,
    end: number,
}

interface RequestPlayMessage {
    playing: boolean,
    progress: number,
    video?: string,
}

export default class SyncedVideo {
    private readonly video: HTMLVideoElement
    private readonly client: SocketClient
    private readonly screen: InstanceType<typeof StageScreen>
    private playback_data: Ref<PlaybackData>

    public constructor(client: SocketClient,
                       video: HTMLVideoElement,
                       screen: InstanceType<typeof StageScreen>,
                       playback_data: Ref<PlaybackData>) {
        this.video = video
        this.client = client
        this.screen = screen
        this.playback_data = playback_data
        video.preload = 'auto'

        video.addEventListener('durationchange', () => {
            this.playback_data.value.duration = video.duration
        })

        video.addEventListener('timeupdate', () => {
            this.playback_data.value.progress = video.currentTime
        })

        video.addEventListener('waiting', async () => {
            if (this.playback_data.value.syncing) {
                return
            }

            await this.send_request_play({
                playing: this.playback_data.value.playing,
                progress: this.playback_data.value.progress,
            })
        })

        client.on<RequestPlayMessage>('request-play', async message => {
            await this.handle_request_play(message)
        })

        client.on('ready', async () => {
            this.playback_data.value.playing
                ? video.play().catch(error => this.set_needs_focus(error))
                : video.pause()
            this.set_syncing(false)
        })
    }

    private async set_needs_focus(error: DOMException) {
        if (error.name != 'NotAllowedError')
            return

        const give_focus = async () => {
            window.removeEventListener('click', give_focus)
            await this.video.play()
            this.playback_data.value.playing = true
            this.client.send('ready', null)
            this.screen.set_hide_overlay_message(true)
        }

        this.client.send('request-play', {
            playing: this.playback_data.value.playing,
            progress: this.playback_data.value.progress,
        })

        window.addEventListener('click', give_focus)
        this.screen.set_hide_overlay_message(false)
        this.screen.set_overlay_message('Click Here to Play')
    }

    private async send_when_ready() {
        if (this.video.src == '') {
            this.screen.set_overlay_message('No Video Selected')
            return
        }

        const ready_to_play = async () => {
            this.client.send('ready', null)
            this.video.removeEventListener('canplaythrough', ready_to_play)
        }

        this.video.addEventListener('canplaythrough', ready_to_play)
        if ((this.video.readyState ?? 0) >= HTMLMediaElement.HAVE_FUTURE_DATA) {
            await ready_to_play()
        }
    }

    private async handle_request_play(message: RequestPlayMessage) {
        if (!this.video.paused) {
            this.video.pause()
        }
        this.set_syncing(true)

        if (message.video ?? '' != '') {
            this.video.src = `/vids/${ message.video }`
        }
        this.video.currentTime = message.progress
        this.playback_data.value.playing = message.playing
        this.playback_data.value.progress = message.progress
        await this.send_when_ready()
    }

    public async send_request_play(message: RequestPlayMessage) {
        this.client.send('request-play', message)
        await this.handle_request_play(message)
    }

    public async send_toggle_play_pause() {
        await this.send_request_play({
            playing: !this.playback_data.value.playing,
            progress: this.playback_data.value.progress,
        })
    }

    public async force_sync() {
        await this.send_request_play({
            playing: this.playback_data.value.playing,
            progress: this.playback_data.value.progress,
        })
    }

    public seek(progress: number) {
        if (progress < 0) {
            this.playback_data.value.progress = 0
        } else if (progress > 1) {
            this.playback_data.value.progress = this.playback_data.value.duration
        } else {
            this.playback_data.value.progress = this.playback_data.value.duration * progress
        }

        if ('fastSeek' in this.video) {
            this.video.fastSeek(this.playback_data.value.progress)
        } else {
            // NOTE: Chrome does not have a `fastSeek` function
            (this.video as any).currentTime = this.playback_data.value.progress
        }
    }

    public change_volume(volume: number) {
        this.video.volume = volume
        this.video.muted = (volume == 0)
    }

    public get_updated_buffered_segments(): BufferedSegmentData[] {
        const buffered_segments = []

        for (let i = 0; i < this.video.buffered.length; i++) {
            const start = this.video.buffered.start(i) / this.playback_data.value.duration
            const end = this.video.buffered.end(i) / this.playback_data.value.duration

            const length = end - start
            if (length > 0.1) {
                buffered_segments.push({start: start, end: end})
            }
        }

        return buffered_segments
    }

    private set_syncing(value: boolean) {
        this.playback_data.value.syncing = value
        if (value) {
            this.screen.set_overlay_message('Syncing Playback...')
            this.screen.set_hide_overlay_message(false)
        } else {
            this.screen.set_overlay_message(undefined)
            this.screen.set_hide_overlay_message(true)
        }
    }

}
