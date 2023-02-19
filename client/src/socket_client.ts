
export interface Message {
    type: string,
}

export class SocketClient {
    private socket: WebSocket
    private message_listeners: Map<string, Map<number, (data: object) => void>>
    private listener_id_to_type_map: Map<number, string>

    constructor(socket: WebSocket) {
        this.socket = socket
        this.message_listeners = new Map()
        this.listener_id_to_type_map = new Map()
        this.start_message_listener()
    }
    private start_message_listener() {
        this.socket.addEventListener('message', event => {
            console.log(`Got message ${ event.data }`)

            const message = JSON.parse(event.data)
            if (!('type' in message))
                return

            const listeners = this.message_listeners.get(message.type)
            if (listeners === undefined)
                return

            const data = JSON.parse(atob(message.data))
            for (const listener of listeners.values())
                listener(data)
        })
    }
    private generateID(): number {
        for (;;) {
            const id = Math.random() * Number.MAX_VALUE
            if (!this.listener_id_to_type_map.has(id)) {
                return id
            }
        }
    }
    public on<T extends object>(type: string, callback: (t: T) => void): number {
        const id = this.generateID()

        if (!this.message_listeners.has(type))
            this.message_listeners.set(type, new Map())

        this.listener_id_to_type_map.set(id, type)
        this.message_listeners.get(type)?.set(id, data => {
            callback(data as T)
        })

        return id
    }
    public off(id: number) {
        const type = this.listener_id_to_type_map.get(id)
        if (type === undefined)
            return

        const listeners = this.message_listeners.get(type)
        if (listeners === undefined)
            return

        this.listener_id_to_type_map.delete(id)
        listeners.delete(id)
    }

    public send<T>(type: string, message: T) {
        this.socket.send(JSON.stringify({
            type: type,
            data: btoa(JSON.stringify(message)),
        }))
    }

}

export function open_socket_client(): Promise<SocketClient> {
    return new Promise((resolve, reject) => {
        const address = `${ document.domain }:${ window.location.port }/socket`
        const socket = new WebSocket(`wss://${ address }`)

        socket.onopen = () => { resolve(new SocketClient(socket)) }
        socket.onerror = () => {
            const socket = new WebSocket(`ws://${ address }`)
            socket.onopen = () => { resolve(new SocketClient(socket)) }
            socket.onerror = error => { reject(error) }
        }

        setTimeout(() => {
            reject(new Error('Connection attempt timeout'))
        }, 1000)
    })
}
