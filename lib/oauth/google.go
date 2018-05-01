package oauth

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
	_google "golang.org/x/oauth2/google"
)

// NewGoogleProvider returns a AuthN integration for Google OAuth
func NewGoogleProvider(credentials *Credentials) *Provider {
	config := &oauth2.Config{
		ClientID:     credentials.ID,
		ClientSecret: credentials.Secret,
		Scopes:       []string{"email"},
		Endpoint:     _google.Endpoint,
	}

	return &Provider{
		config: config,
		UserInfo: func(t *oauth2.Token) (*UserInfo, error) {
			client := config.Client(context.TODO(), t)
			resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			// {
			//   "id": "1234567890",
			//   "email": "user@example.com",
			//   "verified_email": true,
			//   "name": "Example User",
			//   "given_name": "Example",
			//   "family_name": "User",
			//   "link": "https://plus.google.com/1234567890",
			//   "picture": "https://lh6.googleusercontent.com/path/to/photo.jpg"
			// }
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			var user UserInfo
			err = json.Unmarshal(body, &user)
			return &user, err
		},
	}
}