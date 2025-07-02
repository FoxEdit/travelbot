package telegram

import "sync"

type UserState int

const (
	StateNone UserState = iota
	StateAwaitingDepositAmount
	StateAwaitingDepositCurrency
	StateAwaitingWithdrawAmount
	StateAwaitingWithdrawCurrency
	StateAwaitingChooseBaseCurrency
	StateAwaitingAddForeignCurrency
	StateAwaitingRemoveForeignCurrency
	StateAwaitingChangeBaseCurrency
	StateAwaitingHelpRequest
)

type userSession struct {
	State UserState
	Data  map[string]string
}

type StateManager struct {
	mu       sync.Mutex
	sessions map[int64]*userSession
}

func NewStateManager() *StateManager {
	return &StateManager{
		sessions: make(map[int64]*userSession),
	}
}

func (sm *StateManager) SetState(chatID int64, state UserState, data map[string]string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if data == nil {
		data = make(map[string]string)
	}
	sm.sessions[chatID] = &userSession{State: state, Data: data}
}

func (sm *StateManager) GetState(chatID int64) (UserState, map[string]string, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session, ok := sm.sessions[chatID]
	if !ok {
		return StateNone, nil, false
	}
	return session.State, session.Data, true
}

func (sm *StateManager) ClearState(chatID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, chatID)
}
