package model

// AccountI ----------------------------------------------------------------------------------------------------
type AccountI interface {
	ID() string
	Load()
	Save()
}

type PlayerAccount struct {
	AccountId string
}

func (p PlayerAccount) ID() string {
	return p.AccountId
}

func (p PlayerAccount) Load() {
	panic("implement me")
}

func (p PlayerAccount) Save() {
	panic("implement me")
}
