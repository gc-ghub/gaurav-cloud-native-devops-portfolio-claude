import { writable } from 'svelte/store'

const STORAGE_KEY = 'portfolio-theme'

function getInitialTheme() {
  if (typeof localStorage === 'undefined') return 'dark'
  return localStorage.getItem(STORAGE_KEY) === 'light' ? 'light' : 'dark'
}

export const theme = writable(getInitialTheme())

theme.subscribe((value) => {
  if (typeof document === 'undefined') return
  document.documentElement.classList.toggle('light', value === 'light')
  localStorage.setItem(STORAGE_KEY, value)
})

export function toggleTheme() {
  theme.update((value) => (value === 'dark' ? 'light' : 'dark'))
}
