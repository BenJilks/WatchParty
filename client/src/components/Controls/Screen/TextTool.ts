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

    public on_click(x: number, y: number) {
        const text = document.createElement('text')
        text.innerText = 'Hello, Tools!!!'
        text.style.position = 'absolute'
        text.style.top = `${ y }px`
        text.style.left = `${ x }px`
        text.style.color = 'white'
        this.screen.append(text)
    }

}
