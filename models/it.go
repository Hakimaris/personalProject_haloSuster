package models

type ItModel struct {
	Id        string `db:"id" json:"id"`
	NIP       int64  `db:"nip" json:"nip"`
	Name      string `db:"name" json:"name"`
	Password  string `db:"password" json:"password"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
}
