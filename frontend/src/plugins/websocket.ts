import { useSocketStore } from '@/stores/socketStore';

export default function registerSocket(): void {
  const clientSocket = useSocketStore();
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const development = import.meta.env.DEV;
  if (development) {
    clientSocket.init('ws://localhost:8080/ws');
    return;
  }
  clientSocket.init(`${protocol}://${window.location.host}/ws`);
}

