package cached

import (
	"context"
	"errors"

	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
)

type UserRepo struct {
	next     outbound.UserRepository
	cache    outbound.UserCache
	reporter outbound.ErrorReporter
}

func NewUserRepo(next outbound.UserRepository, cache outbound.UserCache,
	reporter outbound.ErrorReporter) *UserRepo {
	return &UserRepo{
		next:     next,
		cache:    cache,
		reporter: reporter,
	}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	if err := r.next.Create(ctx, u); err != nil {
		return err
	}
	if err := r.cache.DeleteNotFound(ctx, u.ID); err != nil {
		r.reporter.Report(ctx, err, map[string]string{"op": "cache_delete_not_found", "user_id": string(u.ID)})
	}
	return nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, u *domain.User) error {
	if err := r.next.UpdateProfile(ctx, u); err != nil {
		return err
	}
	if err := r.cache.Delete(ctx, u.ID); err != nil {
		r.reporter.Report(ctx, err, map[string]string{"op": "cache_delete", "user_id": string(u.ID)})
	}
	return nil
}
func (r *UserRepo) UpdateRoles(ctx context.Context, u *domain.User) error {
	if err := r.next.UpdateRoles(ctx, u); err != nil {
		return err
	}
	if err := r.cache.Delete(ctx, u.ID); err != nil {
		r.reporter.Report(ctx, err, map[string]string{"op": "cache_delete", "user_id": string(u.ID)})
	}
	return nil
}
func (r *UserRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	if errors.Is(err, er.ErrUserNotFound) {
		return nil, err
	}
	if !errors.Is(err, er.ErrCacheMiss) {
		r.reporter.Report(ctx, err, map[string]string{"op": "cache_get", "user_id": string(id)})
	}

	u, err = r.next.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, er.ErrUserNotFound) {
			if err := r.cache.SetNotFound(ctx, id); err != nil {
				r.reporter.Report(ctx, err, map[string]string{"op": "cache_set_not_found", "user_id": string(id)})
			}
		}
		return nil, err
	}

	if err := r.cache.Set(ctx, u); err != nil {
		r.reporter.Report(ctx, err, map[string]string{"op": "cache_set", "user_id": string(id)})
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.next.GetByEmail(ctx, email) //не кэшируем по email - операция редкая, кэширование может создать проблемы при логине
}
