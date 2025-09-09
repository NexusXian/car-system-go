package llm

import (
	"car-system-go/model"
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
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
	// 1. 检查环境变量是否存在
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("MODEL")
	fmt.Println(modelName)
	fmt.Println(apiKey)
	if apiKey == "" || modelName == "" {
		panic("请设置环境变量 ARK_API_KEY 和 MODEL")
	}

	aiModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: os.Getenv("BASEURL"),
	})
	if err != nil {
		panic(fmt.Sprintf("初始化模型失败: %v", err))
		return "", err
	}

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role},你的性格是{character}"),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我分析一下这位驾驶员{user}的驾驶行为,他/她最近三次的违规记录在这里{records}，并给出相关建议,这个建议是给后台管理员看的，不是对驾驶员说的",
		},
	) // 3. 构造输入消息

	params := map[string]any{
		"role":      "汽车驾驶行为分析师,小助手",
		"user":      user,
		"character": "俏皮灵动，专业",
		"records":   records,
	}

	messages, err := template.Format(ctx, params)

	// 4. 调用模型生成回复
	response, err := aiModel.Generate(ctx, messages)
	if err != nil {
		panic(fmt.Sprintf("生成回复失败: %v", err))
		return "", err
	}
	return response.Content, nil
}
