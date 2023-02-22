import type { CursorProperty } from 'csstype'

export default abstract class AnnotationTool {

    private is_enabled: boolean
    protected readonly screen: HTMLDivElement

    public constructor(screen: HTMLDivElement) {
        this.is_enabled = false
        this.screen = screen
    }

    public toggle() {
        if (!this.enabled)
            console.log(`Selected tool ${ this.name }`)
        else
            console.log(`Deselected tool ${ this.name }`)
        this.is_enabled = !this.is_enabled
    }

    public get enabled(): boolean {
        return this.is_enabled
    }

    public abstract get cursor(): CursorProperty
    public abstract get name(): string
    public abstract get icon(): string

    public abstract on_click(x: number, y: number): void

}
