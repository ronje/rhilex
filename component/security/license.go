// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func InitSecurityLicense() {
	GenLocalSecurityLicense()
}

// generateRSAKeyPair 生成RSA密钥对并将私钥和公钥保存到文件
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// savePEMKey 将密钥保存为PEM格式的文件
func savePEMKey(fileName string, keyType string, keyBytes []byte) error {
	// 创建PEM块
	pemKey := &pem.Block{
		Type:  keyType,
		Bytes: keyBytes,
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return pem.Encode(file, pemKey)
}

// encryptWithPublicKey 使用公钥加密数据
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	label := []byte("")
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		msg,
		label,
	)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

// decryptWithPrivateKey 使用私钥解密数据
func DecryptWithPrivateKey(encryptedBytes []byte, priv *rsa.PrivateKey) ([]byte, error) {
	label := []byte("")
	decryptedBytes, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		priv,
		encryptedBytes,
		label,
	)
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}

// GenLocalSecurityLicense
func GenLocalSecurityLicense() error {
	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		fmt.Println("Error generating RSA key pair:", err)
		return err
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err = savePEMKey("./.encrypt.priv", "RSA PRIVATE KEY", privateKeyBytes)
	if err != nil {
		fmt.Println("Error saving private key:", err)
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error marshaling public key:", err)
		return err
	}
	err = savePEMKey("./.encrypt.pub", "RSA PUBLIC KEY", publicKeyBytes)
	if err != nil {
		fmt.Println("Error saving public key:", err)
		return err
	}
	return nil
}

// ReadLocalKeys 读取本地的私钥和公钥文件，并返回它们的字符串表示
func ReadLocalKeys() (string, string, error) {
	privateKeyPEM, err := os.ReadFile("./.encrypt.priv")
	if err != nil {
		return "", "", fmt.Errorf("failed to read private key file: %v", err)
	}
	publicKeyPEM, err := os.ReadFile("./.encrypt.pub")
	if err != nil {
		return "", "", fmt.Errorf("failed to read public key file: %v", err)
	}
	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// CheckLicense 检查本地证书文件是否存在
func CheckLicense() bool {
	_, err := os.Stat("./.encrypt.priv")
	if os.IsNotExist(err) {
		fmt.Printf("Private key file '%s' does not exist.\n", ".encrypt.priv")
		return false
	}
	_, err = os.Stat("./.encrypt.pub")
	if os.IsNotExist(err) {
		fmt.Printf("Public key file '%s' does not exist.\n", ".encrypt.pub")
		return false
	}
	return true
}
