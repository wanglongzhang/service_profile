package service

import (
	"errors"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/model"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/repository"
)

// UserService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserService interface {
	GetAll() []model.User
	GetByID(id int64) (model.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (model.User, bool)
	DeleteByID(id int64) bool

	Update(id int64, user model.User) (model.User, error)
	UpdatePassword(id int64, newPassword string) (model.User, error)
	UpdateUsername(id int64, newUsername string) (model.User, error)

	Create(userPassword string, user model.User) (model.User, error)
}

// NewUserService returns the default user service.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repository.UserRepository
}

// GetAll returns all users.
func (s *userService) GetAll() []model.User {
	return s.repo.SelectMany(func(_ model.User) bool {
		return true
	}, -1)
}

// GetByID returns a user based on its id.
func (s *userService) GetByID(id int64) (model.User, bool) {
	return s.repo.Select(func(m model.User) bool {
		return m.ID == id
	})
}

// GetByUsernameAndPassword returns a user based on its username and passowrd,
// used for authentication.
func (s *userService) GetByUsernameAndPassword(username, userPassword string) (model.User, bool) {
	if username == "" || userPassword == "" {
		return model.User{}, false
	}

	return s.repo.Select(func(m model.User) bool {
		if m.Username == username {
			hashed := m.HashedPassword
			if ok, _ := model.ValidatePassword(userPassword, hashed); ok {
				return true
			}
		}
		return false
	})
}

// Update updates every field from an existing User,
// it's not safe to be used via public API,
// however we will use it on the web/controllers/user_controller.go#PutBy
// in order to show you how it works.
func (s *userService) Update(id int64, user model.User) (model.User, error) {
	user.ID = id
	return s.repo.InsertOrUpdate(user)
}

// UpdatePassword updates a user's password.
func (s *userService) UpdatePassword(id int64, newPassword string) (model.User, error) {
	// update the user and return it.
	hashed, err := model.GeneratePassword(newPassword)
	if err != nil {
		return model.User{}, err
	}

	return s.Update(id, model.User{
		HashedPassword: hashed,
	})
}

// UpdateUsername updates a user's username.
func (s *userService) UpdateUsername(id int64, newUsername string) (model.User, error) {
	return s.Update(id, model.User{
		Username: newUsername,
	})
}

// Create inserts a new User,
// the userPassword is the client-typed password
// it will be hashed before the insertion to our repository.
func (s *userService) Create(userPassword string, user model.User) (model.User, error) {
	if user.ID > 0 || userPassword == "" || user.Firstname == "" || user.Username == "" {
		return model.User{}, errors.New("unable to create this user")
	}

	hashed, err := model.GeneratePassword(userPassword)
	if err != nil {
		return model.User{}, err
	}
	user.HashedPassword = hashed

	return s.repo.InsertOrUpdate(user)
}

// DeleteByID deletes a user by its id.
//
// Returns true if deleted otherwise false.
func (s *userService) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m model.User) bool {
		return m.ID == id
	}, 1)
}
