package adal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestSignJwt_HashSwitching(t *testing.T) {
	// 1. Setup: Generate a mock RSA key and self-signed certificate
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test-cert"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	certBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	cert, _ := x509.ParseCertificate(certBytes)

	spt := &ServicePrincipalToken{} // Mock SPT

	tests := []struct {
		name         string
		enableSha256 bool
		wantHeader   string
	}{
		{
			name:         "Default behavior (SHA-1)",
			enableSha256: false,
			wantHeader:   "x5t",
		},
		{
			name:         "Opt-in behavior (SHA-256)",
			enableSha256: true,
			wantHeader:   "x5t#S256",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := &ServicePrincipalCertificateSecret{
				Certificate:  cert,
				PrivateKey:   priv,
				EnableSha256: tt.enableSha256,
			}

			tokenString, err := secret.SignJwt(spt)
			if err != nil {
				t.Fatalf("SignJwt failed: %v", err)
			}

			// 2. Decode the token (without verifying) to check the header
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				t.Fatalf("Failed to parse generated JWT: %v", err)
			}

			// 3. Assert the correct header key exists
			if _, ok := token.Header[tt.wantHeader]; !ok {
				t.Errorf("Expected header %s not found. Headers present: %v", tt.wantHeader, token.Header)
			}

			// 4. Assert the WRONG header key does NOT exist
			wrongHeader := "x5t#S256"
			if tt.enableSha256 {
				wrongHeader = "x5t"
			}
			if _, ok := token.Header[wrongHeader]; ok {
				t.Errorf("Found unexpected header %s in %s mode", wrongHeader, tt.name)
			}
		})
	}
}
