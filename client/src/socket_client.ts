
function code_point_length(byte: number) {
    if ((byte & 0b11000000) == 0b10000000) return 0
    if ((byte & 0b11100000) == 0b11000000) return 2
    if ((byte & 0b11110000) == 0b11100000) return 3
    if ((byte & 0b11111000) == 0b11110000) return 4
    return 1
}

function decode_utf8_base64(encoded_data: string): string {
    let decoded_string = ''

    const decoded_bytes = atob(encoded_data)
    for (let i = 0; i < decoded_bytes.length; i++) {
        const byte = decoded_bytes.charCodeAt(i)
        const length = code_point_length(byte)
        if (length == 0) {
            console.error(`Invalid codepoint at index ${ i }`)
            continue
        }

        if (length == 1) {
            decoded_string += String.fromCharCode(byte)
            continue
        }

        let codepoint = (byte & 0xff >> length+1) << (length - 1) * 6
        for (let j = length - 2; j >= 0; --j)
            codepoint |= (decoded_bytes.charCodeAt(++i) & 0b00111111) << (j * 6)
        decoded_string += String.fromCharCode(codepoint)
    }

    return decoded_string
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
        this.socket.addEventListener('message', async event => {
            const message = JSON.parse(await event.data.text())
            if (!('type' in message))
                return

            const listeners = this.message_listeners.get(message.type)
            if (listeners === undefined)
                return

            const text_data = decode_utf8_base64(message.data)
            console.log(`Got '${ message.type }': ${ text_data }`)

            const data = JSON.parse(text_data)
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
        try {
            this.socket.send(JSON.stringify({
                type: type,
                data: btoa(JSON.stringify(message)),
            }))
        } catch (error) {
            console.error(error)
        }
    }
}

export function open_socket_client(): Promise<SocketClient> {
    return new Promise((resolve, reject) => {
        const address = `${ document.domain }:${ window.location.port }/socket`
        let socket = new WebSocket(`wss://${ address }`)

        socket.onopen = () => { resolve(new SocketClient(socket)) }
        socket.onerror = () => {
            socket.close()
            socket = new WebSocket(`ws://${ address }`)
            socket.onopen = () => { resolve(new SocketClient(socket)) }
            socket.onerror = error => { reject(error) }
        }

        setTimeout(() => {
            reject(new Error('Connection attempt timeout'))
        }, 1000)
    })
}
