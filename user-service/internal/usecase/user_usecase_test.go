package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/user-service/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	usersByID    map[string]entity.User
	usersByEmail map[string]entity.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByID:    make(map[string]entity.User),
		usersByEmail: make(map[string]entity.User),
	}
}

func (r *fakeUserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	user.ID = "user-1"
	user.CreatedAt = time.Now()

	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user

	return user, nil
}

func (r *fakeUserRepository) GetByID(ctx context.Context, id string) (entity.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeUserRepository) Update(ctx context.Context, user entity.User) (entity.User, error) {
	existing, ok := r.usersByID[user.ID]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}

	existing.FullName = user.FullName
	existing.Email = user.Email

	r.usersByID[user.ID] = existing
	r.usersByEmail[user.Email] = existing

	return existing, nil
}

func (r *fakeUserRepository) Delete(ctx context.Context, id string) error {
	user, ok := r.usersByID[id]
	if ok {
		delete(r.usersByEmail, user.Email)
	}
	delete(r.usersByID, id)
	return nil
}

func TestUserUsecaseRegisterUserSuccess(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepository()
	usecase := NewUserUsecase(repo, nil)

	user, err := usecase.RegisterUser(ctx, "Zhantore Gaineden", "zhantore@example.com", "123456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == "" {
		t.Fatal("expected user id to be generated")
	}

	if user.PasswordHash == "123456" {
		t.Fatal("password must be hashed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("123456")); err != nil {
		t.Fatalf("expected valid bcrypt hash, got %v", err)
	}
}

func TestUserUsecaseRegisterUserRequiresEmail(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepository()
	usecase := NewUserUsecase(repo, nil)

	_, err := usecase.RegisterUser(ctx, "Zhantore Gaineden", "", "123456")
	if err == nil {
		t.Fatal("expected error when email is empty")
	}
}

func TestUserUsecaseLoginUserSuccess(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepository()
	usecase := NewUserUsecase(repo, nil)

	registered, err := usecase.RegisterUser(ctx, "Zhantore Gaineden", "zhantore@example.com", "123456")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	token, loggedUser, err := usecase.LoginUser(ctx, "zhantore@example.com", "123456")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if token == "" {
		t.Fatal("expected token")
	}

	if loggedUser.ID != registered.ID {
		t.Fatalf("expected user id %s, got %s", registered.ID, loggedUser.ID)
	}
}

func TestUserUsecaseLoginUserWrongPassword(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepository()
	usecase := NewUserUsecase(repo, nil)

	_, err := usecase.RegisterUser(ctx, "Zhantore Gaineden", "zhantore@example.com", "123456")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, _, err = usecase.LoginUser(ctx, "zhantore@example.com", "wrong-password")
	if err == nil {
		t.Fatal("expected login error with wrong password")
	}
}
