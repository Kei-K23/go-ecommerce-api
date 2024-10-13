package dto

import "time"

type ReviewRequest struct {
	ID        int       `json:"id"`
	ProductId int       `json:"product_id"`
	UserId    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
