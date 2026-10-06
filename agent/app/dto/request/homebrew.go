package request

type HomebrewSearchReq struct {
	Query string `json:"query" validate:"required,min=1"`
	Type  string `json:"type" validate:"omitempty,oneof=formula cask all"`
}

type HomebrewListReq struct {
	Page     int    `json:"page" form:"page" validate:"required,min=1"`
	PageSize int    `json:"pageSize" form:"pageSize" validate:"required,min=1,max=500"`
	Info     string `json:"info" form:"info"`
	Type     string `json:"type" form:"type" validate:"omitempty,oneof=formula cask all"`
}

type HomebrewPackageReq struct {
	Name   string `json:"name"`
	Type   string `json:"type" validate:"omitempty,oneof=formula cask"`
	TaskID string `json:"taskID"`
}

type HomebrewMirrorUpdateReq struct {
	BottleDomain  string `json:"bottleDomain"`
	APIDomain     string `json:"apiDomain"`
	BrewGitRemote string `json:"brewGitRemote"`
	CoreGitRemote string `json:"coreGitRemote"`
	CaskGitRemote string `json:"caskGitRemote"`
}

type HomebrewTaskReq struct {
	TaskID string `json:"taskID"`
}
