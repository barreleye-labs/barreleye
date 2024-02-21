package dto

type Signer struct {
	X string `json:"x"`
	Y string `json:"y"`
}

func CreateSigner(x string, y string) Signer {
	return Signer{
		X: x,
		Y: y,
	}
}
