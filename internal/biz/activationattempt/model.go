package activationattempt

import "time"

type ActivationAttempt struct {
	ID           uint      `gorm:"primary_key" json:"id" structs:"id"`
	CardValue    string    `gorm:"type:varchar(255);not null" json:"card_value" structs:"card_value"`
	ActivationAt time.Time `gorm:"type:datetime;not null" json:"activation_at" structs:"activation_at"`
	Success      bool      `gorm:"type:tinyint(1);not null" json:"success" structs:"success"`
	ErrorMessage string    `gorm:"type:varchar(255)" json:"error_message" structs:"error_message"`
	RequestData  string    `gorm:"type:json" json:"request_data" structs:"request_data"`
	ResponseData string    `gorm:"type:json" json:"response_data" structs:"response_data"`
	CreatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at" structs:"created_at"`
}

func (a *ActivationAttempt) TableName() string {
	return "activation_attempts"
}
