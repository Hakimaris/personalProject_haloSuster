package models

type UserModel struct {
	ID                  string `json:"id" db:"id"`
	NIP                 int64     `json:"nip" db:"nip"`
	Name                string    `json:"name" db:"name"`
	Password            string    `json:"password" db:"password"`
	IdentityCardScanning string   `json:"identityCardScanImg" db:"identityCardScanning"`
	CreatedAt           string `json:"created_at" db:"createdAt"`
	UpdatedAt           string `json:"updated_at" db:"updatedAt"`
}