
export interface MonkeyData {
    seat: number,
    bottom: number,
    height: number,
    x_offset: number,
    your_token?: string,
}

export enum MonkeyAction {
    Ready = 'ready',
    Clap = 'clap',

    LeanLeft = 'lean-left',
    LeanRight = 'lean-right',
    NoLean = 'no-lean',

    Jump = 'jump',
}

export interface MonkeyActionMessage {
    action: MonkeyAction,
    token: string,
}

export interface MonkeyActionResponseMessage {
    action: MonkeyAction,
    row: number,
    column: number,
}

const monkey_seat_data_table = [
    { bottom: 14, height: 7, seat_offsets:
            [-45, -39.5, -34, -28, -23, -17.5, -12, -6, -1, 5, 11, 16.5, 22, 33, 38.5, 44] },
    { bottom: 12, height: 8, seat_offsets:
            [-45.2, -39, -32.8, -26.5, -20.2, -13.8, -7.5, -1, 5.4, 11.5, 18, 24.3, 30.7, 37, 43.3, 49.5] },
    { bottom: 11, height: 8.5, seat_offsets:
            [-45.5, -38.5, -31, -23.7, -16.3, -9, -1.4, 5.2, 12.5, 20, 27.4, 34.6, 42, 49.5] },
    { bottom: 8.5, height: 9.5, seat_offsets:
            [-46, -37.5, -28.8, -20.3, -11.5, -2.8, 5.5, 14.5, 23, 31.8, 40.5, 49] },
    { bottom: 5, height: 12, seat_offsets:
            [-46.5, -36, -25, -15, -4, 6, 17, 27.5, 38, 49] },
    { bottom: 0, height: 15, seat_offsets:
            [-48, -34.5, -21, -7, 7, 20.2, 34.5, 48] },
    { bottom: -8, height: 19, seat_offsets:
            [-50, -31, -11, 8, 27.5, 47.5] },
]

export function seats_in_row(row: number): number {
    if (row < 0 || row >= monkey_seat_data_table.length) {
        return 0
    }

    const { seat_offsets } = monkey_seat_data_table[row]
    return seat_offsets.length
}

export function create_monkey(row: number, seat: number, your_token?: string): MonkeyData {
    const { bottom, height, seat_offsets } = monkey_seat_data_table[row]
    return {
        seat: seat,
        bottom: bottom,
        height: height,
        x_offset: seat_offsets[seat],
        your_token: your_token,
    }
}
