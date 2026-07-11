import type { HostEvent } from '@/shared/types/host'
import { detectProfile } from '@/profiles/detect'

type Handler = (ev: HostEvent) => void

/** Single owner of realtime stream (UI-TEC-03, UI-RT-01). */
class EventBridgeImpl {
  private ws: WebSocket | null = null
  private handlers = new Set<Handler>()
  private seq = 0
  private missionFilter?: string
  private reconnectMs = 500
  private timer: ReturnType<typeof setTimeout> | null = null
  private mockTimer: ReturnType<typeof setInterval> | null = null
  private connected = false

  get lastSeq() {
    return this.seq
  }

  get isConnected() {
    return this.connected
  }

  subscribe(h: Handler) {
    this.handlers.add(h)
    return () => {
      this.handlers.delete(h)
    }
  }

  connect(missionId?: string) {
    this.missionFilter = missionId
    this.teardown()
    const profile = detectProfile()
    if (profile.useMocks) {
      this.startMockStream()
      return
    }
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const base = profile.baseUrl.replace(/^http/, 'ws').replace(/\/$/, '')
    // Prefer same-origin proxy path
    const url =
      base.startsWith('ws') && !base.includes('://')
        ? `${proto}://${window.location.host}${base}/v1/events?since=${this.seq}${missionId ? `&mission=${missionId}` : ''}`
        : `${proto}://${window.location.host}/api/v1/events?since=${this.seq}${missionId ? `&mission=${missionId}` : ''}`
    try {
      this.ws = new WebSocket(url)
      this.ws.onopen = () => {
        this.connected = true
        this.reconnectMs = 500
      }
      this.ws.onmessage = (m) => {
        try {
          const ev = JSON.parse(m.data as string) as HostEvent
          if (typeof ev.seq === 'number') this.seq = Math.max(this.seq, ev.seq)
          this.dispatch(ev)
        } catch {
          /* ignore */
        }
      }
      this.ws.onclose = () => {
        this.connected = false
        this.scheduleReconnect()
      }
      this.ws.onerror = () => {
        this.ws?.close()
      }
    } catch {
      this.startMockStream()
    }
  }

  private scheduleReconnect() {
    if (this.timer) clearTimeout(this.timer)
    this.timer = setTimeout(() => this.connect(this.missionFilter), this.reconnectMs)
    this.reconnectMs = Math.min(this.reconnectMs * 1.6, 8000)
  }

  private startMockStream() {
    this.connected = true
    if (this.mockTimer) clearInterval(this.mockTimer)
    this.mockTimer = setInterval(() => {
      this.seq += 1
      this.dispatch({
        seq: this.seq,
        type: 'mission.updated',
        ts: new Date().toISOString(),
        missionId: this.missionFilter ?? 'mis_demo_running',
        data: { tick: this.seq },
      })
    }, 4000)
  }

  private dispatch(ev: HostEvent) {
    this.handlers.forEach((h) => h(ev))
  }

  teardown() {
    if (this.timer) clearTimeout(this.timer)
    if (this.mockTimer) clearInterval(this.mockTimer)
    this.mockTimer = null
    this.ws?.close()
    this.ws = null
  }
}

export const EventBridge = new EventBridgeImpl()
