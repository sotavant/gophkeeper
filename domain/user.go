package domain

type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
