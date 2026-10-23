package auth

// SignInWithPasscodeRequest represents the expected payload for a user sign-in with passcode request.
type SignInWithPasscodeRequest struct {
	ID       string `json:"id"       binding:"required,lte=50"       example:"000027b0e9a211ef95a30242ac180002"`
	Passcode string `json:"passcode" binding:"required,number,len=6" example:"123456"`
}
