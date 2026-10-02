package local_conf

type Jwt struct {
	Iss             string   `bson:"iss"`
	Aud             string   `bson:"aud"`
	ExpiryInSeconds int64    `bson:"expiryInSeconds"`
	Scopes          []string `bson:"scopes"`
	PrivateKey      string   `bson:"privateKey"`
	PublicKey       string   `bson:"publicKey"`
}
