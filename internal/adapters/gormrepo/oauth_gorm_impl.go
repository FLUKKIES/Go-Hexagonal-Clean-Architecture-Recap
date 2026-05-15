package gormrepo

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo/model"
	"gorm.io/gorm"
)

type oauthRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuthRepositoryImpl(db *gorm.DB) repositories.IOAuthRepository {
	return &oauthRepositoryImpl{db: db}
}

func (r *oauthRepositoryImpl) FindByProviderAndID(provider, providerID string) (*entities.OAuthAccount, error) {
	m := new(model.OAuthAccountGormModel)
	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.OAuthAccount{
		ID:         m.ID,
		UserID:     m.UserID,
		Provider:   entities.OAuthProvider(m.Provider),
		ProviderID: m.ProviderID,
		Email:      m.Email,
		CreatedAt:  m.CreatedAt,
	}, nil
}

func (r *oauthRepositoryImpl) Create(account *entities.OAuthAccount) error {
	m := &model.OAuthAccountGormModel{
		ID:         account.ID,
		UserID:     account.UserID,
		Provider:   string(account.Provider),
		ProviderID: account.ProviderID,
		Email:      account.Email,
		CreatedAt:  account.CreatedAt,
	}
	return r.db.Create(m).Error
}
