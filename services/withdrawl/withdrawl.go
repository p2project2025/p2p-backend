package withdrawl

import (
	"p2p/models"
	"p2p/repo/withdrawl"
)

type WithdrawlServiceInterface interface {
	CreateWithdrawl(req models.WithdrawlRequest) error
	UpdateWithdrawStatus(withdrawID string, approve bool) error
	ListWithdrawls() ([]models.WithdrawlRes, error)
	GetWithdrawlByID(id string) (*models.WithdrawlRes, error)
	SearchWithdrawlsByUsername(username string) ([]models.WithdrawlRes, error)
	GetWithdrawlsByUserID(userID string) ([]models.WithdrawlRes, error)
}
type WithdrawlService struct {
}

// Create new withdrawl request
func (s *WithdrawlService) CreateWithdrawl(req models.WithdrawlRequest) error {
	repo := withdrawl.WithdrawlRepository(&withdrawl.WithdrawlRepo{})
	return repo.WithdrawlRequest(req)
}

func (s *WithdrawlService) UpdateWithdrawStatus(withdrawID string, approve bool) error {
	repo := withdrawl.WithdrawlRepository(&withdrawl.WithdrawlRepo{})
	return repo.UpdateWithdrawStatus(withdrawID, approve)
}

// GetWithdrawlsByUserID - paginated withdrawls for a given user
func (s *WithdrawlService) GetWithdrawlsByUserID(userID string) ([]models.WithdrawlRes, error) {
	repo := withdrawl.WithdrawlRepository(&withdrawl.WithdrawlRepo{})
	return repo.GetAllByUserID(userID)
}

// List withdrawls with pagination
func (s *WithdrawlService) ListWithdrawls() ([]models.WithdrawlRes, error) {
	repo := withdrawl.WithdrawlRepository(&withdrawl.WithdrawlRepo{})
	return repo.GetAll()
}

// Get withdrawl by ID
func (s *WithdrawlService) GetWithdrawlByID(id string) (*models.WithdrawlRes, error) {
	repo := withdrawl.WithdrawlRepository(&withdrawl.WithdrawlRepo{})
	return repo.GetByID(id)
}

// Search withdrawls by username
func (s *WithdrawlService) SearchWithdrawlsByUsername(username string) ([]models.WithdrawlRes, error) {
	repo := withdrawl.WithdrawlRepository(&withdrawl.WithdrawlRepo{})
	return repo.SearchByUsername(username)
}
