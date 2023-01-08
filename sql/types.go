package sql

type BaseEntity struct {
	ID string `json:"id"`
}

type BaseRelationship struct {
	BaseEntity
	In  string `json:"in"`
	Out string `json:"out"`
}
