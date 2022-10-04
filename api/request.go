package api

type PageReq struct {
	Page int `json:"page" validate:"required,min=1"`
	Size int `json:"size" validate:"required,min=1"`
}

type UserReq struct {
	Username string `json:"username" validate:"required,min=1,max=20"`
	Password string `json:"password" validate:"required,min=1,max=60"`
}

type UserDelReq struct {
	UserID int `json:"userid" validate:"required,min=2"`
}

type UserChgReq struct {
	UserID   int    `json:"userid" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1,max=60"`
}
