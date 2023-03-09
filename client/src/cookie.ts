
export function getCookie(name: string): string | null {
    let component = document.cookie
        .split(';')
        .map(component => component.trim())
        .find(component => component.indexOf(`${ name }=`) == 0)
    if (component == undefined) {
        return null
    }

    const value = component.substring(name.length+1, component.length)
    return decodeURIComponent(value)
}

export function getCookieNumber(name: string): number | null {
    const value = getCookie(name)
    if (value == null) {
        return null
    }

    return parseFloat(value)
}

export function setCookie(name: string, value: string) {
    const cookies = document.cookie
        .split(';')
        .filter(component =>
            component.trim().indexOf(`${ name }=`) != 0)

    cookies.push(`${ name }=${ encodeURIComponent(value) }`)
    document.cookie = cookies.join(';')
}
