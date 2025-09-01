let socket = null;

const connect = () => {
    if (socket !== null && (socket.readyState === 0 || socket.readyState === 1)) {
        return;
    }
    const client_id = localStorage.getItem('player_id')
    console.warn(client_id);
    socket = new WebSocket('ws://localhost:8081/ws');
    socket.onopen = () => {
        socket.send(JSON.stringify({
            type: 'client_id',
            id: client_id,
        }))
    }
    socket.onmessage = (e) => {
        messageHandler(e.data);
    }
    socket.onclose = (e) => {
        console.error(e);
        setTimeout(() => {
            connect();
        }, 5000);
    }
    socket.onerror = (e) => {
        console.error(e);
        setTimeout(() => {
            connect();
        }, 5000);
    }
}

function messageHandler(message) {
    if (message.type === 'new_player_id') {
        localStorage.setItem('player_id', message.id);
        console.log('Received new ID: ' + message.id);
    }
}
