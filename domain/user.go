package domain

type User struct {
	ID       uint64 `json:"-"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
