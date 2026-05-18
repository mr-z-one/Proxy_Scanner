package Model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"proxyScanner/dataType"
	"proxyScanner/db"
	"strings"

	"gorm.io/gorm"
)

type HttpRequest struct {
	CreatedAt int64  `gorm:"autoCreateTime;index" json:"created_at"`
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Signature string `gorm:"unique;type:varchar(128);not null"json:"signature"`

	Host string `gorm:"type:varchar(150);not null"json:"host"`

	Path   string `gorm:"not null;index:idx_search;size:500"json:"path"`
	Method string `gorm:"not null;index:idx_search;size:10"json:"method"`

	Parameters dataType.JSONMap `gorm:"type:jsonb;index:,type:gin" json:"parameters"`

	RequestHeaders  dataType.JSONMap `gorm:"type:jsonb;not null;index:,type:gin" json:"requestHeaders"`
	RequestBodyRaw  string           `gorm:"type:text"json:"requestBodyRaw"`
	RequestBodyJson dataType.JSONMap `gorm:"type:jsonb;index:,type:gin"json:"requestBodyJson"`

	StatusCode       int              `gorm:"index"json:"statusCode"`
	ResponseHeaders  dataType.JSONMap `gorm:"type:jsonb;not null;index:,type:gin" json:"responseHeaders"`
	ResponseBodyJson dataType.JSONMap `gorm:"type:jsonb;index:,type:gin" json:"responseBodyJson"`
	ResponseBodyRaw  string           `gorm:"type:text"json:"responseBodyRaw"`
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
		"binary/octet-stream",
	}
	contentTypes := r.ResponseHeaders.GetKeyValue("Content-Type")
	lowerContentType := strings.ToLower(contentTypes)
	for _, prefix := range binaryPrefixes {
		if strings.HasPrefix(contentTypes, prefix) || strings.HasPrefix(lowerContentType, prefix) {
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

func GetRequestById(id uint) *HttpRequest {
	database := db.GetActiveDatabaseSession()
	result := HttpRequest{}
	database.Where("id = ?", id).First(&result)
	return &result
}
