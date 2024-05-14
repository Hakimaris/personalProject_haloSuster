package models

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type PatientModel struct {
	Id                   string `db:"id" json:"id"`
	IdentityNumber       int64  `db:"identityNumber" json:"identityNumber"`
	Name                 string `db:"name" json:"name"`
	PhoneNumber          string `db:"phoneNumber" json:"phoneNumber"`
	BirthDate            string `db:"birthDate" json:"birthDate"`
	Gender               Gender `db:"gender" json:"gender"`
	IdentityCardScanning string `db:"identityCardScanning" json:"identityCardScanning"`
	CreatedAt            string `db:"createdAt" json:"createdAt"`
}