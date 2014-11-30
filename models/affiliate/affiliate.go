package affiliate

import (
	"time"
)

type Link struct {
	Id          string    `json:"id"`
	Namespace   string    `json:"namespace"`
	Description string    `json:"description"`
	Code        []byte    `json:"code"`
	CreatedAt   time.Time `json:"updated_at"`
	Link        string    `json:"link"`
}
