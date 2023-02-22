import type SubMenu from '@/components/Controls/SubMenu/SubMenu.vue'

export type RatioButtonClick = (event: Event) => void

export class RatioButtons {

  private readonly sub_menus: SubMenu[]
  private selected_button?: HTMLElement
  private selected_sub_menu?: SubMenu

  public constructor() {
    this.sub_menus = []
  }

  private hide_other_sub_menus(sub_menu: SubMenu) {
    for (const other_sub_menu of this.sub_menus) {
      if (other_sub_menu == sub_menu || !other_sub_menu.show)
        continue
      other_sub_menu.toggle()
    }
  }

  public add(sub_menu: SubMenu): RatioButtonClick {
    this.sub_menus.push(sub_menu)
    return (event: Event) => {
      if (this.selected_sub_menu == sub_menu) {
        this.close_current()
        return
      }

      this.selected_button?.classList.remove('selected')
      sub_menu.toggle()

      this.hide_other_sub_menus(sub_menu)
      this.selected_button = event.target as HTMLElement
      this.selected_sub_menu = sub_menu
      this.selected_button.classList.add('selected')
    }
  }

  public close_current() {
    this.selected_button?.classList.remove('selected')
    this.selected_button = undefined

    this.selected_sub_menu?.toggle()
    this.selected_sub_menu = undefined
  }

  public is_sub_menu_open(): boolean {
    return this.selected_sub_menu != undefined
  }

}
