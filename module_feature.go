package ai

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/tinh-tinh/tinhtinh/core"
)

func GetModelName(name string) core.Provide {
	return core.Provide(fmt.Sprintf("model-%s", name))
}

func ForFeature(name string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		modelModule := module.New(core.NewModuleOptions{})

		modelModule.NewProvider(core.ProviderOptions{
			Name: GetModelName(name),
			Factory: func(param ...interface{}) interface{} {
				client := param[0].(*genai.Client)
				model := client.GenerativeModel(name)
				return model
			},
			Inject: []core.Provide{GEN_AI_CLIENT},
		})
		modelModule.Export(GetModelName(name))
		return modelModule
	}
}

func InjectModel(module *core.DynamicModule, name string) *genai.GenerativeModel {
	model, ok := module.Ref(GetModelName(name)).(*genai.GenerativeModel)
	if !ok {
		return nil
	}
	return model
}
