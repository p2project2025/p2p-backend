package admin

import "p2p/repo/admin"

type DashboardServiceInterface interface {
	GetCounts() (map[string]int64, error)
}
type DashboardService struct{}

func (s *DashboardService) GetCounts() (map[string]int64, error) {
	repo := admin.DashboardRepository(&admin.DashboardRepo{})

	return repo.GetCounts()
}
