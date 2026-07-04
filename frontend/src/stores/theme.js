import { writable } from 'svelte/store'

const STORAGE_KEY = 'portfolio-theme'
export const THEME_ORDER = ['green', 'dark', 'light']

function getInitialTheme() {
  if (typeof localStorage === 'undefined') return 'green'
  const stored = localStorage.getItem(STORAGE_KEY)
  return THEME_ORDER.includes(stored) ? stored : 'green'
}

export const theme = writable(getInitialTheme())

theme.subscribe((value) => {
  if (typeof document === 'undefined') return
  document.documentElement.classList.toggle('light', value === 'light')
  document.documentElement.classList.toggle('green', value === 'green')
  localStorage.setItem(STORAGE_KEY, value)
})

export function nextTheme(current) {
  return THEME_ORDER[(THEME_ORDER.indexOf(current) + 1) % THEME_ORDER.length]
}

export function toggleTheme() {
  theme.update(nextTheme)
}
