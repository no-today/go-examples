package domain

import "time"

type Model struct {
	CreatedDate      time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	LastModifiedDate time.Time `json:"last_modified_date,omitempty" bson:"last_modified_date,omitempty"`
}

type AuditModel struct {
	CreatedBy      string `json:"created_by,omitempty" bson:"created_by,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty" bson:"last_modified_by,omitempty"`
}
