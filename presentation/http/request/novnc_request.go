package request

type NoVNCRequestForm struct {
	Address string `form:"address" binding:"required"`
	Port    string `form:"port"  binding:"required"`
}
