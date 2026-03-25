class DaifugoClient {
    constructor() {
        this.ws = null;
        this.handlers = {};
    }

    connect(playerName, roomId) {
        this.ws = new WebSocket(`ws://${location.hostname}:5001/ws`);
        
        this.ws.onopen = () => {
            console.log('Connected to server!');
            this.send('join_game', { player_name: playerName, room_id: roomId });
        };

        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (this.handlers[data.type]) {
                this.handlers[data.type](data.data);
            }
        };

        this.ws.onclose = () => {
            console.log('Disconnected');
        };
    }

    send(type, data) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({ type, data }));
        }
    }

    on(type, handler) {
        this.handlers[type] = handler;
    }
}

window.socketClient = new DaifugoClient();
