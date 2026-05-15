package gormrepo

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"gorm.io/gorm"
)

type oauthRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuthRepositoryImpl(db *gorm.DB) repositories.IOAuthRepository {
	return &oauthRepositoryImpl{db: db}
}

func (r *oauthRepositoryImpl) FindByProviderAndID(provider, providerID string) (*entities.OAuthAccount, error) {
	m := new(entities.OAuthAccount)
	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *oauthRepositoryImpl) Create(account *entities.OAuthAccount) error {
	return r.db.Create(account).Error
}
