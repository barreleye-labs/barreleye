package dto

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
}

func CreateSignature(r string, s string) Signature {
	return Signature{
		R: r,
		S: s,
	}
}
