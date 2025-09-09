package setup

import (
	"car-system-go/utils"
	"fmt"
)

func InitAvatar() {
	err := utils.CreateAvatarDirectory()
	if err != nil {
		fmt.Printf("创建头像目录失败: %v\n", err)
	}
}
