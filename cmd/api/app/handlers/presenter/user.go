package presenter

import "github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"

type jsonUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Rol       string `json:"rol"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy string `json:"created_by"`
}

func User(d *user.User) *jsonUser {
	toReturn := &jsonUser{
		ID:        d.ID,
		Username:  d.Username,
		Rol:       getRol(d.IDRol), // TODO
		CreatedAt: d.CreatedAt.Format("02-01-06 15:04"),
		UpdatedAt: d.UpdatedAt.Format("02-01-06 15:04"),
		CreatedBy: "testuser",
	}

	return toReturn
}

func getRol(id uint) string {
	if id == 1 {
		return "Admin"
	}

	return "User"
}
