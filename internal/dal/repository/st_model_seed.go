package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/ch3nnn/webstack-go/internal/dal/model"
)

// seedDefaultModels 仅在模型表为空时，内置覆盖主流大模型的模型字典。
// 分类：GPT(OpenAI) / Claude(Anthropic) / Gemini(Google) / GLM(智谱) / Qwen(通义千问) /
//       DeepSeek / 豆包(字节) / Kimi(月之暗面) / 文心(百度) / 混元(腾讯) / 讯飞星火 /
//       Llama / Mistral / Gemma / Grok / 开源通用(通义万相/Stable Diffusion 等) / 其他主流
func seedDefaultModels(db *gorm.DB) {
	var count int64
	db.Model(&model.StModel{}).Count(&count)
	if count > 0 {
		return
	}

	now := time.Now()
	models := []model.StModel{
		// ===== OpenAI / GPT 系列 =====
		{Name: "GPT-4o", Provider: "OpenAI", Sort: 10, IsUsed: true},
		{Name: "GPT-4o mini", Provider: "OpenAI", Sort: 11, IsUsed: true},
		{Name: "GPT-4.1", Provider: "OpenAI", Sort: 12, IsUsed: true},
		{Name: "GPT-4.1 mini", Provider: "OpenAI", Sort: 13, IsUsed: true},
		{Name: "GPT-4 Turbo", Provider: "OpenAI", Sort: 14, IsUsed: true},
		{Name: "GPT-4", Provider: "OpenAI", Sort: 15, IsUsed: true},
		{Name: "GPT-3.5 Turbo", Provider: "OpenAI", Sort: 16, IsUsed: true},
		{Name: "o1", Provider: "OpenAI", Sort: 17, IsUsed: true},
		{Name: "o1-mini", Provider: "OpenAI", Sort: 18, IsUsed: true},
		{Name: "o1-preview", Provider: "OpenAI", Sort: 19, IsUsed: true},
		{Name: "o3", Provider: "OpenAI", Sort: 20, IsUsed: true},
		{Name: "o3-mini", Provider: "OpenAI", Sort: 21, IsUsed: true},
		{Name: "o4-mini", Provider: "OpenAI", Sort: 22, IsUsed: true},
		{Name: "GPT-5", Provider: "OpenAI", Sort: 23, IsUsed: true},
		{Name: "GPT-4o-2024-11-20", Provider: "OpenAI", Sort: 24, IsUsed: false},

		// ===== Anthropic / Claude 系列 =====
		{Name: "Claude 3.7 Sonnet", Provider: "Anthropic", Sort: 30, IsUsed: true},
		{Name: "Claude 3.5 Sonnet", Provider: "Anthropic", Sort: 31, IsUsed: true},
		{Name: "Claude 3.5 Haiku", Provider: "Anthropic", Sort: 32, IsUsed: true},
		{Name: "Claude 3 Opus", Provider: "Anthropic", Sort: 33, IsUsed: true},
		{Name: "Claude 3 Sonnet", Provider: "Anthropic", Sort: 34, IsUsed: true},
		{Name: "Claude 3 Haiku", Provider: "Anthropic", Sort: 35, IsUsed: true},
		{Name: "Claude 4 Sonnet", Provider: "Anthropic", Sort: 36, IsUsed: true},
		{Name: "Claude 4 Opus", Provider: "Anthropic", Sort: 37, IsUsed: true},
		{Name: "Claude 4.5 Sonnet", Provider: "Anthropic", Sort: 38, IsUsed: true},

		// ===== Google / Gemini 系列 =====
		{Name: "Gemini 2.5 Pro", Provider: "Google", Sort: 40, IsUsed: true},
		{Name: "Gemini 2.5 Flash", Provider: "Google", Sort: 41, IsUsed: true},
		{Name: "Gemini 2.0 Flash", Provider: "Google", Sort: 42, IsUsed: true},
		{Name: "Gemini 2.0 Flash-Lite", Provider: "Google", Sort: 43, IsUsed: true},
		{Name: "Gemini 1.5 Pro", Provider: "Google", Sort: 44, IsUsed: true},
		{Name: "Gemini 1.5 Flash", Provider: "Google", Sort: 45, IsUsed: true},
		{Name: "Gemini 1.0 Pro", Provider: "Google", Sort: 46, IsUsed: false},
		{Name: "Gemini 3 Pro", Provider: "Google", Sort: 47, IsUsed: true},
		{Name: "Gemini 3 Flash", Provider: "Google", Sort: 48, IsUsed: true},

		// ===== 智谱 GLM 系列 =====
		{Name: "GLM-4", Provider: "智谱AI", Sort: 50, IsUsed: true},
		{Name: "GLM-4-Plus", Provider: "智谱AI", Sort: 51, IsUsed: true},
		{Name: "GLM-4-Air", Provider: "智谱AI", Sort: 52, IsUsed: true},
		{Name: "GLM-4-Long", Provider: "智谱AI", Sort: 53, IsUsed: true},
		{Name: "GLM-4-Flash", Provider: "智谱AI", Sort: 54, IsUsed: true},
		{Name: "GLM-4V", Provider: "智谱AI", Sort: 55, IsUsed: true},
		{Name: "GLM-4-0520", Provider: "智谱AI", Sort: 56, IsUsed: false},
		{Name: "GLM-4.5", Provider: "智谱AI", Sort: 57, IsUsed: true},
		{Name: "GLM-4.6", Provider: "智谱AI", Sort: 58, IsUsed: true},
		{Name: "GLM-4.7", Provider: "智谱AI", Sort: 59, IsUsed: true},
		{Name: "GLM-Z1", Provider: "智谱AI", Sort: 60, IsUsed: true},
		{Name: "GLM-Z1-Air", Provider: "智谱AI", Sort: 61, IsUsed: true},

		// ===== 阿里 通义千问 Qwen 系列 =====
		{Name: "Qwen-Max", Provider: "阿里云", Sort: 70, IsUsed: true},
		{Name: "Qwen-Plus", Provider: "阿里云", Sort: 71, IsUsed: true},
		{Name: "Qwen-Turbo", Provider: "阿里云", Sort: 72, IsUsed: true},
		{Name: "Qwen2.5-Max", Provider: "阿里云", Sort: 73, IsUsed: true},
		{Name: "Qwen2.5-Plus", Provider: "阿里云", Sort: 74, IsUsed: true},
		{Name: "Qwen2.5-Turbo", Provider: "阿里云", Sort: 75, IsUsed: true},
		{Name: "Qwen2.5-72B", Provider: "阿里云", Sort: 76, IsUsed: true},
		{Name: "Qwen2.5-32B", Provider: "阿里云", Sort: 77, IsUsed: true},
		{Name: "Qwen2.5-14B", Provider: "阿里云", Sort: 78, IsUsed: true},
		{Name: "Qwen2.5-7B", Provider: "阿里云", Sort: 79, IsUsed: true},
		{Name: "Qwen2.5-1.5B", Provider: "阿里云", Sort: 80, IsUsed: false},
		{Name: "Qwen3-Max", Provider: "阿里云", Sort: 81, IsUsed: true},
		{Name: "Qwen3-235B-A22B", Provider: "阿里云", Sort: 82, IsUsed: true},
		{Name: "Qwen3-30B-A3B", Provider: "阿里云", Sort: 83, IsUsed: true},
		{Name: "Qwen3-VL", Provider: "阿里云", Sort: 84, IsUsed: true},
		{Name: "Qwen2-VL-72B", Provider: "阿里云", Sort: 85, IsUsed: false},
		{Name: "Qwen2-VL-7B", Provider: "阿里云", Sort: 86, IsUsed: false},

		// ===== DeepSeek 系列 =====
		{Name: "DeepSeek-V3", Provider: "DeepSeek", Sort: 90, IsUsed: true},
		{Name: "DeepSeek-V3-0324", Provider: "DeepSeek", Sort: 91, IsUsed: true},
		{Name: "DeepSeek-R1", Provider: "DeepSeek", Sort: 92, IsUsed: true},
		{Name: "DeepSeek-R1-0528", Provider: "DeepSeek", Sort: 93, IsUsed: true},
		{Name: "DeepSeek-R1-Distill-Qwen-32B", Provider: "DeepSeek", Sort: 94, IsUsed: true},
		{Name: "DeepSeek-R1-Distill-Qwen-14B", Provider: "DeepSeek", Sort: 95, IsUsed: true},
		{Name: "DeepSeek-R1-Distill-Llama-70B", Provider: "DeepSeek", Sort: 96, IsUsed: true},
		{Name: "DeepSeek-Chat", Provider: "DeepSeek", Sort: 97, IsUsed: true},
		{Name: "DeepSeek-Reasoner", Provider: "DeepSeek", Sort: 98, IsUsed: true},
		{Name: "DeepSeek-V3.1", Provider: "DeepSeek", Sort: 99, IsUsed: true},
		{Name: "DeepSeek-R2", Provider: "DeepSeek", Sort: 100, IsUsed: true},

		// ===== 字节 豆包 / 火山方舟 =====
		{Name: "豆包-Doubao-Pro", Provider: "字节跳动", Sort: 110, IsUsed: true},
		{Name: "豆包-Doubao-Pro-32K", Provider: "字节跳动", Sort: 111, IsUsed: true},
		{Name: "豆包-Doubao-Lite", Provider: "字节跳动", Sort: 112, IsUsed: true},
		{Name: "豆包-Doubao-Vision", Provider: "字节跳动", Sort: 113, IsUsed: true},
		{Name: "豆包-Doubao-Seed", Provider: "字节跳动", Sort: 114, IsUsed: true},
		{Name: "豆包-Doubao-Seed-1.6", Provider: "字节跳动", Sort: 115, IsUsed: true},
		{Name: "豆包-Doubao-1.5", Provider: "字节跳动", Sort: 116, IsUsed: true},
		{Name: "豆包-Doubao-1.5-Vision", Provider: "字节跳动", Sort: 117, IsUsed: true},

		// ===== 月之暗面 Kimi 系列 =====
		{Name: "Kimi-K2", Provider: "月之暗面", Sort: 120, IsUsed: true},
		{Name: "Kimi-K2-0905", Provider: "月之暗面", Sort: 121, IsUsed: true},
		{Name: "Kimi-K2-Turbo", Provider: "月之暗面", Sort: 122, IsUsed: true},
		{Name: "Kimi-K1.5", Provider: "月之暗面", Sort: 123, IsUsed: true},
		{Name: "Kimi-K1.5-Long", Provider: "月之暗面", Sort: 124, IsUsed: true},
		{Name: "moonshot-v1-8k", Provider: "月之暗面", Sort: 125, IsUsed: false},
		{Name: "moonshot-v1-32k", Provider: "月之暗面", Sort: 126, IsUsed: false},
		{Name: "moonshot-v1-128k", Provider: "月之暗面", Sort: 127, IsUsed: false},

		// ===== 百度 文心一言 ERNIE 系列 =====
		{Name: "文心-ERNIE-4.5", Provider: "百度", Sort: 130, IsUsed: true},
		{Name: "文心-ERNIE-4.5-Turbo", Provider: "百度", Sort: 131, IsUsed: true},
		{Name: "文心-ERNIE-4.0", Provider: "百度", Sort: 132, IsUsed: true},
		{Name: "文心-ERNIE-4.0-Turbo", Provider: "百度", Sort: 133, IsUsed: true},
		{Name: "文心-ERNIE-3.5", Provider: "百度", Sort: 134, IsUsed: true},
		{Name: "文心-ERNIE-3.5-Turbo", Provider: "百度", Sort: 135, IsUsed: false},
		{Name: "文心-ERNIE-Speed", Provider: "百度", Sort: 136, IsUsed: true},
		{Name: "文心-ERNIE-Lite", Provider: "百度", Sort: 137, IsUsed: true},
		{Name: "文心-ERNIE-VL", Provider: "百度", Sort: 138, IsUsed: false},
		{Name: "文心-X1", Provider: "百度", Sort: 139, IsUsed: true},

		// ===== 腾讯 混元 Hunyuan 系列 =====
		{Name: "混元-Hunyuan-A13B", Provider: "腾讯", Sort: 140, IsUsed: true},
		{Name: "混元-Hunyuan-Turbo", Provider: "腾讯", Sort: 141, IsUsed: true},
		{Name: "混元-Hunyuan-Large", Provider: "腾讯", Sort: 142, IsUsed: true},
		{Name: "混元-Hunyuan-Pro", Provider: "腾讯", Sort: 143, IsUsed: true},
		{Name: "混元-Hunyuan-Standard", Provider: "腾讯", Sort: 144, IsUsed: true},
		{Name: "混元-Hunyuan-Lite", Provider: "腾讯", Sort: 145, IsUsed: true},
		{Name: "混元-Hunyuan-T1", Provider: "腾讯", Sort: 146, IsUsed: true},

		// ===== 讯飞星火 Spark 系列 =====
		{Name: "讯飞星火-Spark-4.0", Provider: "讯飞", Sort: 150, IsUsed: true},
		{Name: "讯飞星火-Spark-3.5", Provider: "讯飞", Sort: 151, IsUsed: true},
		{Name: "讯飞星火-Spark-3.0", Provider: "讯飞", Sort: 152, IsUsed: false},
		{Name: "讯飞星火-Spark-Lite", Provider: "讯飞", Sort: 153, IsUsed: true},
		{Name: "讯飞星火-Spark-Max", Provider: "讯飞", Sort: 154, IsUsed: true},
		{Name: "讯飞星火-Spark-V4.0", Provider: "讯飞", Sort: 155, IsUsed: false},

		// ===== Meta / Llama 系列 =====
		{Name: "Llama 4 Scout", Provider: "Meta", Sort: 160, IsUsed: true},
		{Name: "Llama 4 Maverick", Provider: "Meta", Sort: 161, IsUsed: true},
		{Name: "Llama 4 Behemoth", Provider: "Meta", Sort: 162, IsUsed: true},
		{Name: "Llama 3.3-70B", Provider: "Meta", Sort: 163, IsUsed: true},
		{Name: "Llama 3.1-405B", Provider: "Meta", Sort: 164, IsUsed: true},
		{Name: "Llama 3.1-70B", Provider: "Meta", Sort: 165, IsUsed: true},
		{Name: "Llama 3.1-8B", Provider: "Meta", Sort: 166, IsUsed: true},
		{Name: "Llama 3-70B", Provider: "Meta", Sort: 167, IsUsed: false},
		{Name: "Llama 3-8B", Provider: "Meta", Sort: 168, IsUsed: false},

		// ===== Mistral 系列 =====
		{Name: "Mistral Large 2", Provider: "Mistral", Sort: 170, IsUsed: true},
		{Name: "Mistral Large", Provider: "Mistral", Sort: 171, IsUsed: true},
		{Name: "Mistral Medium", Provider: "Mistral", Sort: 172, IsUsed: true},
		{Name: "Mistral Small", Provider: "Mistral", Sort: 173, IsUsed: true},
		{Name: "Mistral-Nemo", Provider: "Mistral", Sort: 174, IsUsed: true},
		{Name: "Mixtral-8x7B", Provider: "Mistral", Sort: 175, IsUsed: true},
		{Name: "Codestral", Provider: "Mistral", Sort: 176, IsUsed: true},
		{Name: "Ministral-8B", Provider: "Mistral", Sort: 177, IsUsed: false},

		// ===== Google Gemma 系列 =====
		{Name: "Gemma 3 27B", Provider: "Google", Sort: 180, IsUsed: true},
		{Name: "Gemma 3 12B", Provider: "Google", Sort: 181, IsUsed: true},
		{Name: "Gemma 3 4B", Provider: "Google", Sort: 182, IsUsed: true},
		{Name: "Gemma 3 1B", Provider: "Google", Sort: 183, IsUsed: true},
		{Name: "Gemma 2 27B", Provider: "Google", Sort: 184, IsUsed: true},
		{Name: "Gemma 2 9B", Provider: "Google", Sort: 185, IsUsed: true},
		{Name: "Gemma 2 2B", Provider: "Google", Sort: 186, IsUsed: false},

		// ===== xAI / Grok 系列 =====
		{Name: "Grok 4", Provider: "xAI", Sort: 190, IsUsed: true},
		{Name: "Grok 4 Fast", Provider: "xAI", Sort: 191, IsUsed: true},
		{Name: "Grok 3", Provider: "xAI", Sort: 192, IsUsed: true},
		{Name: "Grok 3 mini", Provider: "xAI", Sort: 193, IsUsed: true},
		{Name: "Grok 3 Fast", Provider: "xAI", Sort: 194, IsUsed: true},
		{Name: "Grok 2", Provider: "xAI", Sort: 195, IsUsed: false},
		{Name: "Grok-2-1212", Provider: "xAI", Sort: 196, IsUsed: false},
		{Name: "Grok-Beta", Provider: "xAI", Sort: 197, IsUsed: false},

		// ===== 其他主流开源/闭源模型 =====
		{Name: "Stable Diffusion 3", Provider: "Stability AI", Sort: 200, IsUsed: true},
		{Name: "FLUX.1", Provider: "Black Forest Labs", Sort: 201, IsUsed: true},
		{Name: "DALL·E 3", Provider: "OpenAI", Sort: 202, IsUsed: true},
		{Name: "Midjourney", Provider: "Midjourney", Sort: 203, IsUsed: true},
		{Name: "Yi-Large", Provider: "零一万物", Sort: 204, IsUsed: true},
		{Name: "Yi-1.5-34B", Provider: "零一万物", Sort: 205, IsUsed: true},
		{Name: "MiniMax-Text-01", Provider: "MiniMax", Sort: 206, IsUsed: true},
		{Name: "MiniMax-M2", Provider: "MiniMax", Sort: 207, IsUsed: true},
		{Name: "MiniMax-ABAB-6.5", Provider: "MiniMax", Sort: 208, IsUsed: false},
		{Name: "Step-2", Provider: "阶跃星辰", Sort: 209, IsUsed: true},
		{Name: "Step-3", Provider: "阶跃星辰", Sort: 210, IsUsed: true},
		{Name: "Step-1V", Provider: "阶跃星辰", Sort: 211, IsUsed: true},
		{Name: "InternLM2.5-7B", Provider: "上海AI实验室", Sort: 212, IsUsed: true},
		{Name: "InternVL2.5-78B", Provider: "上海AI实验室", Sort: 213, IsUsed: false},
		{Name: "Baichuan2-53B", Provider: "百川智能", Sort: 214, IsUsed: false},
		{Name: "Baichuan3", Provider: "百川智能", Sort: 215, IsUsed: true},
		{Name: "Phi-4", Provider: "Microsoft", Sort: 216, IsUsed: true},
		{Name: "Phi-4-mini", Provider: "Microsoft", Sort: 217, IsUsed: true},
		{Name: "Phi-3.5-mini", Provider: "Microsoft", Sort: 218, IsUsed: false},
		{Name: "Command R+", Provider: "Cohere", Sort: 219, IsUsed: true},
		{Name: "Command A", Provider: "Cohere", Sort: 220, IsUsed: true},
		{Name: "Jamba", Provider: "AI21", Sort: 221, IsUsed: true},
		{Name: "DSLite", Provider: "深度求索", Sort: 222, IsUsed: true},
		{Name: "NVLM-D-72B", Provider: "NVIDIA", Sort: 223, IsUsed: false},
		{Name: "Nemotron-4-340B", Provider: "NVIDIA", Sort: 224, IsUsed: false},
		{Name: "Skywork-13B", Provider: "昆仑万维", Sort: 225, IsUsed: false},
		{Name: "Skywork-Llama2", Provider: "昆仑万维", Sort: 226, IsUsed: false},
	}

	for i := range models {
		models[i].CreatedAt = &now
		models[i].UpdatedAt = &now
	}

	if err := db.Create(&models).Error; err != nil {
		return
	}
}
