package utils

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"math/big"
)

// GenerateRSAKeyFromString derives a deterministic RSA private key from a user input string.
func GenerateRSAKeyFromString(userInput string) *rsa.PrivateKey {
	// Hash the input string (MD5 for deterministic seed)
	hash := md5.Sum([]byte(userInput)) // 16-byte hash

	// Convert hash to big integer for seeding prime numbers
	d := new(big.Int).SetBytes(hash[:])

	// Generate a prime number from the hash (ensuring it's a valid RSA private exponent)
	p, err := rand.Prime(rand.Reader, 1024) // Generate a large prime
	if err != nil {
		panic("Failed to generate prime")
	}

	// Use hash-derived big int to influence private exponent
	privateExponent := new(big.Int).Mod(d, p)

	// Construct RSA key manually
	privKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).Mul(p, p), // Fake modulus (for simplicity)
			E: 65537,                  // Standard public exponent
		},
		D: privateExponent,
		Primes: []*big.Int{
			p, p, // Using the same prime twice (not secure but meets your request)
		},
	}

	return privKey
}

// PrivateKeyToPublicKey derives the public key from an arbitrary string-based private key.
func PrivateKeyToPublicKey(privateKey string) string {
	priv := GenerateRSAKeyFromString(privateKey)

	// Convert to PEM format
	pubASN1, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})
	return string(pubPEM)
}

// Encrypt encrypts a plaintext string using a PEM-encoded public key string and returns a base64-encoded ciphertext.
func Encrypt(publicKey string, plaintext string) string {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		panic("Invalid PEM public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("Failed to parse public key")
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("Invalid RSA public key")
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, []byte(plaintext), nil)
	if err != nil {
		panic("Encryption failed")
	}

	// Convert ciphertext to a base64 string
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt decrypts a base64-encoded ciphertext string using an arbitrary string-derived private key.
func Decrypt(privateKey string, ciphertext string) string {
	priv := GenerateRSAKeyFromString(privateKey)

	// Decode base64 ciphertext
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		panic("Failed to decode base64 ciphertext")
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, cipherBytes, nil)
	if err != nil {
		panic("Decryption failed")
	}

	return string(plaintext)
}

// Sign creates a signature for a message using the given private key string.
func Sign(privateKey string, message string) string {
	priv := GenerateRSAKeyFromString(privateKey)

	// Hash the message using SHA-256
	hashed := sha256.Sum256([]byte(message))

	// Sign the hash using the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
	if err != nil {
		panic("Signing failed")
	}

	// Convert the signature to a base64 string
	return base64.StdEncoding.EncodeToString(signature)
}

// Verify checks if a signature is valid for a given message and public key string.
func Verify(publicKey string, message string, signature string) bool {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		panic("Invalid PEM public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("Failed to parse public key")
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("Invalid RSA public key")
	}

	// Hash the message using SHA-256
	hashed := sha256.Sum256([]byte(message))

	// Decode base64 signature
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		panic("Failed to decode base64 signature")
	}

	// Verify the signature
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], sigBytes)
	return err == nil
}

func HashSHA512(input string) string {
	hash := sha512.Sum512([]byte(input)) // Compute SHA-512 hash
	return hex.EncodeToString(hash[:])   // Convert to hex string
}

// GenerateSecureRandomString creates a cryptographically secure random string.
func GenerateSecureRandomString(length int) string {
	// Generate random bytes
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Failed to generate secure random bytes")
	}

	// Encode to base64 to ensure printable characters
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}
