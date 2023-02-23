import AnnotationTool from '@/components/Controls/Drawing/AnnotationTool'
import type { Position } from '@/components/Controls/Drawing/AnnotationTool'
import type { CursorProperty } from 'csstype'

export default class LineTool extends AnnotationTool {

    public get name(): CursorProperty {
        return 'line'
    }

    public get cursor(): CursorProperty {
        return 'pointer'
    }

    public get icon(): string {
        return 'up.svg'
    }

    public on_mouse_up({ x, y }: Position) {
        console.log('press', x, y)
    }

    public on_mouse_down({ x, y }: Position) {
        console.log('release', x, y)
    }

    public on_mouse_move({ x, y }: Position) {
        console.log('drag', x, y)
    }

}
