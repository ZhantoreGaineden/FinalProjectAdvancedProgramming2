package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/entity"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id string) error
}

type UserUsecase struct {
	repo UserRepository
	nats *nats.Conn
}

func NewUserUsecase(repo UserRepository, natsConn *nats.Conn) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		nats: natsConn,
	}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, fullName, email, password string) (entity.User, error) {
	if fullName == "" {
		return entity.User{}, errors.New("full name is required")
	}
	if email == "" {
		return entity.User{}, errors.New("email is required")
	}
	if password == "" {
		return entity.User{}, errors.New("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	created, err := u.repo.Create(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	u.publishUserRegistered(created)

	return created, nil
}

func (u *UserUsecase) LoginUser(ctx context.Context, email, password string) (string, entity.User, error) {
	if email == "" {
		return "", entity.User{}, errors.New("email is required")
	}
	if password == "" {
		return "", entity.User{}, errors.New("password is required")
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", entity.User{}, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", entity.User{}, errors.New("invalid email or password")
	}

	token := fmt.Sprintf("demo-token-%s", user.ID)

	return token, user, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New("user id is required")
	}

	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if user.ID == "" {
		return entity.User{}, errors.New("user id is required")
	}
	if user.FullName == "" {
		return entity.User{}, errors.New("full name is required")
	}
	if user.Email == "" {
		return entity.User{}, errors.New("email is required")
	}

	return u.repo.Update(ctx, user)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user id is required")
	}

	return u.repo.Delete(ctx, id)
}

func (u *UserUsecase) publishUserRegistered(user entity.User) {
	if u.nats == nil {
		return
	}

	event := map[string]string{
		"user_id":   user.ID,
		"full_name": user.FullName,
		"email":     user.Email,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	_ = u.nats.Publish("user.registered", data)
}
