package security

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"

	"configuration-management/global"

	"github.com/pkg/errors"
)

func LoadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

func LoadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("unsupported public key type")
	}
}

// IsValidSignature 验证签名是否有效
func IsValidSignature(signature, timestamp, encryptedData string) bool {
	publicKey := global.PublicKey // 假设 global.PublicKey 是事先加载好的 RSA 公钥

	// 将 timestamp 和 encryptedData 拼接，然后使用公钥验证签名
	dataToSign := timestamp + encryptedData
	return verifySignature(dataToSign, signature, publicKey)
}

// VerifySignature 验证签名是否有效
func verifySignature(data, signature string, publicKey *rsa.PublicKey) bool {
	// 对数据进行 SHA-256 哈希
	hashed := sha256.Sum256([]byte(data))

	// 对签名进行 Base64 解码
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	// 使用 RSA 公钥验证签名
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
	return err == nil
}

func getSignature(data string, privateKey *rsa.PrivateKey) (string, error) {
	// 对数据进行 SHA-256 哈希
	hashed := sha256.Sum256([]byte(data))

	// 使用 RSA 私钥签名
	signatureBytes, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	// 对签名进行 Base64 编码
	return base64.StdEncoding.EncodeToString(signatureBytes), nil
}

func GetSignature(data string) (string, error) {
	privateKey := global.PrivateKey // 假设 global.PrivateKey 是事先加载好的 RSA 私钥

	// 对数据进行 SHA-256 哈希
	hashed := sha256.Sum256([]byte(data))

	// 使用 RSA 私钥签名
	signatureBytes, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	// 对签名进行 Base64 编码
	return base64.StdEncoding.EncodeToString(signatureBytes), nil
}
