package models

type DataResponseModel struct {
	ResponseModel
	Data []Todo `json:"data"`
}

type SignleDataResponseModel struct {
	ResponseModel
	Data Todo `json:"data"`
}

type ResponseModel struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
