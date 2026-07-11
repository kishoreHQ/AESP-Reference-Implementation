import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import App from './app/App'
import './index.css'
import { detectProfile } from './profiles/detect'
import { mountMissionControl, mountHermesUI } from './embed'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 2000,
      refetchOnWindowFocus: false,
    },
  },
})

async function bootstrap() {
  const profile = detectProfile()
  if (profile.useMocks) {
    const { worker } = await import('./mocks/browser')
    await worker.start({
      onUnhandledRequest: 'bypass',
      serviceWorker: { url: '/mockServiceWorker.js' },
    })
  }

  const rootEl = document.getElementById('root')
  if (!rootEl) throw new Error('root missing')

  createRoot(rootEl).render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </QueryClientProvider>
    </StrictMode>,
  )
}

// P3 embed entry is exported; default SPA bootstraps on #root
void bootstrap()

export { mountMissionControl, mountHermesUI }
