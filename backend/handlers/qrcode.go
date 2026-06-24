package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

var verifyPIN string

func InitQR() {
	verifyPIN = os.Getenv("VERIFY_PIN")
	if verifyPIN == "" {
		verifyPIN = "xinhang2026"
	}
}

func deriveKey() []byte {
	h := sha256.Sum256(jwtSecret)
	return h[:]
}

func encryptPayload(plaintext string) (string, error) {
	block, err := aes.NewCipher(deriveKey())
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptPayload(encoded string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(deriveKey())
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	plaintext, err := aesGCM.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func GeneratePermitQR(c *gin.Context) {
	appNo := c.Query("no")
	student := c.Query("student")
	if appNo == "" || student == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing no or student"})
		return
	}

	plaintext := fmt.Sprintf("XINHANG|%s|%s", appNo, student)
	encrypted, err := encryptPayload(plaintext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
		return
	}

	png, err := qrcode.Encode(encrypted, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "qr generation failed"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

func VerifyPermitQR(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		PIN     string `json:"pin" binding:"required"`
	}
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content and pin required"})
		return
	}

	if req.PIN != verifyPIN {
		c.JSON(http.StatusForbidden, gin.H{"valid": false, "message": "验证密码错误"})
		return
	}

	doVerify(c, req.Content)
}

func VerifyPermitQRGet(c *gin.Context) {
	code := c.Query("code")
	pin := c.Query("pin")
	if code == "" || pin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code and pin required"})
		return
	}
	if pin != verifyPIN {
		c.JSON(http.StatusForbidden, gin.H{"valid": false, "message": "验证密码错误"})
		return
	}
	doVerify(c, code)
}

func doVerify(c *gin.Context, encrypted string) {
	plaintext, err := decryptPayload(encrypted)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": "无法解密，非本校签发的二维码"})
		return
	}

	parts := strings.SplitN(plaintext, "|", 3)
	if len(parts) < 3 || parts[0] != "XINHANG" {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": "二维码格式无效"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "准考证验证通过",
		"appNo":   parts[1],
		"student": parts[2],
	})
}
