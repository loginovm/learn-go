package storage

type Repo interface {
	Close() error
	EventRepo
}
