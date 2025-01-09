package user

// User is struct used to decode user data from requests.
type User struct {
	// The email of the user.
	Email string `json:"email"`
	// The password of the user.
	Password string `json:"password"`
}
