package domain

import "time"

type Short_Link struct {
	ID           int        `json:"id"`
	Short_code   string     `json:"short_code"`
	Original_URL string     `json:"original_url"`
	User_ID      int        `json:"user_id"`
	Created_at   time.Time  `json:"created_at"`
	Expires_at   *time.Time `json:"expires_at,omitempty"`
}
