package domain

// Data структура для хранения данных в памяти
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

// DataName структура для хранения данных в памяти в кратком виде
type DataName struct {
	Name string
	ID   uint64
}
