package entities

type PharmacyShippingMethod struct {
	Id                    int64
	PharmacyId            int64
	OfficialShippingId    int64
	NonOfficialShippingId int64
}
