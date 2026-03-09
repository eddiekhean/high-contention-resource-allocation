package service

type TokenStorageService struct {
	// Add fields as needed, e.g., redis client
}

func NewTokenStorageService() *TokenStorageService {
	return &TokenStorageService{}
}

func (s *TokenStorageService) IsValid(jti string) bool {
	// Basic implementation or stub for now
	return true
}
