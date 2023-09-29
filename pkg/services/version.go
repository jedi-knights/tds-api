package services

import "github.com/jedi-knights/tds-api/pkg"

type VersionService struct{}

func NewVersion() *VersionService {
	return &VersionService{}
}

func (s *VersionService) GetVersion() (string, error) {
	return pkg.Version, nil
}
