document.addEventListener('DOMContentLoaded', () => {
    let myPlayerId = null;

    document.getElementById('joinBtn').addEventListener('click', () => {
        const pName = document.getElementById('playerName').value || 'Guest';
        const rId = document.getElementById('roomId').value || 'room1';
        
        window.socketClient.connect(pName, rId);
        window.uiManager.showGame();
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.isComposing) {
            const lobby = document.getElementById('lobby');
            if (lobby && !lobby.classList.contains('hidden')) {
                document.getElementById('joinBtn').click();
            }
        }
    });

    window.socketClient.on('state_update', (state) => {
        if (!myPlayerId) {
            const pName = document.getElementById('playerName').value || 'Guest';
            const me = state.players.find(p => p.name === pName);
            if (me) myPlayerId = me.id;
        }

        window.uiManager.renderState(state, myPlayerId);
    });

    document.getElementById('startBtn').addEventListener('click', () => {
        window.socketClient.send('start_game', {});
    });

    document.getElementById('playBtn').addEventListener('click', () => {
        const cards = window.uiManager.selectedCards;
        if (cards.length > 0) {
            window.socketClient.send('play_cards', { cards });
            window.uiManager.selectedCards = []; // reset select
        }
    });

    document.getElementById('passBtn').addEventListener('click', () => {
        window.socketClient.send('pass_turn', {});
        window.uiManager.selectedCards = [];
    });

    document.getElementById('leaveBtn').addEventListener('click', () => {
        window.location.reload();
    });

    document.getElementById('closeModalBtn').addEventListener('click', () => {
        document.getElementById('winModal').classList.add('hidden');
    });
});
