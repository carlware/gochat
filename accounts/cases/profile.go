package cases

import (
	"context"
	"errors"

	"github.com/carlware/gochat/accounts/models"
	"github.com/carlware/gochat/common/auth"
)

func List(ctx context.Context) []*models.Profile {
	profiles := []*models.Profile{}
	for _, v := range USERS {
		profiles = append(profiles, v)
	}
	return profiles
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(req *LoginRequest) (*auth.JWTToken, error) {
	if profile, ok := USERS[req.Username]; ok {
		if profile.PasswordMatch(req.Password) {
			fields := map[string]interface{}{
				"uid":      profile.ID,
				"username": profile.Name,
			}
			return auth.GenerateJWT(fields)
		}
		return nil, errors.New("password does not match")
	}
	return nil, errors.New("User does not exist")
}
