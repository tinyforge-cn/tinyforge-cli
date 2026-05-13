package commands

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/tinycs-cn/cli/internal/client"
	"github.com/tinycs-cn/cli/internal/config"
)

// TestCommand implements "tinycs test".
//
// Flow:
//  1. Load auth token — required for stage inference.
//  2. Resolve course + language + project dir (same as submit).
//  3. Resolve tester: --tester-path forces binary; --docker forces Docker;
//     otherwise auto-select (Windows → Docker, macOS/Linux → binary).
//  4. Resolve stage slug:
//     a. --stage <slug>  → use as-is
//     b. --all           → omit -s, tester runs all stages
//     c. default         → GET /v1/current-stage?course=&language= (requires auth)
//  5. exec tester with resolved args, inheriting stdout/stderr.
func TestCommand(args []string) error {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	stage := flags.String("stage", "", "指定评测关卡 (slug)")
	all := flags.Bool("all", false, "测试所有关卡")
	localTester := flags.String("tester-path", "", "直接指定本地 tester 二进制路径（调试用）")
	useDocker := flags.Bool("docker", false, "使用 Docker 容器运行 tester（macOS/Linux 可选，Windows 默认）")
	apiURL := flags.String("api-url", "", "API 地址（内部测试用）")
	if err := flags.Parse(args); err != nil {
		return err
	}

	// 1. Auth
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	token := cfg.GetToken()
	if token == "" {
		return errors.New("未登录，请先运行: tinycs login")
	}

	// 2. Resolve course + language + project dir
	course, language, projectDir, err := resolveProject()
	if err != nil {
		return err
	}

	// 3. Resolve tester runner
	var runner testerRunner
	if *localTester != "" {
		runner = &binaryRunner{path: *localTester}
	} else {
		runner, err = ensureTester(course, *useDocker)
		if err != nil {
			return err
		}
	}

	// 4. Resolve stage slug
	stageSlug := *stage
	if stageSlug == "" && !*all {
		stageSlug, err = resolveCurrentStageFromAPI(course, language, cfg, *apiURL)
		if err != nil {
			return err
		}
	}

	// 5. Run tester
	fmt.Println("🧪 本地评测（结果不计入提交记录）")
	fmt.Println()
	return runner.Run(projectDir, stageSlug, *all)
}

// resolveCurrentStageFromAPI calls GET /v1/current-stage?course=&language= and
// returns the current stage slug. Produces user-friendly errors for 400/404.
func resolveCurrentStageFromAPI(course, language string, cfg *config.Config, apiURL string) (string, error) {
	baseURL := cfg.GetAPIURL(apiURL)
	c := client.New(baseURL, cfg.GetToken())

	item, err := c.GetCurrentStage(course, language)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case http.StatusBadRequest:
				return "", fmt.Errorf(
					"该课程为随机挑战模式，请用 --stage <slug> 指定关卡\n运行 tinycs stages 查看可用关卡",
				)
			case http.StatusNotFound:
				return "", fmt.Errorf(
					"未找到该课程的注册记录，请先到 https://www.tinycs.cn 注册仓库后再运行",
				)
			}
		}
		return "", fmt.Errorf("获取当前关卡失败: %w", err)
	}
	return item.Slug, nil
}
