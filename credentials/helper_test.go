package credentials

const (
	testCredentialsAsMap = `{
	"test": {
		"type": "secret",
		"description": "test-desc",
		"secret": "my secret"
	},
	"test2": {
		"type": "usernamepassword",
		"description": "test2-desc",
		"username": "my",
		"password": "secret"
	}
}`
	testCredentialsAsList = `[
	{
		"id": "test",
		"type": "secret",
		"description": "test-desc",
		"secret": "my secret"
	},
	{
		"id": "test2",
		"type": "usernamepassword",
		"description": "test2-desc",
		"username": "my",
		"password": "secret"
	}
]`
)

var testCredentials = []Credentials{
	func() (creds *SecretTextCredentials) {
		creds = NewSecretText()
		creds.ID = "test"
		creds.Secret = "my secret"
		creds.Description = "test-desc"
		return
	}(),
	func() (creds *UsernamePasswordCredentials) {
		creds = NewUsernamePassword()
		creds.ID = "test2"
		creds.Username = "my"
		creds.Password = "secret"
		creds.Description = "test2-desc"
		return
	}(),
}
