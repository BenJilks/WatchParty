export type RatioButtonClick = (event: Event) => void

interface Tool {
  toggle: () => void,
  enabled: boolean,
}

export class RatioButtons<T extends Tool> {

  private readonly tools: T[]
  private selected_button?: HTMLElement
  private selected_tool?: T

  public constructor() {
    this.tools = []
  }

  private hide_other_sub_menus(tool: T) {
    for (const other_tool of this.tools) {
      if (other_tool == tool || !other_tool.enabled)
        continue
      other_tool.toggle()
    }
  }

  public add(tool: T): RatioButtonClick {
    this.tools.push(tool)
    return (event: Event) => {
      if (this.selected_tool == tool) {
        this.close_current()
        return
      }

      this.selected_button?.classList.remove('selected')
      tool.toggle()

      this.hide_other_sub_menus(tool)
      this.selected_button = event.target as HTMLElement
      this.selected_tool = tool
      this.selected_button.classList.add('selected')
    }
  }

  public close_current() {
    this.selected_button?.classList.remove('selected')
    this.selected_button = undefined

    this.selected_tool?.toggle()
    this.selected_tool = undefined
  }

  public is_any_selected(): boolean {
    return this.selected_tool != undefined
  }

}
