package entities

type GoogleAuthUser struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GiveName      string `json:"given_name"`
	Id            string `json:"id"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}
