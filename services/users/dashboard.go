package users

import (
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
		// Handle empty DB gracefully
		log.Println("Warning: could not fetch admin config, using defaults:", err)
		cnf = &models.AdminConfigData{} // empty defaults
	}

	// Convert USD rate safely
	usdRate := 0.00 // default if cnf.USDTRate is empty or invalid
	if cnf.USDTRate != "" {
		if f, err := strconv.ParseFloat(cnf.USDTRate, 64); err == nil {
			usdRate = f
		} else {
			log.Println("Error converting USDTRate to float64:", err)
		}
	}

	res.SellPrice = cnf.USDTRate
	res.INRBalance = res.Balance * usdRate
	res.WalletAddress = cnf.SecureWalletAddress

	return res, nil
}
