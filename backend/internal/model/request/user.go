package request

// User Login Structure
type Login struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`   // Verification code
	CaptchaId string `json:"captchaId"` // Verification code ID
}
