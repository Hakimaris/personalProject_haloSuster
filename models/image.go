package models

type ImageModel struct {
    Id        string `db:"id" json:"id"`
    Path      string `db:"path" json:"path"`
    CreatedAt string `db:"createdAt" json:"createdAt"`
}
