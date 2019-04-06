package credentials

type UsernamePasswordCredentials struct {
	CredentialsBase `mapstructure:",squash"`
	Username        string
	Password        string
}

func NewUsernamePassword() *UsernamePasswordCredentials {
	cred := &UsernamePasswordCredentials{}
	cred.CredType = "Username/Password"
	return cred
}

func (cred *UsernamePasswordCredentials) ToString() string {
	return cred.BaseToString()
}
