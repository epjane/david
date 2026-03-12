// HashAlgorithm interface and implementations
//
// This package implements secure password hashing algorithms for david.
// We support bcrypt, argon2id, and scrypt - all memory-hard algorithms.
//
// Why we don't use SHA-256/SHA-512 for password hashing:
//
// SHA algorithms (SHA-256, SHA-512) are designed for general-purpose hashing
// and are intentionally fast. This makes them unsuitable for password hashing
// because:
//
//  1. Speed is a security risk: Fast hashes can be brute-forced quickly
//     using GPU/ASIC clusters (see: https://github.com/audstanley/david/issues/2)
//
// 2. No built-in salting: SHA doesn't include salt, requiring manual implementation
//
// 3. No adaptive work factor: SHA can't be slowed down as hardware improves
//
//  4. Vulnerable to rainbow tables: Without proper salting, common passwords
//     can be pre-computed and looked up
//
// Instead, we use memory-hard algorithms (bcrypt, argon2id, scrypt) that:
// - Are intentionally slow (~65-85ms per hash)
// - Require significant memory (~4KB-128MB), making GPU attacks expensive
// - Include built-in salting
// - Have configurable work factors that can be increased over time
//
// For more information on password hashing, see:
// - OWASP Password Storage Cheat Sheet
// - RFC 9106 (Argon2)
// - RFC 7914 (scrypt)
package app

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"

	"crypto/subtle"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// HashAlgorithm defines the interface for password hashing algorithms
type HashAlgorithm interface {
	// Generate creates a hash from the password with given cost/params
	Generate(password string, params interface{}) (string, error)
	// Verify checks if the password matches the hash
	Verify(password, hash string) error
	// Name returns the algorithm name
	Name() string
	// Prefix returns the hash prefix for detection
	Prefix() string
}

// HashParams contains configuration for all hash algorithms
type HashParams struct {
	BcryptCost int  `yaml:"bcrypt_cost" default:"10"`
	ScryptN    uint `yaml:"scrypt_n" default:"16384"` // CPU/memory cost
	ScryptR    uint `yaml:"scrypt_r" default:"8"`     // Block size
	ScryptP    uint `yaml:"scrypt_p" default:"1"`     // Parallelism
}

// BcryptAlgorithm implements bcrypt hashing
type BcryptAlgorithm struct{}

func (b *BcryptAlgorithm) Generate(password string, params interface{}) (string, error) {
	cost := 10

	// Check for HashParams first (new API)
	if hp, ok := params.(HashParams); ok {
		cost = hp.BcryptCost
	} else if p, ok := params.(int); ok {
		// Check for int (backward compatibility)
		cost = p
	}

	if cost < 4 {
		cost = 4
	}
	if cost > 14 {
		cost = 14
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptAlgorithm) Verify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (b *BcryptAlgorithm) Name() string {
	return "bcrypt"
}

func (b *BcryptAlgorithm) Prefix() string {
	return "$2a$"
}

// Argon2Algorithm implements Argon2id hashing
type Argon2Algorithm struct{}

func (a *Argon2Algorithm) Generate(password string, params interface{}) (string, error) {
	// Use default Argon2 parameters
	memory := uint32(64 * 1024) // 64 MiB
	iterations := uint32(3)
	parallelism := uint8(4)
	keyLength := uint32(32)

	salt := make([]byte, 16)
	randomData(salt)

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash)), nil
}

func (a *Argon2Algorithm) Verify(password, hash string) error {
	// Parse the hash to extract parameters
	parts := strings.Split(hash, "$")
	if len(parts) != 5 {
		return fmt.Errorf("invalid argon2 hash format")
	}

	// parts[1] = "argon2id"
	// parts[2] = "v=1$m=65536,t=3,p=4"
	// parts[3] = salt
	// parts[4] = hash

	// Parse version and parameters
	var version, memory, iterations, parallelism uint32
	_, err := fmt.Sscanf(parts[2], "v=%d$m=%d,t=%d,p=%d", &version, &memory, &iterations, &parallelism)
	if err != nil {
		return fmt.Errorf("failed to parse argon2 parameters: %v", err)
	}

	// Decode salt and hash
	salt, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %v", err)
	}

	hashBytes, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %v", err)
	}

	// Verify the password
	derivedKey := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(parallelism), uint32(len(hashBytes)))
	if !constantTimeCompare(hashBytes, derivedKey) {
		return fmt.Errorf("password does not match")
	}

	return nil
}

