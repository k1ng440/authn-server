package services_test

import (
	"testing"

	"github.com/keratin/authn-server/config"
	"github.com/keratin/authn-server/data/mock"
	"github.com/keratin/authn-server/services"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsVerifierSuccess(t *testing.T) {
	username := "myname"
	password := "mysecret"
	bcrypted := []byte("$2a$04$lzQPXlov4RFLxps1uUGq4e4wmVjLYz3WrqQw4bSdfIiJRyo3/fk3C")

	cfg := config.Config{BcryptCost: 4}
	store := mock.NewAccountStore()
	store.Create(username, bcrypted)

	acc, errs := services.CredentialsVerifier(store, &cfg, username, password)
	if len(errs) > 0 {
		for _, err := range errs {
			assert.NoError(t, err)
		}
	} else {
		assert.NotEqual(t, 0, acc.Id)
		assert.Equal(t, username, acc.Username)
	}
}

func TestCredentialsVerifierFailure(t *testing.T) {
	password := "mysecret"
	bcrypted := []byte("$2a$04$lzQPXlov4RFLxps1uUGq4e4wmVjLYz3WrqQw4bSdfIiJRyo3/fk3C")

	cfg := config.Config{BcryptCost: 4}
	store := mock.NewAccountStore()
	store.Create("known", bcrypted)
	acc, _ := store.Create("locked", bcrypted)
	store.Lock(acc.Id)
	acc, _ = store.Create("expired", bcrypted)
	store.RequireNewPassword(acc.Id)

	testCases := []struct {
		username string
		password string
		errors   []services.Error
	}{
		{"", "", []services.Error{{"credentials", "FAILED"}}},
		{"unknown", "unknown", []services.Error{{"credentials", "FAILED"}}},
		{"known", "unknown", []services.Error{{"credentials", "FAILED"}}},
		{"unknown", password, []services.Error{{"credentials", "FAILED"}}},
		{"locked", password, []services.Error{{"account", "LOCKED"}}},
		{"expired", password, []services.Error{{"credentials", "EXPIRED"}}},
	}

	for _, tc := range testCases {
		_, errs := services.CredentialsVerifier(store, &cfg, tc.username, tc.password)
		assert.Equal(t, tc.errors, errs)
	}
}