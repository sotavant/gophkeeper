package domain

type Data struct {
	Name string
	ID,
	Version,
	UID uint64
	Pass,
	CardNum,
	Text,
	Meta,
	Login *string
	FileID *uint64
}
