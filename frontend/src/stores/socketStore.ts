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
            // do not init if exists
            if (this.socket && (this.socket.readyState === 0 || this.socket.readyState === 1)) {
                return;
            }
            // create connection
            this.socket = new WebSocket(url);
            // onopen callback
            this.socket.onopen = this.onOpen;
            // onclose callback
            this.socket.onclose = this.onClose;
            // onerror callback
            this.socket.onerror = this.onError;
            // onmessage callback
            this.socket.onmessage = this.onMessage;
        },
        onOpen(): void {
            this.connected = true;
        },
        onClose(e: CloseEvent): void {
            console.warn(e);
            this.connected = true;
            setTimeout(() => {
                this.init(this.url!);
            }, 5000);
        },
        onError(e: Event): void {
            console.warn(e);
            this.connected = true;
            setTimeout(() => {
                this.init(this.url!);
            }, 5000);
        },
        onMessage(e: MessageEvent<string>): unknown {
            try {
                const data = JSON.parse(e.data);
                return data
            } catch {
                console.warn('[WS] Invalid message received', e.data);
            }
        },
    },
});
