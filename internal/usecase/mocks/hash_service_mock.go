package mocks

type HashServiceMock struct {
	HashPasswordFunc      func(password string) (string, error)
	CheckPasswordHashFunc func(password, hash string) bool
}

func (m *HashServiceMock) HashPassword(password string) (string, error) {
	if m.HashPasswordFunc != nil {
		return m.HashPasswordFunc(password)
	}
	return "hashed_password", nil
}

func (m *HashServiceMock) CheckPasswordHash(password, hash string) bool {
	if m.CheckPasswordHashFunc != nil {
		return m.CheckPasswordHashFunc(password, hash)
	}
	return true
}
