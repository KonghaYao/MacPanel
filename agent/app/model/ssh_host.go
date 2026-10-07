package model

const (
	SSHHostAuthAuthorized = "authorized"
	SSHHostAuthFailed     = "failed"
)

type SSHHost struct {
	BaseModel
	Alias      string `json:"alias" gorm:"uniqueIndex;not null;size:64"`
	HostName   string `json:"hostName" gorm:"not null;size:255"`
	User       string `json:"user" gorm:"not null;size:64"`
	Port       int    `json:"port" gorm:"not null;default:22"`
	AuthStatus string `json:"authStatus" gorm:"not null;size:32;default:authorized"`
}
