package request

type MirrorApplyReq struct {
	Ecosystem string            `json:"ecosystem" validate:"required"`
	PresetID  string            `json:"presetId"`
	Values    map[string]string `json:"values"`
}
