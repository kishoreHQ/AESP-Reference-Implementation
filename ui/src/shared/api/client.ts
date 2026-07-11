import type { ApiError, Envelope } from '@/shared/types/host'
import { detectProfile } from '@/profiles/detect'

export class HostApiError extends Error {
  code: string
  remediation?: string
  constructor(err: ApiError) {
    super(err.message)
    this.code = err.code
    this.remediation = err.remediation
  }
}

function baseUrl() {
  return detectProfile().baseUrl.replace(/\/$/, '')
}

function authHeaders(): Record<string, string> {
  const headers: Record<string, string> = { Accept: 'application/json' }
  const token = detectProfile().token
  if (token) headers.Authorization = `Bearer ${token}`
  return headers
}

export async function apiGet<T>(path: string, params?: Record<string, string | undefined>): Promise<T> {
  const qs = new URLSearchParams()
  if (params) {
    Object.entries(params).forEach(([k, v]) => {
      if (v != null && v !== '') qs.set(k, v)
    })
  }
  const q = qs.toString()
  const res = await fetch(`${baseUrl()}${path}${q ? `?${q}` : ''}`, { headers: authHeaders() })
  return parse<T>(res)
}

export async function apiPost<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${baseUrl()}${path}`, {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: body == null ? undefined : JSON.stringify(body),
  })
  return parse<T>(res)
}

export async function apiPut<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${baseUrl()}${path}`, {
    method: 'PUT',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: body == null ? undefined : JSON.stringify(body),
  })
  return parse<T>(res)
}

async function parse<T>(res: Response): Promise<T> {
  let json: Envelope<T> | null = null
  try {
    json = (await res.json()) as Envelope<T>
  } catch {
    throw new HostApiError({
      code: 'parse_error',
      message: `Invalid JSON (${res.status})`,
      remediation: 'Check Host Interface connectivity and response format.',
    })
  }
  if (!res.ok || json.error) {
    throw new HostApiError(
      json.error ?? {
        code: `http_${res.status}`,
        message: res.statusText || 'Request failed',
        remediation: 'Retry or inspect daemon logs.',
      },
    )
  }
  return json.data
}
