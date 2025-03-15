package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Convert an arbitrary string to a big.Int
func StringToBigInt(s string) (*big.Int, error) {
	return new(big.Int).SetBytes([]byte(s)), nil
}

// Convert a big.Int back to a string
func BigIntToString(b *big.Int) string {
	return string(b.Bytes())
}

func BigIntToString_Key(b *big.Int) string {
	return b.Text(16)
}

func StringToBigInt_Key(s string) (*big.Int, error) {
	x, _ := new(big.Int).SetString(s, 16)
	if x == nil {
		return nil, fmt.Errorf("invalid big.Int string")
	}
	return x, nil
}

func StringToBigInt_Ciphertext(s string) (*big.Int, error) {
	return StringToBigInt_Key(s)
}

func BigIntToString_Ciphertext(b *big.Int) string {
	return BigIntToString_Key(b)
}

////////////////////////////////
///// ELGAMAL CRYPTOGRAPHY /////
////////////////////////////////

func GetPrime() *big.Int {
	num := new(big.Int)
	num.SetString("297138044045261491373545160878212841086754727264356427752279134139093393432479232426981708915566446228157986733178792358394397732218358050476422997921755483991327064324936568323843882602983055776986371112569464915603732489164224351368123939071577554288188235637692777077409044054418418789205827922754919745591", 10)
	return num
}

func GetGenerator() *big.Int {
	num := new(big.Int)
	num.SetString("2", 10)
	return num
}

// This function is only used when there is need
// for generating temporary private keys
func GeneratePrivateKey() (*big.Int, error) {
	p := GetPrime()

	a, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(169122285939163)))
	if err != nil {
		return nil, err
	}

	return a, nil
}

func GeneratePrivateKeyString() (string, error) {
	a, err := GeneratePrivateKey()
	if err != nil {
		return "", err
	}

	return BigIntToString_Key(a), nil
}

func ComputePublicKey(privateKey *big.Int) (*big.Int, error) {
	p := GetPrime()
	alpha := GetGenerator()
	a := privateKey

	beta := new(big.Int).Exp(alpha, a, p)

	return beta, nil
}

func ComputePublicKeyString(privateKeyString string) (string, error) {
	privateKey, err := StringToBigInt_Key(privateKeyString)
	if err != nil {
		return "", err
	}

	beta, err := ComputePublicKey(privateKey)
	if err != nil {
		return "", err
	}

	return BigIntToString_Key(beta), nil
}

func ElGamalEncrypt(beta, x *big.Int) (*big.Int, *big.Int, error) {
	p := GetPrime()
	alpha := GetGenerator()

	// Generate random k
	k, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(169122285939163)))
	if err != nil {
		return nil, nil, err
	}

	y1 := new(big.Int).Exp(alpha, k, p)

	beta_to_power_k := new(big.Int).Exp(beta, k, p)
	y2 := new(big.Int).Mod(new(big.Int).Mul(x, beta_to_power_k), p)

	return y1, y2, nil
}

func ElGamalEncryptString(publicKeyString, xString string) (string, string, error) {
	beta, err := StringToBigInt_Key(publicKeyString)
	if err != nil {
		return "", "", err
	}

	x, err := StringToBigInt(xString)
	if err != nil {
		return "", "", err
	}

	y1, y2, err := ElGamalEncrypt(beta, x)
	if err != nil {
		return "", "", err
	}

	return BigIntToString_Ciphertext(y1), BigIntToString_Ciphertext(y2), nil
}

func ElGamalDecrypt(a, y1, y2 *big.Int) (*big.Int, error) {
	p := GetPrime()
	y1_to_power_a := new(big.Int).Exp(y1, a, p)
	y1_to_power_minus_a := new(big.Int).ModInverse(y1_to_power_a, p)
	x := new(big.Int).Mod(new(big.Int).Mul(y2, y1_to_power_minus_a), p)
	return x, nil
}

func ElGamalDecryptString(privateKeyString, y1String, y2String string) (string, error) {
	a, err := StringToBigInt_Key(privateKeyString)
	if err != nil {
		return "", err
	}

	y1, err := StringToBigInt_Ciphertext(y1String)
	if err != nil {
		return "", err
	}

	y2, err := StringToBigInt_Ciphertext(y2String)
	if err != nil {
		return "", err
	}

	x, err := ElGamalDecrypt(a, y1, y2)
	if err != nil {
		return "", err
	}

	return BigIntToString(x), nil
}

/////////////////////////////////
/// PROXY RE-ENCRYPTION (PRE) ///
/////////////////////////////////

