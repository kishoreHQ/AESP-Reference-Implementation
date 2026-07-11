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
  private sse: EventSource | null = null

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
    this.connectWS(missionId)
  }

  private connectWS(missionId?: string) {
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const qs = new URLSearchParams()
    qs.set('since', String(this.seq))
    if (missionId) qs.set('mission', missionId)
    const url = `${proto}://${window.location.host}/api/v1/events?${qs.toString()}`
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

  /** Mock stream simulates live mission progression for demos/tests. */
  private startMockStream() {
    this.connected = true
    if (this.mockTimer) clearInterval(this.mockTimer)
    const types = [
      'mission.updated',
      'node.updated',
      'log.append',
      'artifact.created',
      'memory.written',
    ] as const
    this.mockTimer = setInterval(() => {
      this.seq += 1
      const type = types[this.seq % types.length]
      this.dispatch({
        seq: this.seq,
        type,
        ts: new Date().toISOString(),
        missionId: this.missionFilter ?? 'mis_demo_running',
        data: { tick: this.seq, mock: true },
      })
    }, 2500)
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
    this.sse?.close()
    this.sse = null
  }
}

export const EventBridge = new EventBridgeImpl()
