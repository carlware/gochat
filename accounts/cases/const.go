package cases

import "github.com/carlware/gochat/accounts/models"

// List of fixed users
var USERS = map[string]*models.Profile{}

func init() {
	u1 := models.NewProfile("carlos", "", "1234")
	u2 := models.NewProfile("john", "", "1234")
	u3 := models.NewProfile("gerard", "", "1234")
	USERS[u1.Name] = u1
	USERS[u2.Name] = u2
	USERS[u3.Name] = u3
}
