class UIManager {
    constructor() {
        this.lobby = document.getElementById('lobby');
        this.game = document.getElementById('game');
        this.myCards = document.getElementById('myCards');
        this.boardCards = document.getElementById('boardCards');
        this.gameStatus = document.getElementById('gameStatus');
        this.selectedCards = [];
    }

    showGame() {
        this.lobby.classList.add('hidden');
        this.game.classList.remove('hidden');
    }

    renderState(state, myPlayerId) {
        const me = state.players.find(p => p.id === myPlayerId);
        if (me) {
            this.renderHand(me.hand);
        }

        this.renderPlayers(state, myPlayerId);

        if (state.game && state.game.board && state.game.board.length > 0) {
            const lastPlay = state.game.board[state.game.board.length - 1];
            this.renderBoard(lastPlay.cards);
        } else {
            this.boardCards.innerHTML = '';
        }
        
        if (state.state === 'waiting') {
            this.gameStatus.innerText = `Waiting... (${state.players.length}/4)`;
            if (me && me.is_host && state.players.length >= 2) {
                document.getElementById('startBtn').classList.remove('hidden');
            } else {
                document.getElementById('startBtn').classList.add('hidden');
            }
        } else if (state.state === 'playing') {
            document.getElementById('startBtn').classList.add('hidden');
            let isMyTurn = state.game && state.game.current_turn === myPlayerId;
            this.gameStatus.innerText = isMyTurn ? 'YOUR TURN' : 'Waiting for turn';
            
            document.getElementById('playBtn').disabled = !isMyTurn;
            document.getElementById('passBtn').disabled = !isMyTurn;

            if (isMyTurn) {
                this.myCards.classList.add('my-turn');
            } else {
                this.myCards.classList.remove('my-turn');
            }
        }
    }

    renderPlayers(state, myPlayerId) {
        const oppArea = document.getElementById('opponentsArea');
        oppArea.innerHTML = '';
        
        let order = [];
        if (state.state === 'playing' && state.game && state.game.player_order) {
            order = state.game.player_order.map(id => state.players.find(p => p.id === id));
        } else {
            order = state.players;
        }

        order.forEach(p => {
            if (!p) return;
            const div = document.createElement('div');
            div.className = 'player-info glass-panel';
            if (state.state === 'playing' && state.game.current_turn === p.id) {
                div.classList.add('active-turn');
            }
            let text = p.id === myPlayerId ? `${p.name} (You)` : p.name;
            text += ` - Cards: ${p.card_count}`;
            if (p.is_host) text += ' 👑';
            div.innerText = text;
            oppArea.appendChild(div);
        });
    }

    renderHand(hand) {
        this.myCards.innerHTML = '';
        if (!hand) return;
        
        hand.forEach((card) => {
            const el = document.createElement('div');
            el.className = 'card';
            el.style.backgroundImage = `url('/public/cards/${card.suit}_${card.rank}.png')`;
            el.onclick = () => this.toggleSelection(el, card);
            this.myCards.appendChild(el);
        });
    }

    renderBoard(cards) {
        this.boardCards.innerHTML = '';
        cards.forEach((card, i) => {
            const el = document.createElement('div');
            el.className = 'card';
            el.style.backgroundImage = `url('/public/cards/${card.suit}_${card.rank}.png')`;
            el.style.transform = `rotate(${(i - (cards.length-1)/2)*8}deg)`;
            this.boardCards.appendChild(el);
        });
    }

    toggleSelection(el, card) {
        el.classList.toggle('selected');
        if (el.classList.contains('selected')) {
            this.selectedCards.push(card);
        } else {
            this.selectedCards = this.selectedCards.filter(c => c.suit !== card.suit || c.rank !== card.rank);
        }
    }
}

window.uiManager = new UIManager();
