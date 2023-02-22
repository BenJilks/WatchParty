export type RatioButtonClick = (event: Event) => void

interface Tool {
  toggle: () => void
  enabled: boolean
}

export class RatioButtons<T extends Tool> {

  private selected_button?: HTMLElement
  private selected_tool?: T

  public add(tool: T): RatioButtonClick {
    return (event: Event) => {
      if (this.selected_tool != tool && this.selected_tool !== undefined)
        this.close_current()
      tool.toggle()

      this.selected_tool = tool
      this.selected_button = event.target as HTMLElement
      if (tool.enabled)
        this.selected_button.classList.add('selected')
      else
        this.selected_button.classList.remove('selected')
    }
  }

  public close_current() {
    if (this.selected_tool?.enabled) {
      this.selected_button?.classList.remove('selected')
      this.selected_tool?.toggle()
    }

    this.selected_button = undefined
    this.selected_tool = undefined
  }

  public is_any_selected(): boolean {
    return this.selected_tool != undefined
  }

}
