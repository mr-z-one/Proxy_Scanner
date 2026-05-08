package Model

import (
	"fmt"
	"proxyScanner/dataType"

	"gorm.io/gorm"
)

type HttpRequest struct {
	CreatedAt int64  `gorm:"autoCreateTime;index"`
	ID        uint   `gorm:"primaryKey"`
	Signature string `gorm:"uniqueIndex;type:varchar(64);not null"`

	Path   string `gorm:"not null;index:idx_search;size:500"`
	Method string `gorm:"not null;index:idx_search;size:10"`

	RequestHeaders  dataType.JSONMap `gorm:"type:jsonb;not null"`
	RequestBodyRaw  string           `gorm:"type:text"`
	RequestBodyJson dataType.JSONMap `gorm:"type:jsonb"`

	StatusCode       int              `gorm:"index"`
	ResponseHeaders  dataType.JSONMap `gorm:"type:jsonb;not null"`
	ResponseBodyJson dataType.JSONMap `gorm:"type:jsonb"`
	ResponseBodyRaw  string           `gorm:"type:text"`
}

func (r *HttpRequest) BeforeCreate(tx *gorm.DB) error {

	if c, ok := r.ResponseHeaders["Content-Type"]; ok {
		cc, ok := c.(string)
		if ok {
			fmt.Println(cc)
		}
	}
	panic("shit")
	return nil
}
