
export type MonkeyData = {
    seat: number,
    bottom: number,
    height: number,
    x_offset: number,
}

const monkey_seat_data_table = [
    { bottom: 14, height: 7, seat_offsets:
            [-45, -39.5, -34, -28, -23, -17.5, -12, -6, -1, 5, 11, 16.5, 22, 33, 38.5, 44] },
    { bottom: 12, height: 8, seat_offsets:
            [-45.2, -39, -32.8, -26.5, -20.2, -13.8, -7.5, -1, 5.4, 11.5, 18, 24.3, 30.7, 37, 43.3, 49.5] },
    { bottom: 11, height: 8.5, seat_offsets:
            [-45.5, -38.5, -32, -24.2, -16.8, -9.3, -1.9, 5.2, 12.4, 20, 27, 34.4, 42, 49.5] },
    { bottom: 10, height: 9, seat_offsets:
            [-46, -38, -29.3, -20.3, -11.5, -2.8, 5.5, 14, 23, 31.8, 40.5, 49.5] },
    { bottom: 7, height: 10, seat_offsets:
            [-47, -36.8, -26, -15, -4, 6, 17, 27.5, 37.5, 48] },
    { bottom: 2, height: 12, seat_offsets:
            [-48, -34.5, -21, -7, 7, 20.2, 34.5, 48] },
    { bottom: -5, height: 16, seat_offsets:
            [-50, -31, -11, 8, 27.5, 47.5] },
]

export function seats_in_row(row: number): number {
    if (!(row in monkey_seat_data_table)) {
        return 0
    }

    const { seat_offsets } = monkey_seat_data_table[row]
    return seat_offsets.length
}

export function create_monkey(row: number, seat: number): MonkeyData {
    const { bottom, height, seat_offsets } = monkey_seat_data_table[row]
    return {
        seat: seat,
        bottom: bottom,
        height: height,
        x_offset: seat_offsets[seat],
    }
}
