package repositories

import "github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"

type IOAuthRepository interface {
	FindByProviderAndID(provider, providerID string) (*entities.OAuthAccount, error)
	Create(account *entities.OAuthAccount) error
}
