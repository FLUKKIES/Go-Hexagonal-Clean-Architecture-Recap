package oauthprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type facebookOAuthProvider struct {
	config *oauth2.Config
}

func NewFacebookOAuthProvider(clientID, clientSecret, redirectURL string) ports.IOAuthProvider {
	return &facebookOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "public_profile"},
			Endpoint:     facebook.Endpoint,
		},
	}
}

func (f *facebookOAuthProvider) ProviderName() string { return "facebook" }

func (f *facebookOAuthProvider) GetAuthURL(state string) string {
	return f.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (f *facebookOAuthProvider) GetUserProfile(code string) (*ports.OAuthUserProfile, error) {
	// แลก Code เป็น Token
	token, err := f.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// เรียก Facebook Graph API
	client := f.config.Client(context.Background(), token)
	url := fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,email,first_name,last_name,picture.type(large)&access_token=%s",
		token.AccessToken,
	)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &ports.OAuthUserProfile{
		ProviderID: info.ID,
		Email:      info.Email,
		FirstName:  info.FirstName,
		LastName:   info.LastName,
		ProfileUrl: info.Picture.Data.URL,
	}, nil
}
