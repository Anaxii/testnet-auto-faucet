package externaldatabase

type VerificationRequest struct {
	Name            string          `json:"name" bson:"name"`
	Status          string          `json:"status" bson:"status"`
	Email           string          `json:"email" bson:"email"`
	WalletAddress   string          `json:"wallet_address" bson:"wallet_address"`
	Country         string          `json:"country" bson:"country"`
	PhysicalAddress string          `json:"physical_address" bson:"physical_address"`
	IdentityNumber  string          `json:"identity_number" bson:"identity_number"`
	DateOfBirth     string          `json:"date_of_birth" bson:"date_of_birth"`
	Signature       SignatureStruct `json:"signature" bson:"signature"`
	Beneficiary     struct {
		Name          string `json:"name" bson:"name"`
		WalletAddress string `json:"wallet_address" bson:"wallet_address"`
	}
}

type SignatureStruct struct {
	Message       string        `json:"message" bson:"message"`
	Account       string        `json:"account" bson:"account"`
	SignatureData SignatureData `json:"signature_data" bson:"signature_data"`
}

type SignatureData struct {
	HashedMessage string `json:"hashed_message" bson:"hashed_message"`
	R             string `json:"r" bson:"r"`
	S             string `json:"s" bson:"s"`
	V             string `json:"v" bson:"v"`
	Sig           string `json:"sig" bson:"sig"`
}
