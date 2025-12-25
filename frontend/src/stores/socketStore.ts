import type { MessageBase } from '@/types/messages';
import { defineStore } from 'pinia';

interface ClientSocketState {
  url: string | null;
  connected: boolean;
  socket: WebSocket | null;
}

export const useSocketStore = defineStore('socketStore', {
  state: (): ClientSocketState => ({
    url: null,
    connected: false,
    socket: null,
  }),
  actions: {
    init(url: string): void {
      this.url = url;
      // do not init if exists
      if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
        return;
      }
      // create connection
      this.socket = new WebSocket(url);
      // onopen callback
      this.socket.onopen = this.onOpen.bind(this);
      // onclose callback
      this.socket.onclose = this.onClose.bind(this);
      // onerror callback
      this.socket.onerror = this.onError.bind(this);
      // onmessage callback
      this.socket.onmessage = this.onMessage.bind(this);
    },
    onOpen(): void {
      this.connected = true;
    },
    onClose(): void {
      this.connected = false;
      setTimeout(() => {
        this.init(this.url!);
      }, 2000);
    },
    onError(): void {
      this.connected = false;
      setTimeout(() => {
        this.init(this.url!);
      }, 2000);
    },
    onMessage(e: MessageEvent<string>): MessageBase | null {
      try {
        const data = JSON.parse(e.data);
        return data;
      } catch {
        console.warn('[WS] Invalid message received', e.data);
        return null;
      }
    },
    sendMessage<T extends MessageBase>(data: T): void {
      if (this.socket && this.connected) {
        this.socket.send(JSON.stringify(data));
      }
    },
  },
});
