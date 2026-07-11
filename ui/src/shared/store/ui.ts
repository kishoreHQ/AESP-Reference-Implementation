import { create } from 'zustand'

type Theme = 'dark' | 'light'

interface UiState {
  theme: Theme
  railExpanded: boolean
  spineOpen: boolean
  density: 'comfortable' | 'compact'
  reducedMotion: boolean
  toasts: { id: string; message: string; tone: 'ok' | 'warn' | 'fail' | 'info' }[]
  setTheme: (t: Theme) => void
  toggleRail: () => void
  setSpineOpen: (v: boolean) => void
  pushToast: (message: string, tone?: UiState['toasts'][0]['tone']) => void
  dismissToast: (id: string) => void
}

export const useUi = create<UiState>((set) => ({
  theme: 'dark',
  railExpanded: false,
  spineOpen: true,
  density: 'comfortable',
  reducedMotion: false,
  toasts: [],
  setTheme: (theme) => {
    document.documentElement.classList.toggle('light', theme === 'light')
    document.documentElement.classList.toggle('dark', theme === 'dark')
    set({ theme })
  },
  toggleRail: () => set((s) => ({ railExpanded: !s.railExpanded })),
  setSpineOpen: (spineOpen) => set({ spineOpen }),
  pushToast: (message, tone = 'info') =>
    set((s) => ({
      toasts: [...s.toasts, { id: String(Date.now()), message, tone }].slice(-5),
    })),
  dismissToast: (id) => set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) })),
}))
