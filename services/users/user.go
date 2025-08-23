package users

import (
	"errors"
	"log"
	"p2p/models"
	"p2p/repo/users"
	"p2p/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceInterface interface {
	RegisterUser(user models.User) (primitive.ObjectID, error)
	SignInUser(user models.Login) (models.User, error)
	BlockUser(userID primitive.ObjectID, block bool) error
	GetAllUsers() ([]models.User, error)
}
type UserService struct{}

func (s *UserService) RegisterUser(user models.User) (primitive.ObjectID, error) {
	user.Role = "user"
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

func (s *UserService) SignInUser(user models.Login) (models.User, error) {

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

	if userData.IsBlocked {
		return userData, errors.New("User Blocked Contact Admin")
	}
	//  password check
	if err := utils.CheckPasswordHash(user.Password, userData.Password); err != nil {
		log.Println("Password Mismatch", err)
		return models.User{}, errors.New("Password Mismatch")
	}

	adminRepo := admin.AdminRepository(&admin.AdminRepo{})
	cnf, err := adminRepo.Fetch()
	if err != nil {
		log.Println("Error fetching admin config:", err)
		return models.User{}, err
	}
	usdRate, err := strconv.ParseFloat(cnf.USDTRate, 64)
	if err != nil {
		log.Println("Error converting string to float64:", err)
		return userData, nil // //cant efffect login
	}

	userData.INRBalance = usdRate * userData.Balance
	return userData, nil

}

func (s *UserService) BlockUser(userID primitive.ObjectID, block bool) error {

	// Repo instance
	repo := users.UserRepository(&users.UserRepo{})

	return repo.BlockUser(userID, block)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	repo := users.UserRepository(&users.UserRepo{})

	return repo.GetAllUsers()
}
