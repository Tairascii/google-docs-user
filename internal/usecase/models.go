package usecase

type Tokens struct {
	Access  string
	Refresh string
}

type SignUpData struct {
	Name          string
	Email         string
	Password      string
	ProfilePicUrl string
}
