package storage

type Storage struct {
	PodStorage *PodStorage
}

func New() *Storage {
	return &Storage{
		PodStorage: &PodStorage{},
	}
}

// Return pods storage
func (s *Storage) Pods() *PodStorage {
	return s.PodStorage
}
