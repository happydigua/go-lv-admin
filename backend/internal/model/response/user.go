package response

import "go-lv-vue-admin/internal/model"

type LoginResponse struct {
	User      model.LvUser `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
}
