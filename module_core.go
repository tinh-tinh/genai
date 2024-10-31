package ai

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/tinh-tinh/tinhtinh/core"
	"google.golang.org/api/option"
)

const GEN_AI_CLIENT core.Provide = "gen-ai-client"

func ForRoot(opts ...option.ClientOption) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		genAiModule := module.New(core.NewModuleOptions{})

		ctx := context.Background()
		client, err := genai.NewClient(ctx, opts...)
		if err != nil {
			panic(err)
		}
		genAiModule.NewProvider(core.ProviderOptions{
			Name:  GEN_AI_CLIENT,
			Value: client,
		})
		genAiModule.Export(GEN_AI_CLIENT)
		return genAiModule
	}
}

func InjectClient(module *core.DynamicModule) *genai.Client {
	client, ok := module.Ref(GEN_AI_CLIENT).(*genai.Client)
	if !ok {
		return nil
	}
	return client
}
