package deposit

import (
	"p2p/models"
	"p2p/repo/deposit"
)

// DepositServiceInterface defines the methods for deposit operations
type DepositServiceInterface interface {
	UpdateDepositStatus(depositID string, approve bool) error
	CreateDeposit(req models.DepositRequest) error
	ListDeposits() ([]models.DepositRes, error)
	GetDepositByID(id string) (*models.DepositRes, error)
	SearchDepositsByUsername(username string) ([]models.DepositRes, error)
	GetDepositsByUserID(userID string) ([]models.DepositRes, error)
}

// DepositService implements DepositServiceInterface
type DepositService struct {
}

// Create new deposit request
func (s *DepositService) CreateDeposit(req models.DepositRequest) error {
	repo := deposit.DepositRepository(&deposit.DepositRepo{})
	return repo.DepositRequest(req)
}

func (s *DepositService) UpdateDepositStatus(depositID string, approve bool) error {
	repo := deposit.DepositRepository(&deposit.DepositRepo{})
	return repo.UpdateDepositStatus(depositID, approve)
}

// GetDepositsByUserID - paginated deposits for a given user
func (s *DepositService) GetDepositsByUserID(userID string) ([]models.DepositRes, error) {
	repo := deposit.DepositRepository(&deposit.DepositRepo{})
	return repo.GetAllByUserID(userID)
}

// List deposits with pagination
func (s *DepositService) ListDeposits() ([]models.DepositRes, error) {
	repo := deposit.DepositRepository(&deposit.DepositRepo{})
	return repo.GetAll()
}

// Get deposit by ID
func (s *DepositService) GetDepositByID(id string) (*models.DepositRes, error) {
	repo := deposit.DepositRepository(&deposit.DepositRepo{})
	return repo.GetByID(id)
}

// Search deposits by username
func (s *DepositService) SearchDepositsByUsername(username string) ([]models.DepositRes, error) {
	repo := deposit.DepositRepository(&deposit.DepositRepo{})
	return repo.SearchByUsername(username)
}
