package users

import (
	"fmt"
	"log"
	"p2p/models"
	"p2p/repo/admin"
	"p2p/repo/users"
	"strconv"

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
	usdRate, err := strconv.ParseFloat(cnf.USDTRate, 64)
	if err != nil {
		fmt.Println("Error converting string to float64:", err)
		return res, fmt.Errorf("error converting string to float64: %s", err)
	}
	res.SellPrice = cnf.USDTRate
	res.INRBalance = res.Balance * usdRate
	res.WalletAddress = cnf.SecureWalletAddress

	return res, nil
}
