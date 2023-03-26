package domain

type ReservationJson struct {
	State string `json:"state"`
	Order struct {
		ID    string `json:"id"`
		State string `json:"state"`
	} `json:"order"`
}
