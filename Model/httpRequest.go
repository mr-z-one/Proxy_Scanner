package Model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"proxyScanner/dataType"
	"strings"

	"gorm.io/gorm"
)

type HttpRequest struct {
	CreatedAt int64  `gorm:"autoCreateTime;index"`
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Signature string `gorm:"unique;type:varchar(128);not null"`

	Host string `gorm:"type:varchar(150);not null"`

	Path   string `gorm:"not null;index:idx_search;size:500"`
	Method string `gorm:"not null;index:idx_search;size:10"`

	Parameters dataType.JSONMap `gorm:"type:jsonb"`

	RequestHeaders  dataType.JSONMap `gorm:"type:jsonb;not null"`
	RequestBodyRaw  string           `gorm:"type:text"`
	RequestBodyJson dataType.JSONMap `gorm:"type:jsonb"`

	StatusCode       int              `gorm:"index"`
	ResponseHeaders  dataType.JSONMap `gorm:"type:jsonb;not null"`
	ResponseBodyJson dataType.JSONMap `gorm:"type:jsonb"`
	ResponseBodyRaw  string           `gorm:"type:text"`
}

func (r *HttpRequest) BeforeCreate(tx *gorm.DB) error {

	err := HandleBinarytData(r)
	if err != nil {
		return err
	}
	r.generateSignature()

	HandleContentType(r)
	//fmt.Println("-----------", r.Signature)
	return nil
}

func HandleBinarytData(r *HttpRequest) error {
	binaryPrefixes := []string{
		"image/",
		"video/",
		"audio/",
		"font/",
		"application/vnd.",
		"application/x-",
		"application/zip",
		"application/pdf",
		"application/msword",
		"application/octet-stream",
		"application/java-archive",
		"application/x-gzip",
		"application/x-tar",
		"application/x-7z-compressed",
		"application/x-rar",
	}
	contentTypes := r.ResponseHeaders.GetKeyValue("Content-Type")
	for _, prefix := range binaryPrefixes {
		if strings.HasPrefix(contentTypes, prefix) {
			return fmt.Errorf("binary data already contains %s", contentTypes)
		}
	}
	return nil
}

func HandleContentType(r *HttpRequest) {
	contentTypes := r.ResponseHeaders.GetKeyValue("Content-Type")
	if contentTypes != "" {
		contentType := strings.TrimSpace(strings.Split(contentTypes, ";")[0])

		if contentType == "application/json" {
			jsonData := dataType.JSONMap{}
			err := json.Unmarshal([]byte(r.RequestBodyRaw), &jsonData)
			if err == nil {
				r.RequestBodyJson = jsonData

			}

			jsonData = dataType.JSONMap{}

			err = json.Unmarshal([]byte(r.ResponseBodyRaw), &jsonData)

			if err == nil {
				r.ResponseBodyJson = jsonData

			}

		}

	}
}
func (r *HttpRequest) generateSignature() {

	input := fmt.Sprintf("%s.%s.%s.%s", r.Path, r.Host, r.Method, r.hash(r.RequestBodyRaw))
	r.Signature = r.hash(input)
}
func (r *HttpRequest) hash(text string) string {
	if text == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
