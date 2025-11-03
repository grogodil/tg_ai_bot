package bot

type UserState struct {
    QuizIndex   int
    QuizScore   int
    WaitingForAI bool
}

var UserStates = make(map[int64]*UserState)

func GetUserState(userID int64) *UserState {
    state, exists := UserStates[userID]
    if !exists {
        state = &UserState{}
        UserStates[userID] = state
    }
    return state
}