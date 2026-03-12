package app

import (
	"testing"
)

func TestBcryptAlgorithm(t *testing.T) {
	alg := &BcryptAlgorithm{}

	// Test generation
	hash, err := alg.Generate("testpassword", 10)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	// Verify it starts with bcrypt prefix
	if len(hash) < 60 || hash[:4] != "$2a$" {
		t.Errorf("Invalid bcrypt hash format: %s", hash)
	}

	// Test verification
	err = alg.Verify("testpassword", hash)
	if err != nil {
		t.Errorf("Verification failed: %v", err)
	}

	// Test wrong password
	err = alg.Verify("wrongpassword", hash)
	if err == nil {
		t.Error("Verification should have failed for wrong password")
	}
}

func TestScryptAlgorithm(t *testing.T) {
	alg := &ScryptAlgorithm{}

	// Test generation
	hash, err := alg.Generate("testpassword", HashParams{
		ScryptN: 16384,
		ScryptR: 8,
		ScryptP: 1,
	})
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	// Verify it starts with scrypt prefix
	if hash[:8] != "$scrypt$" {
		t.Errorf("Invalid scrypt hash format: %s", hash)
	}

	// Test verification
	err = alg.Verify("testpassword", hash)
	if err != nil {
		t.Errorf("Verification failed: %v", err)
	}

	// Test wrong password
	err = alg.Verify("wrongpassword", hash)
	if err == nil {
		t.Error("Verification should have failed for wrong password")
	}
}

func TestGetHashAlgorithm(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"bcrypt", "bcrypt"},
		{"argon2", "argon2"},
		{"argon2id", "argon2"},
		{"scrypt", "scrypt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alg, err := GetHashAlgorithm(tt.name)
			if err != nil {
				t.Errorf("GetHashAlgorithm(%s) error: %v", tt.name, err)
				return
			}
			if alg.Name() != tt.expected {
				t.Errorf("GetHashAlgorithm(%s) = %s, want %s", tt.name, alg.Name(), tt.expected)
			}
		})
	}

	t.Run("unknown algorithm", func(t *testing.T) {
		_, err := GetHashAlgorithm("unknown")
		if err == nil {
			t.Error("Expected error for unknown algorithm")
		}
	})
}

func TestDetectHashAlgorithm(t *testing.T) {
	bcryptHash := "$2a$10$abcdefghijklmnopqrstuv"
	scryptHash := "$scrypt$N=16384,r=8,p=1$salt$hash"

	t.Run("bcrypt", func(t *testing.T) {
		alg, err := DetectHashAlgorithm(bcryptHash)
		if err != nil {
			t.Errorf("DetectHashAlgorithm(bcrypt) error: %v", err)
			return
		}
		if alg.Name() != "bcrypt" {
			t.Errorf("Expected bcrypt, got %s", alg.Name())
		}
	})

	t.Run("scrypt", func(t *testing.T) {
		alg, err := DetectHashAlgorithm(scryptHash)
		if err != nil {
			t.Errorf("DetectHashAlgorithm(scrypt) error: %v", err)
			return
		}
		if alg.Name() != "scrypt" {
			t.Errorf("Expected scrypt, got %s", alg.Name())
		}
	})
}

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		algorithm string
		params    HashParams
		prefix    string
	}{
		{"bcrypt", HashParams{BcryptCost: 10}, "$2a$"},
		{"scrypt", HashParams{ScryptN: 16384, ScryptR: 8, ScryptP: 1}, "$scrypt$"},
	}

	for _, tt := range tests {
		t.Run(tt.algorithm, func(t *testing.T) {
			hash, err := GeneratePasswordHash("password", tt.algorithm, tt.params)
			if err != nil {
				t.Fatalf("GeneratePasswordHash error: %v", err)
			}
			if len(hash) < 60 {
				t.Errorf("Hash too short: %s", hash)
			}
		})
	}
}

func TestVerifyPasswordHash(t *testing.T) {
	// Generate a real bcrypt hash
	alg := &BcryptAlgorithm{}
	realHash, err := alg.Generate("test", 10)
	if err != nil {
		t.Fatalf("Failed to generate bcrypt hash: %v", err)
	}

	err = VerifyPasswordHash("test", realHash)
	if err != nil {
		t.Errorf("VerifyPasswordHash failed: %v", err)
	}
}

func TestArgon2Algorithm(t *testing.T) {
	// Skip if argon2 is not available
	// This is a basic test to ensure the algorithm exists
	alg := &Argon2Algorithm{}

	if alg == nil {
		t.Error("Argon2Algorithm should not be nil")
	}

	if alg.Name() != "argon2" {
		t.Errorf("Expected name 'argon2', got '%s'", alg.Name())
	}

	if alg.Prefix() != "$argon2id$" {
		t.Errorf("Expected prefix '$argon2id$', got '%s'", alg.Prefix())
	}
}
