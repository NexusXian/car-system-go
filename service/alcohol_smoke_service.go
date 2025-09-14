package service

import (
	"car-system-go/request"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// 违规类型定义
const (
	ViolationNone             = iota
	ViolationDrunkDriving     // 酒驾
	ViolationDangerousDriving // 醉驾
	ViolationSmoking          // 抽烟
)

// 用户违规状态
type userViolationState struct {
	lastViolationType int       // 上次违规类型
	lastViolationTime time.Time // 上次违规时间
	alcoholLevel      float64   // 当前酒精水平
	smokeLevel        float64   // 当前烟雾水平
	violationCount    int       // 违规次数统计
}

var (
	// 全局用户状态映射
	userStates = make(map[string]*userViolationState)
	stateMutex = &sync.RWMutex{}

	// 违规判定阈值
	alcoholWarningThreshold = 120.0 // 酒驾阈值
	alcoholDangerThreshold  = 180.0 // 醉驾阈值
	smokeThreshold          = 90.0  // 抽烟阈值

	// 违规冷却时间（同一违规类型在此时间内不重复记录）
	violationCooldown = 30 * time.Minute
)

func AlcoholSmokeService() {
	var user request.UserInfractionCreateRequest
	user.IDCardNumber = "110101199001011234"
	user.RealName = "张三"
	user.LicensePlate = "京A12345"

	listenIP := ""       // 空字符串表示监听所有网络接口
	listenPort := "8080" // 与开发板配置的目标端口一致
	bufferSize := 1024   // 缓冲区大小，根据开发板发送的数据大小调整
	listenAddr := fmt.Sprintf("%s:%s", listenIP, listenPort)
	conn, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		fmt.Printf("无法创建UDP连接: %v\n", err)
		os.Exit(1)
	}
	defer func(conn net.PacketConn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	fmt.Printf("正在监听 UDP %s 端口，等待开发板数据...\n", listenPort)
	fmt.Println("请确保开发板配置的目标IP是本机IP，目标端口是", listenPort)
	fmt.Println("按 Ctrl+C 退出程序")

	buffer := make([]byte, bufferSize)

	sensorRegex := regexp.MustCompile(`Alcohol=([\d.]+),\s*Smoke=([\d.]+)`)

	for {
		n, remoteAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("接收数据错误: %v\n", err)
			continue
		}

		recvTime := time.Now().Format("2006-01-02 15:04:05")
		data := string(buffer[:n])

		fmt.Printf("\n[%s] 从 %s 接收到数据：\n", recvTime, remoteAddr.String())
		fmt.Printf("  文本内容: %s\n", data)

		// 使用正则表达式提取数据
		matches := sensorRegex.FindStringSubmatch(data)
		if len(matches) == 3 {
			// 将提取的字符串转换为浮点数
			alcohol, err1 := strconv.ParseFloat(matches[1], 64)
			smoke, err2 := strconv.ParseFloat(matches[2], 64)

			if err1 == nil && err2 == nil {
				// 获取或创建用户状态
				state := getUserState(user.IDCardNumber)

				// 更新传感器数据
				stateMutex.Lock()
				state.alcoholLevel = alcohol
				state.smokeLevel = smoke
				stateMutex.Unlock()

				// 检查并处理违规
				processViolations(user, state)
			} else {
				fmt.Printf("  数据转换错误: Alcohol=%v, Smoke=%v\n", err1, err2)
			}
		} else {
			fmt.Println("  未能从数据中提取传感器读数")
		}

		ack := fmt.Sprintf("已收到 %d 字节数据", n)
		_, err = conn.WriteTo([]byte(ack), remoteAddr)
		if err != nil {
			fmt.Printf("发送确认消息失败: %v\n", err)
		} else {
			fmt.Println("  已发送确认消息")
		}
	}
}

// 获取用户状态
func getUserState(userID string) *userViolationState {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	if state, exists := userStates[userID]; exists {
		return state
	}

	// 创建新用户状态
	state := &userViolationState{
		lastViolationType: ViolationNone,
		lastViolationTime: time.Now().Add(-violationCooldown), // 设置为冷却时间前，确保第一次检测能正常记录
	}
	userStates[userID] = state
	return state
}

// 处理违规检测
func processViolations(user request.UserInfractionCreateRequest, state *userViolationState) {
	now := time.Now()

	// 检查是否在违规冷却期内
	inCooldown := now.Sub(state.lastViolationTime) < violationCooldown

	// 确定当前最高优先级的违规类型
	currentViolation := determineViolationType(state.alcoholLevel, state.smokeLevel)

	// 如果没有违规，重置状态
	if currentViolation == ViolationNone {
		if state.lastViolationType != ViolationNone {
			fmt.Println("  违规状态已解除")
			stateMutex.Lock()
			state.lastViolationType = ViolationNone
			stateMutex.Unlock()
		}
		return
	}

	// 如果有违规，但类型相同且在冷却期内，不重复记录
	if inCooldown && currentViolation == state.lastViolationType {
		fmt.Printf("  检测到违规但仍在冷却期内 (剩余: %v)\n",
			violationCooldown-now.Sub(state.lastViolationTime))
		return
	}

	// 如果违规类型升级（如从酒驾升级为醉驾），即使冷却期未过也记录
	shouldRecord := !inCooldown || currentViolation > state.lastViolationType

	if shouldRecord {
		// 记录违规
		switch currentViolation {
		case ViolationDrunkDriving:
			user.Record = "开车酒驾"
		case ViolationDangerousDriving:
			user.Record = "开车醉驾"
		case ViolationSmoking:
			user.Record = "开车抽烟,存在单手掌握方向盘"
		default:
			panic("unhandled default case")
		}

		if err := UserInfractionCreateService(user); err != nil {
			fmt.Printf("记录违规失败: %v\n", err)
		} else {
			stateMutex.Lock()
			state.lastViolationType = currentViolation
			state.lastViolationTime = now
			state.violationCount++
			stateMutex.Unlock()

			fmt.Printf("  已记录违规: %s (总违规次数: %d)\n", user.Record, state.violationCount)
		}
	}
}

// 确定违规类型
func determineViolationType(alcohol, smoke float64) int {
	if alcohol >= alcoholDangerThreshold {
		return ViolationDangerousDriving
	} else if alcohol >= alcoholWarningThreshold {
		return ViolationDrunkDriving
	} else if smoke >= smokeThreshold {
		return ViolationSmoking
	}
	return ViolationNone
}
