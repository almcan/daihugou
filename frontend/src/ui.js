class UIManager {
    constructor() {
        this.lobby = document.getElementById('lobby');
        this.game = document.getElementById('game');
        this.myCards = document.getElementById('myCards');
        this.boardCards = document.getElementById('boardCards');
        this.gameStatus = document.getElementById('gameStatus');
        this.selectedCards = [];
        this.finishedPlayers = new Set();
    }

    showGame() {
        this.lobby.classList.add('hidden');
        this.game.classList.remove('hidden');
    }

    renderState(state, myPlayerId) {
        this.currentState = state;
        this.myPlayerId = myPlayerId;
        
        const me = state.players.find(p => p.id === myPlayerId);
        if (me) {
            this.renderHand(me.hand);
        }

        this.renderPlayers(state, myPlayerId);

        if (state.game && state.game.board && state.game.board.length > 0) {
            this.renderBoard(state.game.board);
        } else {
            this.boardCards.innerHTML = '';
        }
        
        if (state.state === 'waiting') {
            this.finishedPlayers.clear();
            this.gameStatus.innerText = `Waiting... (${state.players.length}/6)`;
            document.getElementById('leaveBtn').classList.remove('hidden');
            if (me && me.is_host && state.players.length >= 2) {
                document.getElementById('startBtn').classList.remove('hidden');
            } else {
                document.getElementById('startBtn').classList.add('hidden');
            }
        } else if (state.state === 'playing') {
            document.getElementById('startBtn').classList.add('hidden');
            document.getElementById('leaveBtn').classList.add('hidden');
            
            if (state.game && state.game.interaction) {
                const ix = state.game.interaction;
                if (ix.player_id === this.myPlayerId) {
                    this.gameStatus.innerText = 'ACTION REQUIRED';
                    this.showInteractionModal(ix, state);
                } else {
                    this.gameStatus.innerText = 'Waiting for player action...';
                    document.getElementById('interactionModal').classList.add('hidden');
                    document.getElementById('playBtn').disabled = true;
                    document.getElementById('passBtn').disabled = true;
                }
            } else {
                document.getElementById('interactionModal').classList.add('hidden');
                let isMyTurn = state.game && state.game.current_turn === myPlayerId;
                this.gameStatus.innerText = isMyTurn ? 'YOUR TURN' : 'Waiting for turn';
                this.updateHighlights();
            }
        }

        const ruleStatus = document.getElementById('ruleStatus');
        let statuses = [];
        if (state.game && state.game.is_revolution) {
            statuses.push('🔥 革命中');
        }
        if (state.game && state.game.is_11_back) {
            statuses.push('🔄 11バック');
        }
        if (state.game && state.game.bound_suits && state.game.bound_suits.length > 0) {
            statuses.push('🔒 マーク縛り');
        }
        if (state.game && state.game.is_sequence_bound) {
            statuses.push('📈 連番縛り');
        }
        if (statuses.length > 0) {
            ruleStatus.innerText = statuses.join(' / ');
            ruleStatus.classList.remove('hidden');
        } else {
            ruleStatus.classList.add('hidden');
        }

        if (state.state === 'playing') {
            state.players.forEach(p => {
                if (p.card_count === 0 && p.rank > 0 && !this.finishedPlayers.has(p.id)) {
                    this.finishedPlayers.add(p.id);
                    document.getElementById('winMessage').innerText = `${p.name} がアガりました！ (${p.rank}位)`;
                    document.getElementById('winModal').classList.remove('hidden');
                    setTimeout(() => {
                        document.getElementById('winModal').classList.add('hidden');
                    }, 1500);
                }
            });
        }

        if (state.game && state.game.meme_event && state.game.meme_event !== this.lastMemeEvent) {
            this.lastMemeEvent = state.game.meme_event;
            this.showMemePopup(state.game.meme_event);
        }
    }

    showMemePopup(text) {
        if (!text) return;
        const popup = document.createElement('div');
        popup.className = 'meme-popup';
        popup.innerText = text;
        document.body.appendChild(popup);
        setTimeout(() => popup.remove(), 2500);
    }

    showInteractionModal(ix, state) {
        const modal = document.getElementById('interactionModal');
        const content = document.getElementById('interactionContent');
        modal.classList.remove('hidden');
        content.innerHTML = '';
        
        const myHand = state.players.find(p => p.id === this.myPlayerId).hand;
        const reqCount = ix.count || 1;
        let selectedItems = [];

        const effRev = state.game.is_revolution !== state.game.is_11_back;
        const getStr = (c) => {
            if (c.suit === 'joker') return 99;
            let s = c.rank;
            if (s === 1) s = 14;
            if (s === 2) s = 15;
            return effRev ? 18 - s : s;
        };
        const getRankStr = (r) => {
            let s = r;
            if (s === 1) s = 14;
            if (s === 2) s = 15;
            return effRev ? 18 - s : s;
        };

        const renderPicker = (items, isCard) => {
            const container = document.createElement('div');
            container.className = 'picker-container';
            items.forEach((item) => {
                const el = document.createElement('div');
                el.className = isCard ? 'card picker-card' : 'picker-item';
                if (isCard) {
                    let bg = `/public/cards/${item.suit}_${item.rank}.png`;
                    if (item.suit === 'joker') bg = `/public/cards/joker.png`;
                    el.style.backgroundImage = `url('${bg}')`;
                } else {
                    el.innerText = item;
                }
                
                el.onclick = () => {
                    const idx = selectedItems.findIndex(x => x === item);
                    if (idx > -1) {
                        selectedItems.splice(idx, 1);
                        el.classList.remove('selected');
                    } else if (selectedItems.length < reqCount) {
                        selectedItems.push(item);
                        el.classList.add('selected');
                    }
                    btn.disabled = selectedItems.length !== reqCount;
                };
                container.appendChild(el);
            });
            return container;
        };

        const btn = document.createElement('button');
        btn.className = 'primary-btn mt-4';
        btn.innerText = '決定';
        btn.disabled = true;

        if (ix.type === '7_pass' || ix.type === '10_discard') {
            content.innerHTML = `<h3>${ix.type === '7_pass' ? '7渡し' : '10捨て'}</h3>
                <p>${ix.type === '7_pass' ? '次の人に渡す' : '捨てる'}カードを <strong>${reqCount}</strong> 枚選んでください</p>`;
            content.appendChild(renderPicker(myHand, true));
            btn.onclick = () => {
                window.socketClient.send('resolve_interaction', { cards: selectedItems });
            };
            content.appendChild(btn);
        } else if (ix.type === '9_pick') {
            content.innerHTML = `<h3>9拾い（栗拾い）</h3>
                <p>墓地（これまでに捨てられたカード）から <strong>${reqCount}</strong> 枚選んで手札に加えます</p>`;
            let cem = [...(state.game.cemetery || [])];
            cem.sort((a,b) => getStr(b) - getStr(a));
            
            content.appendChild(renderPicker(cem, true));
            btn.onclick = () => {
                window.socketClient.send('resolve_interaction', { cards: selectedItems });
            };
            content.appendChild(btn);
        } else if (ix.type === '12_bomber') {
            content.innerHTML = `<h3>12ボンバー</h3>
                <p>全員の手札から消し去る数字を <strong>${reqCount}</strong> つ選んでください</p>`;
            let ranks = [1,2,3,4,5,6,7,8,9,10,11,12,13];
            ranks.sort((a,b) => getRankStr(b) - getRankStr(a));
            
            content.appendChild(renderPicker(ranks, false));
            btn.onclick = () => {
                window.socketClient.send('resolve_interaction', { ranks: selectedItems });
            };
            content.appendChild(btn);
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
            let bg = `/public/cards/${card.suit}_${card.rank}.png`;
            if (card.suit === 'joker') bg = `/public/cards/joker.png`;
            el.style.backgroundImage = `url('${bg}')`;
            
            // Re-apply selection state if it matches a previously selected card
            if (this.selectedCards.some(c => c.suit === card.suit && c.rank === card.rank)) {
                el.classList.add('selected');
            }

            el.onclick = () => this.toggleSelection(el, card);
            this.myCards.appendChild(el);
        });
        
        if (this.currentState && this.currentState.state === 'playing') {
            this.updateHighlights();
        }
    }

    updateHighlights() {
        if (!this.currentState || this.currentState.state !== 'playing') return;
        const state = this.currentState;
        const isMyTurn = state.game.current_turn === this.myPlayerId;
        const hand = state.players.find(p => p.id === this.myPlayerId)?.hand;
        if (!hand) return;

        const pBtn = document.getElementById('playBtn');
        const passBtn = document.getElementById('passBtn');
        pBtn.disabled = !isMyTurn || this.selectedCards.length === 0;
        passBtn.disabled = !isMyTurn;

        let lastPlay = null;
        if (state.game.board && state.game.board.length > 0) {
            lastPlay = state.game.board[state.game.board.length - 1];
        }

        const effRev = state.game.is_revolution !== state.game.is_11_back;
        const getStr = (c) => {
            if (c.suit === 'joker') return 99;
            let s = c.rank;
            if (s === 1) s = 14;
            if (s === 2) s = 15;
            return effRev ? 18 - s : s;
        };

        const bStr = lastPlay ? getStr(lastPlay.cards[0]) : -1;
        const reqCount = lastPlay ? lastPlay.cards.length : 1;
        const seqBound = state.game.is_sequence_bound;
        
        let rankCounts = {};
        hand.forEach(c => {
            if (c.suit !== 'joker') rankCounts[c.rank] = (rankCounts[c.rank] || 0) + 1;
        });
        const jp = hand.filter(c => c.suit === 'joker').length;

        let valids = new Set();
        hand.forEach(c => {
            if (c.suit !== 'joker' && getStr(c) > bStr) {
                if (!lastPlay || rankCounts[c.rank] + jp >= reqCount) {
                    if (!seqBound || getStr(c) === bStr + 1) {
                        valids.add(c.rank);
                    }
                }
            }
        });

        const sRanks = [...new Set(this.selectedCards.filter(c => c.suit !== 'joker').map(c => c.rank))];
        const hasJk = this.selectedCards.some(c => c.suit === 'joker');

        Array.from(this.myCards.children).forEach((el, i) => {
            const c = hand[i];
            const sel = el.classList.contains('selected');
            let bright = false;

            if (isMyTurn) {
                if (this.selectedCards.length > 0) {
                    if (sel) bright = true;
                    else if (c.suit === 'joker') bright = true;
                    else if (sRanks.length === 1 && c.rank === sRanks[0]) bright = true;
                    else if (sRanks.length === 0 && hasJk) bright = true;
                } else {
                    if (!lastPlay) bright = true;
                    else if (c.suit === 'joker' || valids.has(c.rank)) bright = true;
                }
                if (state.game.bound_suits && state.game.bound_suits.length > 0) {
                    if (c.suit !== 'joker' && !state.game.bound_suits.includes(c.suit)) bright = false;
                }
            }
            
            el.style.filter = bright ? 'brightness(1.1)' : 'brightness(0.35)';
        });
        
        // Final sanity check for play button if cards are selected
        if (this.selectedCards.length > 0 && lastPlay) {
            if (this.selectedCards.length !== lastPlay.cards.length) pBtn.disabled = true;
        }
    }

    renderBoard(board) {
        this.boardCards.innerHTML = '';
        board.forEach((play, playIdx) => {
            const playContainer = document.createElement('div');
            playContainer.className = 'play-container';
            
            // Random-ish rotation so that multiple plays look like a messy pile
            const groupRotation = Math.sin(playIdx * 12.34) * 20; // -20 to +20 degrees
            const offsetX = Math.cos(playIdx * 43.21) * 20; 
            const offsetY = Math.sin(playIdx * 43.21) * 10;
            
            playContainer.style.transform = `translate(${offsetX}px, ${offsetY}px)`;

            play.cards.forEach((card, i) => {
                const el = document.createElement('div');
                el.className = 'card';
                let bg = `/public/cards/${card.suit}_${card.rank}.png`;
                if (card.suit === 'joker') bg = `/public/cards/joker.png`;
                el.style.backgroundImage = `url('${bg}')`;
                
                const spreadRot = (i - (play.cards.length-1)/2)*8;
                el.style.transform = `rotate(${groupRotation + spreadRot}deg)`;
                playContainer.appendChild(el);
            });
            this.boardCards.appendChild(playContainer);
        });
    }

    toggleSelection(el, card) {
        if (!this.currentState || this.currentState.game.current_turn !== this.myPlayerId) return;
        
        el.classList.toggle('selected');
        if (el.classList.contains('selected')) {
            this.selectedCards.push(card);
        } else {
            this.selectedCards = this.selectedCards.filter(c => c.suit !== card.suit || c.rank !== card.rank);
        }
        this.updateHighlights();
    }
}

window.uiManager = new UIManager();
