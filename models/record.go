package models

type RecordModel struct {
	Id             string `db:"id" json:"id"`
	IdentityNumber int64  `db:"identityNumber" json:"identityNumber"`
	Symptoms       string `db:"symptoms" json:"symptoms"`
	Medications    string `db:"medications" json:"medications"`
	CreatedAt      string `db:"createdAt" json:"createdAt"`
	CreatorId      string `db:"creatorId" json:"creatorId"`
}
