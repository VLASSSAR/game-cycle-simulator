package domain

type State string

const (
	Init       State = "init"
	Reserve    State = "reserve"
	Engine     State = "engine"
	Settlement State = "settlement"
	Wallet     State = "wallet"
	ErrorState State = "error"
)

func AllStates() []State {
	return []State{
		Init,
		Reserve,
		Engine,
		Settlement,
		Wallet,
		ErrorState,
	}
}

func IsTerminalState(state State) bool {
	return state == Wallet || state == ErrorState
}
