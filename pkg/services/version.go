package services

type VersionService struct{}

func NewVersion() *VersionService {
	return &VersionService{}
}

func (s *VersionService) GetVersion() (string, error) {
	return "1.0.0", nil
}
