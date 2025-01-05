package transactions

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

// Wallet represents a user's blockchain wallet
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  string
	Address    string
}

// CreateWallet generates a new blockchain wallet
func CreateWallet() *Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := hex.EncodeToString(append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...))
	address := GenerateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}

// GenerateAddress creates a short hash-based address from the public key
func GenerateAddress(publicKey string) string {
	hash := sha256.Sum256([]byte(publicKey))
	return hex.EncodeToString(hash[:])[:20] // First 20 chars for simplicity
}
