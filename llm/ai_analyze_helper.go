package llm

import (
	"car-system-go/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func AIReportUser(records []model.InfractionRecord) (AnalysisResult, error) {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		return AnalysisResult{}, fmt.Errorf("加载环境变量失败: %v", err)
	}

	// 检查环境变量是否存在
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")

	if apiKey == "" || modelName == "" {
		return AnalysisResult{}, fmt.Errorf("请设置环境变量 ARK_API_KEY 和 MODEL")
	}

	aiModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("初始化模型失败: %v", err)
	}

	// 提示词模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是专业的汽车驾驶行为分析师，负责为司机提供数据支持和决策建议。
你的分析需专业、客观且具有实操性，性格俏皮灵动但不失严谨。
请聚焦于驾驶风险评估和管理措施建议，避免无关内容。
⚠️ 你必须严格按照给定的 JSON 格式输出，不能输出其他任何说明文字。
        `),
		&schema.Message{
			Role: schema.User,
			Content: `请基于以下信息分析驾驶员的驾驶行为并提供管理建议（仅给后台管理员参考）：

违规记录：{records}

请严格输出以下 JSON 格式：
{{
  "code": 200,
  "message": "success",
  "data": [
    "请避免急刹车，保持安全车距",
    "建议在施工路段减速慢行",
    "关注疲劳驾驶风险，定期休息"
  ]
}}

`,
		},
	)

	params := map[string]any{
		"records": records,
	}

	messages, err := template.Format(ctx, params)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("格式化提示词失败: %v", err)
	}

	// 调用模型生成回复
	response, err := aiModel.Generate(ctx, messages)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("生成回复失败: %v", err)
	}

	// 解析 JSON
	var result AnalysisResult
	if err := json.Unmarshal([]byte(response.Content), &result); err != nil {
		return AnalysisResult{}, fmt.Errorf("解析 AI 输出失败: %v\n原始输出: %s", err, response.Content)
	}

	return result, nil
}

func AIModels(user model.User, records []model.InfractionRecord) (string, error) {
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

func AIReport(records []model.InfractionRecord) (ReportResult, error) {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		return ReportResult{}, fmt.Errorf("加载环境变量失败: %v", err)
	}

	// 检查环境变量是否存在
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")

	if apiKey == "" || modelName == "" {
		return ReportResult{}, fmt.Errorf("请设置环境变量 ARK_API_KEY 和 MODEL")
	}

	aiModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		return ReportResult{}, fmt.Errorf("初始化模型失败: %v", err)
	}

	// 提示词模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是专业的汽车驾驶行为分析师，负责为后台管理员提供数据支持和决策建议。
你的分析需专业、客观且具有实操性，性格俏皮灵动但不失严谨。
请聚焦于驾驶风险评估和管理措施建议，避免无关内容。

只包含json的格式仅返回纯JSON字符串，不包含任何额外说明文字、注释或格式标记（不要使用markdown标签，不要为我生成特殊的符号，就是一行json字符串）

⚠️ 你必须严格按照给定的 JSON 格式输出，不能输出其他任何说明文字，
        `),
		&schema.Message{
			Role: schema.User,
			Content: `请基于以下信息分析驾驶员的驾驶行为并提供管理建议（仅给后台管理员参考）：
违规记录：{records}
你是一个数据处理分析大师,应当给出所有司机的共性问题分析：如返回数据：（如“本周急刹车事件显著增加,某个司机在近期的违规率显著增加等”）
请严格输出以下 JSON 格式：
          {{
          "code": 200,
          "message": "success",
          "data": {{
            "coreConclusion": "⚠️ 本周急刹车事件显著增加",(控制在20个字以内)  // 核心结论
            "reason": "通过对数据的分析发现，急刹车事件同比上周增加30%，其中85%发生在XX路南向北方向。这与该路段本周开始的围挡施工时间高度吻合。",  // 根因解释
            "advices": "相关操作建议xxxx"//建议
			"trustIndex:"报告建议的可信度，范围到60%～100%"
              }}
        }}
`,
		},
	)

	params := map[string]any{
		"records": records,
	}

	messages, err := template.Format(ctx, params)
	if err != nil {
		return ReportResult{}, fmt.Errorf("格式化提示词失败: %v", err)
	}

	// 调用模型生成回复
	response, err := aiModel.Generate(ctx, messages)
	if err != nil {
		return ReportResult{}, fmt.Errorf("生成回复失败: %v", err)
	}

	// 解析 JSON
	var result ReportResult
	if err := json.Unmarshal([]byte(response.Content), &result); err != nil {
		return ReportResult{}, fmt.Errorf("解析 AI 输出失败: %v\n原始输出: %s", err, response.Content)
	}

	return result, nil
}

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
3.不要用markdown进行输出，只用输出文字就可以了；
4.只要文字，连** **
`,
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

// DriverClassify 驾驶员分类函数
func DriverClassify(records []model.InfractionRecord) (string, error) {
	ctx := context.Background()

	// 加载环境变量（只在需要时加载，避免重复加载）
	if err := godotenv.Load(); err != nil {
		// 不使用panic，而是返回错误
		return "", fmt.Errorf("加载环境变量失败: %w", err)
	}

	// 检查环境变量是否存在
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")

	if apiKey == "" {
		return "", fmt.Errorf("环境变量 ARK_API_KEY 未设置")
	}
	if modelName == "" {
		return "", fmt.Errorf("环境变量 MODEL 未设置")
	}

	// 初始化AI模型
	aiModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		return "", fmt.Errorf("初始化模型失败: %w", err)
	}

	// 构造更优化的提示词模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一名专业的驾驶行为分析专家。请基于给定的司机违规记录数据进行群体分类分析。

分类标准：
- 安全型：无违规或轻微违规记录的司机
- 激进型：有超速、危险驾驶等激进行为记录的司机  
- 疲劳型：有疲劳驾驶、注意力不集中等记录的司机
- 未知型：数据不足或无法明确分类的司机

请严格按照JSON格式输出，不要添加任何解释文字。
请分析以下司机违规记录数据并进行分类：

违规记录数据：{records}

请输出严格的JSON格式，必须包含以下字段：
{
  "totalDrivers": 数字,
  "categories": [
    { "type": "安全型", "count": 数字 },
    { "type": "激进型", "count": 数字 },
    { "type": "疲劳型", "count": 数字 },
    { "type": "未知型", "count": 数字 }
  ]
}

注意：只返回JSON，不要包含任何其他文字。


`),
		&schema.Message{
			Role: schema.User,
		},
	)

	// 准备模板参数
	params := map[string]any{
		"records": records,
	}

	// 格式化消息
	messages, err := template.Format(ctx, params)
	if err != nil {
		return "", fmt.Errorf("格式化消息失败: %w", err)
	}

	// 调用模型生成回复
	response, err := aiModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("AI生成回复失败: %w", err)
	}

	responseContent := response.Content

	// 清理可能的markdown代码块标记
	responseContent = strings.TrimPrefix(responseContent, "```json")
	responseContent = strings.TrimSuffix(responseContent, "```")
	responseContent = strings.TrimSpace(responseContent)
	return response.Content, nil
}
