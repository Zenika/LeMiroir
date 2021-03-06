package export

import (
	"github.com/allez-chauffe/marcel/pkg/db"
	"github.com/allez-chauffe/marcel/pkg/db/users"
)

type userPassword struct {
	users.User
	PasswordHash string `json:"passwordHash"`
	PasswordSalt string `json:"passwordSalt"`
}

func listUsers(withPassword bool) (interface{}, error) {
	users, err := db.Users().List()
	if err != nil {
		return nil, err
	}

	if withPassword {
		var usersPassword = make([]userPassword, 0, len(users))

		for _, user := range users {
			usersPassword = append(usersPassword, userPassword{
				User:         user,
				PasswordHash: user.PasswordHash,
				PasswordSalt: user.PasswordSalt,
			})
		}

		return usersPassword, nil
	}

	return users, nil
}

func Users(withPassword bool, outputFile string, pretty bool) error {
	return export(func() (interface{}, error) {
		return listUsers(withPassword)
	}, outputFile, pretty)
}
