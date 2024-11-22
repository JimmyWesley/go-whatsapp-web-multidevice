package send

type ChatPresenceRequest struct {
	Phone    string `json:"phone"`
	Presence string `json:"presence"`        // Ex: "composing", "paused"
	Media    string `json:"media,omitempty"` // Opcional: Ex: "audio", ""
}
