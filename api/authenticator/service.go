package authenticator

import (
	"api/models"
	"context"
	"github.com/supabase-community/supabase-go"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email string) (*models.User, error)
	HandleLogin(ctx context.Context, email string) error
	UpdateFailedLogin(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
}

type SupabaseAuthService struct {
	client *supabase.Client
	AuthService
}

func NewSupabaseAuthService(client *supabase.Client) AuthService {
	return &SupabaseAuthService{client: client}
}
