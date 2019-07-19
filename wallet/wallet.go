package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubKeyHash := PublicKeyHash(w.PublicKey)

	versionHash := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionHash)

	fullHash := append(versionHash, checksum...)
	address := Base58Encode(fullHash)

	return address
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

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	HandleErr(err)

	publicRipemd := hasher.Sum(nil)

	return publicRipemd
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
