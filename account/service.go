package account

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	tracer := otel.Tracer("account-service")
	ctx, span := tracer.Start(ctx, "PostAccount")
	defer span.End()

	a := Account{
		Name: name,
		ID:   ksuid.New().String(),
	}
	if err := s.repo.PutAccount(ctx, a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	tracer := otel.Tracer("account-service")
	ctx, span := tracer.Start(ctx, "GetAccount")
	defer span.End()

	return s.repo.GetAccount(ctx, id)
}

func (s *accountService) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	tracer := otel.Tracer("account-service")
	ctx, span := tracer.Start(ctx, "ListAccounts")
	defer span.End()

	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repo.ListAccounts(ctx, skip, take)
}
