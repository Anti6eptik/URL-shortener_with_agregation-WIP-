package domain

import "time"

type Click_Event struct {
	ID            int       `json:"id"`
	Short_link_ID int       `json:"short_link_id"`
	Clicked_at    time.Time `json:"clicked_at"`
	Referer       string    `json:"referer"`
}
