package models

type NurseModel struct {
    Id                 string  `db:"id" json:"id"`
    NIP                int64   `db:"nip" json:"nip"`
    Name               string  `db:"name" json:"name"`
    Password           string  `db:"password" json:"password"`
    IsGranted          bool    `db:"isGranted" json:"isGranted"`
    IdentityCardScanning string `db:"identityCardScanning" json:"identityCardScanning"`
    CreatedAt          string  `db:"createdAt" json:"createdAt"`
    UpdatedAt          string  `db:"updatedAt" json:"updatedAt"`
    DeletedAt          *string `db:"deletedAt" json:"deletedAt"`
}