func PREComputeReEncryptionKey(studentPrivateKey, qrPublicKey *big.Int) (*big.Int, error) {
	a := studentPrivateKey
	g_to_power_b := qrPublicKey
	p := GetPrime()

	H_g_to_power_a_b := new(big.Int).Exp(g_to_power_b, a, p)
	inverse_H_g_to_power_a_b := new(big.Int).ModInverse(H_g_to_power_a_b, p)
	rk := new(big.Int).Mod(new(big.Int).Mul(a, inverse_H_g_to_power_a_b), p)

	return rk, nil
}

func PREComputeReEncryptionKeyString(studentPrivateKeyString, qrPublicKeyString string) (string, error) {
	studentPrivateKey, err := StringToBigInt_Key(studentPrivateKeyString)
	if err != nil {
		return "", err
	}

	qrPublicKey, err := StringToBigInt_Key(qrPublicKeyString)
	if err != nil {
		return "", err
	}

	rk, err := PREComputeReEncryptionKey(studentPrivateKey, qrPublicKey)
	if err != nil {
		return "", err
	}

	return BigIntToString_Key(rk), nil
}

func PREEncrypt(rk, y1, y2 *big.Int) (*big.Int, *big.Int, error) {
	p := GetPrime()

	z1 := new(big.Int).Exp(y1, rk, p)
	z2 := y2

	return z1, z2, nil
}

func PREEncryptString(reEncryptionKeyString, y1String, y2String string) (string, string, error) {
	rk, err := StringToBigInt_Key(reEncryptionKeyString)
	if err != nil {
		return "", "", err
	}

	y1, err := StringToBigInt_Ciphertext(y1String)
	if err != nil {
		return "", "", err
	}

	y2, err := StringToBigInt_Ciphertext(y2String)
	if err != nil {
		return "", "", err
	}

	z1, z2, err := PREEncrypt(rk, y1, y2)
	if err != nil {
		return "", "", err
	}

	return BigIntToString(z1), BigIntToString(z2), nil
}

func PREDecrypt(studentPublicKey, qrPrivateKey, z1, z2 *big.Int) (*big.Int, error) {
	p := GetPrime()
	g_to_power_a := studentPublicKey
	b := qrPrivateKey

	H_g_to_power_a_b := new(big.Int).Exp(g_to_power_a, b, p)
	t := new(big.Int).Exp(z1, H_g_to_power_a_b, p)
	inverse_t := new(big.Int).ModInverse(t, p)
	x := new(big.Int).Mod(new(big.Int).Mul(z2, inverse_t), p)

	return x, nil
}

func PREDecryptString(studentPublicKeyString, qrPrivateKeyString, z1String, z2String string) (string, error) {
	studentPublicKey, err := StringToBigInt_Key(studentPublicKeyString)
	if err != nil {
		return "", err
	}

	qrPrivateKey, err := StringToBigInt_Key(qrPrivateKeyString)
	if err != nil {
		return "", err
	}

	z1, err := StringToBigInt_Ciphertext(z1String)
	if err != nil {
		return "", err
	}

	z2, err := StringToBigInt_Ciphertext(z2String)
	if err != nil {
		return "", err
	}

	x, err := PREDecrypt(studentPublicKey, qrPrivateKey, z1, z2)
	if err != nil {
		return "", err
	}

	return BigIntToString(x), nil
}

////////////////////////////////
////// VIGENÈRE CIPHER /////////
////////////////////////////////

// Function to repeat the key to match the length of the text
func repeatKey(key string, length int) string {
	repeatedKey := []rune{}
	keyRunes := []rune(key)
	keyLen := len(keyRunes)

	for i := 0; i < length; i++ {
		repeatedKey = append(repeatedKey, keyRunes[i%keyLen])
	}
	return string(repeatedKey)
}

// Function to encrypt text using Vigenère Cipher (supports all Unicode characters)
func VigenereEncryptString(key, text string) (string, error) {
	textRunes := []rune(text)
	keyRunes := []rune(repeatKey(key, len(textRunes)))
	encrypted := make([]rune, len(textRunes))

	for i := 0; i < len(textRunes); i++ {
		encrypted[i] = (textRunes[i] + keyRunes[i]) % 0x10FFFF // Ensuring valid Unicode range
	}
	return string(encrypted), nil
}

// Function to decrypt text using Vigenère Cipher (supports all Unicode characters)
func VigenereDecryptString(key, text string) (string, error) {
	textRunes := []rune(text)
	keyRunes := []rune(repeatKey(key, len(textRunes)))
	decrypted := make([]rune, len(textRunes))

	for i := 0; i < len(textRunes); i++ {
		decrypted[i] = (textRunes[i] - keyRunes[i] + 0x10FFFF) % 0x10FFFF // Reverse shift
	}
	return string(decrypted), nil
}

/////////////////////////////////
//// MORE UTILITY FUNCTIONS /////
/////////////////////////////////

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
