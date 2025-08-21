package admin

import (
	"errors"
	"log"
	"p2p/models"
	"p2p/repo/admin"
	"p2p/repo/users"
	"p2p/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminServiceInterface interface {
	RegisterAdmin(user models.User) (primitive.ObjectID, error)
	SignInAdmin(user models.Login) (models.User, error)
	FetchAdminConfig() (*models.AdminConfigData, error)
	UpsertAdminConfig(adminConfig models.AdminConfigData) (primitive.ObjectID, error)
}
type AdminService struct{}

func (s *AdminService) RegisterAdmin(user models.User) (primitive.ObjectID, error) {
	user.Role = "admin"

	user.IsBlocked = false
	user.Balance = 0.00
	user.CreatedAt = time.Now()

	if user.Email == "" || user.Password == "" {
		return primitive.NilObjectID, errors.New("email or password can't be empty")
	}
	if user.Name == "" {
		user.Name = user.Email // Default to email if name is not provided
	}
	if user.PhoneNum == "" {
		user.PhoneNum = "0000000000" // Default phone number if not provided
	}

	// Repo instance
	repo := users.UserRepository(&users.UserRepo{})

	// Duplicate email check
	if exists, err := repo.CheckEmailExists(user.Email); err != nil {
		log.Println("Error checking email:", err)
		return primitive.NilObjectID, err
	} else if exists {
		return primitive.NilObjectID, errors.New("email already registered")
	}

	// Duplicate phone check
	if exists, err := repo.CheckPhoneExists(user.PhoneNum); err != nil {
		log.Println("Error checking phone number:", err)
		return primitive.NilObjectID, err
	} else if exists {
		return primitive.NilObjectID, errors.New("phone number already registered")
	}

	// Hashing password
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return primitive.NilObjectID, err
	}
	user.Password = hashedPwd

	return repo.RegisterUser(user)
}

func (s *AdminService) SignInAdmin(user models.Login) (models.User, error) {

	if user.Email == "" || user.Password == "" {
		return models.User{}, errors.New("email or password can't be empty")
	}

	// Repo instance
	repo := users.UserRepository(&users.UserRepo{})

	userData, err := repo.GetUserByEmail(user.Email)
	if err != nil {
		log.Println("Error fetching user by email:", err)
		return models.User{}, err
	}
	//  password check
	if err := utils.CheckPasswordHash(user.Password, userData.Password); err != nil {
		log.Println("Password Mismatch", err)
		return models.User{}, errors.New("Password Mismatch")
	}

	return userData, nil
}

func (s *AdminService) FetchAdminConfig() (*models.AdminConfigData, error) {
	repo := admin.AdminRepository(&admin.AdminRepo{})
	config, err := repo.Fetch()
	if err != nil {
		log.Println("Error fetching admin config:", err)
		return nil, err
	}
	return config, nil
}

func (s *AdminService) UpsertAdminConfig(adminConfig models.AdminConfigData) (primitive.ObjectID, error) {
	repo := admin.AdminRepository(&admin.AdminRepo{})
	id, err := repo.Upsert(adminConfig)
	if err != nil {
		log.Println("Error updating admin config:", err)
		return primitive.NilObjectID, err
	}
	return id, nil
}
