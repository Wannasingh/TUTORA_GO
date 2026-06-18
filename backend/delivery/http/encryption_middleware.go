package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/config"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

// RequestPayload represents the incoming encrypted format
type RequestPayload struct {
	Data string `json:"data"`
}

// ResponsePayload represents the outgoing encrypted format
type ResponsePayload struct {
	Data string `json:"data"`
}

// DecryptionMiddleware decrypts incoming encrypted JSON request bodies
func DecryptionMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read raw body bytes
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}

		// If body is empty, proceed without decryption
		if len(bodyBytes) == 0 {
			c.Next()
			return
		}

		// Parse the wrapper JSON
		var payload RequestPayload
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload format, expected {data: ...}"})
			c.Abort()
			return
		}

		// Decrypt the ciphertext data
		decryptedBytes, err := utils.DecryptAES(payload.Data, []byte(cfg.PayloadEncryptionKey))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "decryption failed"})
			c.Abort()
			return
		}

		// Restore decrypted body into request
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBytes))
		c.Request.ContentLength = int64(len(decryptedBytes))

		c.Next()
	}
}

// encryptionWriter is a custom writer that intercepts Gin writes to encrypt payloads before sending
type encryptionWriter struct {
	gin.ResponseWriter
	bodyBuffer *bytes.Buffer
	key        []byte
}

func (w *encryptionWriter) Write(b []byte) (int, error) {
	return w.bodyBuffer.Write(b)
}

func (w *encryptionWriter) WriteString(s string) (int, error) {
	return w.bodyBuffer.WriteString(s)
}

// EncryptionMiddleware encrypts outgoing JSON responses using AES-GCM
func EncryptionMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := &bytes.Buffer{}
		writer := &encryptionWriter{
			ResponseWriter: c.Writer,
			bodyBuffer:     buffer,
			key:            []byte(cfg.PayloadEncryptionKey),
		}
		c.Writer = writer

		c.Next()

		// If request is aborted or not successful status (e.g. 4xx, 5xx), we might skip encrypting details,
		// but to be standard, we encrypt all JSON responses.
		if buffer.Len() == 0 {
			return
		}

		responseBytes := buffer.Bytes()

		// Encrypt response
		ciphertext, err := utils.EncryptAES(responseBytes, writer.key)
		if err != nil {
			// Fallback to error response if encryption fails
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(`{"error":"failed to encrypt response"}`))
			return
		}

		// Format output as wrapper JSON
		wrapper := ResponsePayload{Data: ciphertext}
		wrapperBytes, err := json.Marshal(wrapper)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(`{"error":"failed to marshal response"}`))
			return
		}

		// Write to actual client response
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(wrapperBytes)
	}
}
