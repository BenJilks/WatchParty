import AnnotationTool from '@/components/Controls/Screen/AnnotationTool'
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

}
