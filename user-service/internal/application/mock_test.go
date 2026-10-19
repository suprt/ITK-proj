package application_test

import (
	"context"
	"time"

	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

// --- UserRepository ---

type mockUserRepo struct {
	users  map[domain.UserID]*domain.User
	emails map[string]domain.UserID
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:  make(map[domain.UserID]*domain.User),
		emails: make(map[string]domain.UserID),
	}
}

func (m *mockUserRepo) Create(_ context.Context, u *domain.User) error {
	if _, exists := m.emails[u.Email]; exists {
		return er.ErrEmailTaken
	}
	for _, existing := range m.users {
		if existing.Username == u.Username {
			return er.ErrUsernameTaken
		}
	}
	m.users[u.ID] = u
	m.emails[u.Email] = u.ID
	return nil
}

func (m *mockUserRepo) Update(_ context.Context, u *domain.User) error {
	if _, exists := m.users[u.ID]; !exists {
		return er.ErrUserNotFound
	}
	if exiID, taken := m.emails[u.Email]; taken && exiID != u.ID {
		return er.ErrEmailTaken
	}
	m.users[u.ID] = u
	m.emails[u.Email] = u.ID
	return nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id domain.UserID) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, er.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	id, exists := m.emails[email]
	if !exists {
		return nil, er.ErrUserNotFound
	}
	return m.users[id], nil
}

// --- CredentialsRepository ---

type mockCredentialsRepo struct {
	data map[domain.UserID]*domain.Credentials
}

func newMockCredentialsRepo() *mockCredentialsRepo {
	return &mockCredentialsRepo{data: make(map[domain.UserID]*domain.Credentials)}
}

func (m *mockCredentialsRepo) Save(_ context.Context, c *domain.Credentials) error {
	m.data[c.UserID] = c
	return nil
}

func (m *mockCredentialsRepo) GetByID(_ context.Context, userID domain.UserID) (*domain.Credentials, error) {
	c, exists := m.data[userID]
	if !exists {
		return nil, er.ErrUserNotFound
	}
	return c, nil
}

// --- SessionRepository ---

type mockSessionRepo struct {
	data        map[domain.RefreshToken]*domain.RefreshSession
	deleteCalls int
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{data: make(map[domain.RefreshToken]*domain.RefreshSession)}
}

func (m *mockSessionRepo) Save(_ context.Context, s *domain.RefreshSession) error {
	m.data[s.Token] = s
	return nil
}

func (m *mockSessionRepo) GetByToken(_ context.Context, t domain.RefreshToken) (*domain.RefreshSession, error) {
	s, exists := m.data[t]
	if !exists {
		return nil, er.ErrInvalidToken
	}
	return s, nil
}

func (m *mockSessionRepo) Delete(_ context.Context, t domain.RefreshToken) error {
	if _, ok := m.data[t]; ok {
		m.deleteCalls++
		delete(m.data, t)
	}
	return nil
}

func (m *mockSessionRepo) DeleteByUserID(_ context.Context, userID domain.UserID) error {
	for token, s := range m.data {
		if s.UserID == userID {
			delete(m.data, token)
		}
	}
	return nil
}

// --- PasswordHasher ---

type mockHasher struct {
	hash       string
	hashErr    error
	compareErr error
}

func (m *mockHasher) Hash(_ string) (string, error) {
	return m.hash, m.hashErr
}

func (m *mockHasher) Compare(_, _ string) error {
	return m.compareErr
}

// --- TokenIssuer ---

type mockTokenIssuer struct {
	access  domain.AccessToken
	refresh domain.RefreshToken
	err     error
}

func (m *mockTokenIssuer) IssueAccess(_ domain.UserID, _ []domain.Role, _ int64) (domain.AccessToken, error) {
	if m.err != nil {
		return "", m.err
	}
	if m.access == "" {
		return "access-token", nil
	}
	return m.access, nil
}

func (m *mockTokenIssuer) NewRefresh() (domain.RefreshToken, error) {
	if m.err != nil {
		return "", m.err
	}
	if m.refresh == "" {
		return "new-refresh-token", nil
	}
	return m.refresh, nil
}

func (m *mockTokenIssuer) ParseAccess(_ domain.AccessToken) (domain.UserID, []domain.Role, int64, error) {
	return "", nil, 0, nil
}

// --- IDGenerator ---

type mockIDGen struct {
	id string
}

func (m *mockIDGen) New() string {
	return m.id
}

// --- UnitOfWork ---

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// --- helper ---

func mustTime(d time.Duration) time.Time {
	return time.Now().Add(d)
}
