package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte{0x00}
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	HandleErr(err)

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func CreateWallet() *Wallet {
	private, public := NewKeyPair()

	wallet := Wallet{private, public}
	return &wallet
}

func PublicKeyHash(pubkey []byte) []byte {
	pubHash := sha256.Sum256(pubkey)

	hashing := ripemd160.New()
	_, err := hashing.Write(pubHash[:])
	HandleErr(err)

	publicRipemd := hashing.Sum(nil)

	return publicRipemd
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w Wallet) Address() []byte {
	pubKeyHash := PublicKeyHash(w.PublicKey)
	versionHash := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionHash)

	fullHash := append(checksum, versionHash...)
	address := Base58Encoding(fullHash)
	return address
}
