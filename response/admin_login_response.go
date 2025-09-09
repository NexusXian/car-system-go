package response

import "time"

type AdminLoginResponse struct {
	Token      string    `json:"token"`
	ExpireTime time.Time `json:"expireTime"`
	AdminInfo  AdminInfo `json:"adminInfo"`
	TimeStamp  int64     `json:"timeStamp"`
}

type AdminInfo struct {
	WorkID    string `json:"workID"`
	AvatarUrl string `json:"avatarUrl"`
}
