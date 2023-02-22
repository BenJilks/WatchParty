import AnnotationTool from '@/components/Controls/Screen/AnnotationTool'
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

    public on_click(x: number, y: number) {
        console.log(x, y)
    }

}
