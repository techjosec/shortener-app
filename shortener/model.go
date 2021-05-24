package shortener

type Redirect struct {
	Code      string `json:"code" bson:"code" mspack:"code"`
	URL       string `json:"url" bson:"url" mspack:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" mspack:"created_at"`
}
