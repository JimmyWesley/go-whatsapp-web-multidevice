package send

import "mime/multipart"

type PTTRequest struct {
	Phone string                `json:"phone"` // Número do destinatário
	Audio *multipart.FileHeader `json:"Audio" form:"Audio"`
}
