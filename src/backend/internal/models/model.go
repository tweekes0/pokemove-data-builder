package models

type Model interface {
	BulkInsert([]interface{}) error
	RelationsBulkInsert([]interface{}) error
}