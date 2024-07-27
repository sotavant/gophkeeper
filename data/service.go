package data

type Service struct {
	dataRepo DataRepository
}

type DataRepository interface {
}

func NewService(u DataRepository) Service {
	return Service{
		dataRepo: u,
	}
}
