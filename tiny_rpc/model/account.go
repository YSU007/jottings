package model

// AccountI ----------------------------------------------------------------------------------------------------
type AccountI interface {
	LoadAccount()
	SaveAccount()
}

type PlayerAccount struct {
	AccountId string
}

func (p PlayerAccount) LoadAccount() {
	panic("implement me")
}

func (p PlayerAccount) SaveAccount() {
	panic("implement me")
}