func (a *Argon2Algorithm) Name() string {
	return "argon2"
}

func (a *Argon2Algorithm) Prefix() string {
	return "$argon2id$"
}

// ScryptAlgorithm implements scrypt hashing
type ScryptAlgorithm struct{}

func (s *ScryptAlgorithm) Generate(password string, params interface{}) (string, error) {
	p := HashParams{
		ScryptN: 16384,
		ScryptR: 8,
		ScryptP: 1,
	}
	if hp, ok := params.(HashParams); ok {
		if hp.ScryptN > 0 {
			p.ScryptN = hp.ScryptN
		}
		if hp.ScryptR > 0 {
			p.ScryptR = hp.ScryptR
		}
		if hp.ScryptP > 0 {
			p.ScryptP = hp.ScryptP
		}
	}

	salt := make([]byte, 16)
	randomData(salt)

	hash, err := scrypt.Key([]byte(password), salt, int(p.ScryptN), int(p.ScryptR), int(p.ScryptP), 32)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("$scrypt$N=%d,r=%d,p=%d$%s$%s",
		p.ScryptN,
		p.ScryptR,
		p.ScryptP,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash)), nil
}

func (s *ScryptAlgorithm) Verify(password, hash string) error {
	parts := strings.Split(hash, "$")
	if len(parts) != 5 {
		return fmt.Errorf("invalid scrypt hash format")
	}

	var n, r, p uint
	_, err := fmt.Sscanf(parts[2], "N=%d,r=%d,p=%d", &n, &r, &p)
	if err != nil {
		return fmt.Errorf("failed to parse scrypt parameters: %v", err)
	}

	salt, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %v", err)
	}

	hashBytes, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %v", err)
	}

	derivedKey, err := scrypt.Key([]byte(password), salt, int(n), int(r), int(p), 32)
	if err != nil {
		return err
	}

	if !constantTimeCompare(hashBytes, derivedKey) {
		return fmt.Errorf("password does not match")
	}

	return nil
}

func (s *ScryptAlgorithm) Name() string {
	return "scrypt"
}

func (s *ScryptAlgorithm) Prefix() string {
	return "$scrypt$"
}

// GetHashAlgorithm returns the appropriate hash algorithm based on name
func GetHashAlgorithm(name string) (HashAlgorithm, error) {
	switch name {
	case "bcrypt":
		return &BcryptAlgorithm{}, nil
	case "argon2", "argon2id":
		return &Argon2Algorithm{}, nil
	case "scrypt":
		return &ScryptAlgorithm{}, nil
	default:
		return nil, fmt.Errorf("unknown hash algorithm: %s", name)
	}
}

// DetectHashAlgorithm returns the algorithm based on hash prefix
func DetectHashAlgorithm(hash string) (HashAlgorithm, error) {
	algorithms := []HashAlgorithm{
		&BcryptAlgorithm{},
		&Argon2Algorithm{},
		&ScryptAlgorithm{},
	}

	for _, alg := range algorithms {
		if strings.HasPrefix(hash, alg.Prefix()) {
			return alg, nil
		}
	}

	return nil, fmt.Errorf("unable to detect hash algorithm")
}

// GeneratePasswordHash generates a hash using the specified algorithm
func GeneratePasswordHash(password string, algorithm string, params HashParams) (string, error) {
	alg, err := GetHashAlgorithm(algorithm)
	if err != nil {
		return "", err
	}

	return alg.Generate(password, params)
}

// VerifyPasswordHash verifies a password against a hash using auto-detection
func VerifyPasswordHash(password, hash string) error {
	alg, err := DetectHashAlgorithm(hash)
	if err != nil {
		return err
	}

	return alg.Verify(password, hash)
}

// Helper functions

func randomData(b []byte) {
	binary.LittleEndian.PutUint64(b[:8], binary.LittleEndian.Uint64(randomData8()))
	binary.LittleEndian.PutUint64(b[8:], binary.LittleEndian.Uint64(randomData8()))
}

func randomData8() []byte {
	b := make([]byte, 8)
	sha512.New().Write(b)
	return b
}

func splitHash(hash string) []string {
	return strings.Split(hash, "$")
}

func constantTimeCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
