// WebSocket 客户端工具 - 用于与后端实时通信
import { ref } from 'vue'

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'

interface WSMessage {
  type: string
  data?: any
  [key: string]: any
}

const isConnected = ref(false)
const messages = ref<WSMessage[]>([])
let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let reconnectAttempts = 0
const MAX_RECONNECT_ATTEMPTS = 5
const RECONNECT_INTERVAL = 3000
const messageHandlers: Map<string, Array<(data: any) => void>> = new Map()

function connect(): void {
  if (ws?.readyState === WebSocket.OPEN || ws?.readyState === WebSocket.CONNECTING) return

  // 带上 JWT token 进行认证
  const token = localStorage.getItem('token')
  const url = token ? `${WS_URL}?token=${token}` : WS_URL

  ws = new WebSocket(url)

  ws.onopen = () => {
    isConnected.value = true
    reconnectAttempts = 0
    console.log('[WS] Connected')
  }

  ws.onmessage = (event: MessageEvent) => {
    try {
      const msg: WSMessage = JSON.parse(event.data)
      messages.value.push(msg)

      // 按消息类型分发
      const handlers = messageHandlers.get(msg.type)
      if (handlers) {
        handlers.forEach((handler) => handler(msg))
      }
    } catch (e) {
      console.error('[WS] Invalid message:', event.data)
    }
  }

  ws.onclose = () => {
    isConnected.value = false
    console.log('[WS] Disconnected')
    attemptReconnect()
  }

  ws.onerror = (error) => {
    console.error('[WS] Error:', error)
    ws?.close()
  }
}

function attemptReconnect(): void {
  if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
    console.log('[WS] Max reconnect attempts reached')
    return
  }
  reconnectAttempts++
  reconnectTimer = setTimeout(() => {
    console.log(`[WS] Reconnect attempt ${reconnectAttempts}`)
    connect()
  }, RECONNECT_INTERVAL)
}

function send(msg: WSMessage): void {
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(msg))
  } else {
    console.warn('[WS] Not connected, cannot send message')
  }
}

function on(type: string, handler: (data: any) => void): void {
  if (!messageHandlers.has(type)) {
    messageHandlers.set(type, [])
  }
  messageHandlers.get(type)!.push(handler)
}

function off(type: string, handler: (data: any) => void): void {
  const handlers = messageHandlers.get(type)
  if (handlers) {
    const idx = handlers.indexOf(handler)
    if (idx > -1) handlers.splice(idx, 1)
  }
}

function disconnect(): void {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  reconnectAttempts = MAX_RECONNECT_ATTEMPTS // 阻止自动重连
  ws?.close()
  ws = null
  isConnected.value = false
}

export function useWebSocket() {
  return { isConnected, messages, connect, send, on, off, disconnect }
}
