package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

// Share representa uma parte do segredo (x, y)
type Share struct {
	X *big.Int
	Y *big.Int
}

// Interpolação de Lagrange para reconstruir o segredo (avaliado em x = 0)
func interpolate(shares []Share, prime *big.Int) *big.Int {
	secret := big.NewInt(0)

	for i, shareI := range shares {
		num := big.NewInt(1)
		den := big.NewInt(1)

		for j, shareJ := range shares {
			if i == j {
				continue
			}
			num.Mul(num, new(big.Int).Neg(shareJ.X))
			num.Mod(num, prime)
			den.Mul(den, new(big.Int).Sub(shareI.X, shareJ.X))
			den.Mod(den, prime)
		}

		// Compute term = y * num * den^-1 mod prime
		denInv := new(big.Int).ModInverse(den, prime)
		term := new(big.Int).Mul(shareI.Y, num)
		term.Mul(term, denInv)
		term.Mod(term, prime)

		secret.Add(secret, term)
		secret.Mod(secret, prime)
	}

	return secret
}

func main() {

	prime, _ := new(big.Int).SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed", 16)

	// Shares
	share2 := Share{
		X: big.NewInt(2),
		Y: new(big.Int),
	}
	share2.Y.SetString("0679aa25b8fc83fd391cc282edf686bdcd8c422c799326aaf12052b6d74c3249", 16)

	share3 := Share{
		X: big.NewInt(3),
		Y: new(big.Int),
	}
	share3.Y.SetString("0ac0a3eca4e806e45ce12cc64068407b774e6e75dfd8e80aa543aa5b53832b45", 16)

	// Reconstruir o segredo usando Lagrange
	shares := []Share{share2, share3}
	secret := interpolate(shares, prime)

	// SHA256 do segredo como seed
	secretBytes := secret.Bytes()

	// Garante que tenha 32 bytes (preenche com zeros à esquerda se necessário)
	seed := sha256.Sum256(append(make([]byte, 32-len(secretBytes)), secretBytes...))

	// Gerar chave privada e pública
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	// Codificar chave pública em Base64
	publicKeyB64 := base64.StdEncoding.EncodeToString(publicKey)

	fmt.Println("Chave pública base64:", publicKeyB64)
}
