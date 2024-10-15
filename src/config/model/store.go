package model

type Store struct {
	Snapshot string `json:"snapshot"`
}

func (store Store) Init() {
	if '/' == store.Snapshot[len(store.Snapshot)-1] {
		store.Snapshot = store.Snapshot[:len(store.Snapshot)-1]
	}
}
