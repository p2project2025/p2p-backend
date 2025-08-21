package users

import (
	"log"
	"p2p/models"
	"p2p/repo/admin"
	"p2p/repo/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardServiceInterface interface {
	GetUserDashboard(userId primitive.ObjectID) (*models.UserDash, error)
}
type DashboardService struct{}

func (s *DashboardService) GetUserDashboard(userId primitive.ObjectID) (*models.UserDash, error) {
	repo := users.DashboardRepository(&users.DashboardRepo{})
	res, err := repo.GetUserDashboard(userId)
	if err != nil {
		log.Println("Error fetching user dashboard:", err)
		return nil, err
	}

	adminRepo := admin.AdminRepository(&admin.AdminRepo{})
	cnf, err := adminRepo.Fetch()
	if err != nil {
		log.Println("Error fetching admin config:", err)
		return nil, err
	}
	res.SellPrice = cnf.USDTRate
	res.WalletAddress = cnf.SecureWalletAddress

	return res, nil
}
