let socket = null;

const connect = () => {
    if (socket !== null && (socket.readyState === 0 || socket.readyState === 1)) {
        return;
    }
    socket = new WebSocket("ws://localhost:8081/ws");
    socket.onmessage = (e) => {
        console.log(e.data);
    }
    socket.onclose = (e) => {
        console.error(e);
        setTimeout(() => {
            connect()
        }, 5000);
    }
    socket.onerror = (e) => {
        console.error(e);
        setTimeout(() => {
            connect()
        }, 5000);
    }
}
