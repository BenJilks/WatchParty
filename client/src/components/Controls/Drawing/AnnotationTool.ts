import type { CursorProperty } from 'csstype'

export interface Position {
    x: number,
    y: number,
}

export default abstract class AnnotationTool {

    private is_enabled: boolean
    protected readonly context: CanvasRenderingContext2D

    public constructor(context: CanvasRenderingContext2D) {
        this.is_enabled = false
        this.context = context
    }

    public toggle() {
        this.is_enabled = !this.is_enabled
    }

    public get enabled(): boolean {
        return this.is_enabled
    }

    public abstract get cursor(): CursorProperty
    public abstract get name(): string
    public abstract get icon(): string

    public on_mouse_down(_mouse: Position) {}
    public on_mouse_up(_mouse: Position) {}
    public on_mouse_move(_mouse: Position) {}

}
