package models

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type PatientModel struct {
	Id                   string `db:"id" json:"id"`
	IdentityNumber       int64  `db:"identityNumber" json:"identityNumber"`
	Name                 string `db:"name" json:"name"`
	PhoneNumber          string `db:"phoneNumber" json:"phoneNumber"`
	BirthDate            string `db:"birthDate" json:"birthDate"`
	Gender               Gender `db:"gender" json:"gender"`
	IdentityCardScanImg string `db:"identityCardScanImg" json:"identityCardScanImg"`
	CreatedAt            string `db:"createdAt" json:"createdAt"`
}
