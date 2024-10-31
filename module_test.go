package ai

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
	"google.golang.org/api/option"
)

func Test_Module(t *testing.T) {
	userController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("users")

		ctrl.Get("", func(ctx core.Ctx) error {
			model := InjectModel(module, "gemini-1.5-flash")
			resp, err := model.GenerateContent(context.Background(), genai.Text("Write a story about a magic backpack."))
			if err != nil {
				fmt.Println(err)
				return err
			}
			return ctx.JSON(core.Map{
				"data": resp,
			})
		})
		return ctrl
	}

	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		return module.New(core.NewModuleOptions{
			Imports: []core.Module{
				ForFeature("gemini-1.5-flash"),
			},
			Controllers: []core.Controller{
				userController,
			},
		})
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{
				ForRoot(option.WithAPIKey(os.Getenv("API_KEY"))),
				userModule,
			},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/users")
	require.Nil(t, err)
	require.Equal(t, 200, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	fmt.Println(string(data))
}
