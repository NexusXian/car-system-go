package llm

import (
	"car-system-go/model"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func AIModel(user model.User, records []model.InfractionRecord) (string, error) {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		panic(err)
		return "", err
	}

	// 检查环境变量是否存在
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")

	if apiKey == "" || modelName == "" {
		panic("请设置环境变量 ARK_API_KEY 和 MODEL")
	}

	aiModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		panic(fmt.Sprintf("初始化模型失败: %v", err))
		return "", err
	}

	// 优化后的提示词模板：更明确的分析目标和输出结构
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是专业的汽车驾驶行为分析师，负责为后台管理员提供数据支持和决策建议。
你的分析需专业、客观且具有实操性，性格俏皮灵动但不失严谨。
请聚焦于驾驶风险评估和管理措施建议，避免无关内容。`),
		&schema.Message{
			Role: schema.User,
			Content: `请基于以下信息分析驾驶员的驾驶行为并提供管理建议（仅给后台管理员参考）：
1. 驾驶员信息：{user}
2. 最近三次违规记录：{records}

请按以下结构输出：
- 违规特征分析：总结违规类型、频率、时间/地点规律等关键特征
- 风险等级评估：基于违规情况评定风险等级（低/中/高）并说明理由
- 管理建议：提出具体可执行的措施（如重点监控、违规预警、培训干预等）
- 注意事项：提醒管理员需要关注的特殊情况或潜在问题

分析需简洁明了，重点突出对管理工作的实际指导价值。`,
		},
	)

	params := map[string]any{
		"user":    user,
		"records": records,
	}

	messages, err := template.Format(ctx, params)
	if err != nil {
		panic(fmt.Sprintf("格式化提示词失败: %v", err))
		return "", err
	}

	// 调用模型生成回复
	response, err := aiModel.Generate(ctx, messages)
	if err != nil {
		panic(fmt.Sprintf("生成回复失败: %v", err))
		return "", err
	}
	return response.Content, nil
}

// 新增流式输出方法
func AIModelStream(user model.User, records []model.InfractionRecord) (<-chan string, error) {
	ctx := context.Background()

	// 加载环境变量（避免 panic，改为返回错误）
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("加载环境变量失败: %v", err)
	}

	// 检查环境变量
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")
	if apiKey == "" || modelName == "" {
		return nil, fmt.Errorf("请设置环境变量 ARK_API_KEY 和 MODEL")
	}

	// 初始化模型
	aiModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化模型失败: %v", err)
	}

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是专业的汽车驾驶行为分析师，负责为后台管理员提供数据支持和决策建议。
你的分析需专业、客观且具有实操性，性格俏皮灵动但不失严谨。
请聚焦于驾驶风险评估和管理措施建议，避免无关内容。`),
		&schema.Message{
			Role: schema.User,
			Content: `请基于以下信息分析驾驶员的驾驶行为并提供管理建议（仅给后台管理员参考）：
1. 驾驶员信息：{user}
2. 最近三次违规记录：{records}

请按以下结构输出：
- 违规特征分析：总结违规类型、频率、时间/地点规律等关键特征
- 风险等级评估：基于违规情况评定风险等级（低/中/高）并说明理由
- 管理建议：提出具体可执行的措施（如重点监控、违规预警、培训干预等）
- 注意事项：提醒管理员需要关注的特殊情况或潜在问题

分析需简洁明了，重点突出对管理工作的实际指导价值。`,
		},
	)

	params := map[string]any{
		"user":    user,
		"records": records,
	}

	messages, err := template.Format(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("格式化提示词失败: %v", err)
	}

	// 调用模型的流式接口
	streamReader, err := aiModel.Stream(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("开启流式输出失败: %v", err)
	}

	// 创建通道用于传递流式片段
	streamChan := make(chan string)

	// 启动 goroutine 读取流并发送到通道
	// 启动 goroutine 读取流并发送到通道
	go func() {
		defer close(streamChan) // 确保通道最终关闭

		// 循环读取流中的内容
		for {
			// 读取下一个流式片段
			msg, err := streamReader.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, schema.ErrNoValue) {
					break // 直接退出循环，不发送任何内容
				}
				// 只有非预期错误（如网络中断、API 异常）才发送错误提示
				streamChan <- fmt.Sprintf("流式输出错误: %v", err)
				break
			}

			// 将片段内容发送到通道
			if msg.Content != "" {
				streamChan <- msg.Content
			}
		}
	}()

	return streamChan, nil
}

// 根据全部数据为管理员提供问答服务
func AIQuestionStream(records []model.InfractionRecord, message string) (<-chan string, error) {
	ctx := context.Background()

	// 加载环境变量（避免 panic，改为返回错误）
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("加载环境变量失败: %v", err)
	}

	// 检查环境变量
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")
	if apiKey == "" || modelName == "" {
		return nil, fmt.Errorf("请设置环境变量 ARK_API_KEY 和 MODEL")
	}

	// 初始化模型
	aiModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化模型失败: %v", err)
	}

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是专业的汽车驾驶行为分析师，负责为后台管理员提供数据支持和决策建议。
你的分析需专业、客观且具有实操性，性格俏皮灵动但不失严谨。
请聚焦于驾驶风险评估和管理措施建议，避免无关内容。`),
		&schema.Message{
			Role: schema.User,
			Content: `请基于以下信息分析驾驶员的驾驶行为并提供管理建议（仅给后台管理员参考）：
1. 所有的违规记录：{records}
2. 管理员的问题: {questions}

请按以下结构输出：
- 违规特征分析：总结违规类型、频率、时间/地点规律等关键特征
- 风险等级评估：基于违规情况评定风险等级（低/中/高）并说明理由
- 管理建议：提出具体可执行的措施（如重点监控、违规预警、培训干预等）
- 注意事项：提醒管理员需要关注的特殊情况或潜在问题

分析需简洁明了，重点突出对管理工作的实际指导价值。`,
		},
	)

	params := map[string]any{
		"records":   records,
		"questions": message,
	}

	messages, err := template.Format(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("格式化提示词失败: %v", err)
	}

	// 调用模型的流式接口
	streamReader, err := aiModel.Stream(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("开启流式输出失败: %v", err)
	}

	// 创建通道用于传递流式片段
	streamChan := make(chan string)

	// 启动 goroutine 读取流并发送到通道
	// 启动 goroutine 读取流并发送到通道
	go func() {
		defer close(streamChan) // 确保通道最终关闭

		// 循环读取流中的内容
		for {
			// 读取下一个流式片段
			msg, err := streamReader.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, schema.ErrNoValue) {
					break // 直接退出循环，不发送任何内容
				}
				// 只有非预期错误（如网络中断、API 异常）才发送错误提示
				streamChan <- fmt.Sprintf("流式输出错误: %v", err)
				break
			}

			// 将片段内容发送到通道
			if msg.Content != "" {
				streamChan <- msg.Content
			}
		}
	}()

	return streamChan, nil
}
