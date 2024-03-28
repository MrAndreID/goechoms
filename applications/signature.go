package applications

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"hash"
	"io"

	"github.com/MrAndreID/goechoms/configs"

	"github.com/sirupsen/logrus"
)

// RSA Ciphers : RSA/NONE/OAEPWithSHA1AndMGF1Padding
// For Testing : https://8gwifi.org/RSAFunctionality?keysize=512
func (app *Application) DecryptSignature(cfg *configs.Config, encryptedData string) (*string, error) {
	var (
		tag    string    = "Applications.Signature.DecryptSignature."
		hash   hash.Hash = sha1.New()
		random io.Reader = rand.Reader
	)

	privateKeyBlock, _ := pem.Decode([]byte(cfg.RSAOAEPKey))

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to load private key")

		return nil, err
	}

	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to decode base64 for encrypted data")

		return nil, err
	}

	decryptedData, err := rsa.DecryptOAEP(hash, random, privateKey, decodedData, nil)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": err.Error(),
		}).Error("failed to decrypt data")

		return nil, err
	}

	resultData := string(decryptedData)

	return &resultData, nil
}

func (app *Application) VerifySignature(cfg *configs.Config, payload string) (*string, error) {
	mac := hmac.New(sha512.New, []byte(cfg.SecretKey))

	_, err := mac.Write([]byte(payload))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   "Applications.Signature.VerifySignature.01",
			"error": err.Error(),
		}).Error("failed to write from payload")

		return nil, err
	}

	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &sign, nil
}
