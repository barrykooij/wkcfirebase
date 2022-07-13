package wkcfirebase

type StateDocument struct {
	Round    int    `firestore:"round"`
	VoteOpen bool   `firestore:"vote_open"`
	Question string `firestore:"question"`
}

type StateChangedListener func(*StateDocument)
