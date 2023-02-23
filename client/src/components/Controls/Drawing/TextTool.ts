import AnnotationTool from '@/components/Controls/Drawing/AnnotationTool'
import type { Position } from '@/components/Controls/Drawing/AnnotationTool'
import type { CursorProperty } from 'csstype'

export default class TextTool extends AnnotationTool {

    public get name(): CursorProperty {
        return 'text';
    }

    public get cursor(): CursorProperty {
        return 'text';
    }

    public get icon(): string {
        return 'text.svg'
    }

    public override on_mouse_down(position: Position) {
        this.context.lineWidth = 10
        this.context.fillStyle = '#fff'
        this.context.fillText('Test', position.x, position.y)
        this.context.fill()
    }

}
