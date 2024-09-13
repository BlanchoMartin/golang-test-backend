// models/message.go
package models

type Brainee struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    Text   string `json:"text"`
    Author string `json:"author"`
    Brand  string `json:"brand"`
}
