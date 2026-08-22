package envelope

type DataEnvelopeDTO[T any] struct {
	Data T `json:"data"`
}
