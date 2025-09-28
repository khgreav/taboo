import { useSocketStore } from '@/stores/socketStore';

export default function registerSocket(): void {
  const clientSocket = useSocketStore();
  clientSocket.init('ws://localhost:8081/ws');
}

