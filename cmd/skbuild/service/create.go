package service

import (
	"github.com/iancoleman/strcase"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skema-dev/skemabuild/internal/generator"
	"github.com/skema-dev/skemabuild/internal/pkg/console"
	"github.com/skema-dev/skemabuild/internal/pkg/io"
	"github.com/spf13/cobra"
)

const (
	createDescription     = "Create service code from protocol buffers definition"
	createLongDescription = "skbuild service create --proto=<protobuf_uri>"
)

func newCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: createDescription,
		Long:  createLongDescription,
		Run: func(c *cobra.Command, args []string) {
			protoUrl := c.Flag("proto").Value.String()
			goModule := c.Flag("module").Value.String()
			goVersion, _ := c.Flags().GetString("goversion")
			serviceName, _ := c.Flags().GetString("service")
			output, _ := c.Flags().GetString("output")
			tpl, _ := c.Flags().GetString("tpl")
			s, _ := c.Flags().GetString("http")
			httpEnabled, _ := strconv.ParseBool(s)
			values, _ := c.Flags().GetString("value")

			userValues := map[string]string{}
			if values != "" {
				ss := strings.Split(values, ",")
				for _, s := range ss {
					kv := strings.Split(s, ":")
					if len(kv) != 2 {
						console.Fatalf("Invalid parameter: %s", s)
					}
					k := strcase.ToCamel(kv[0])
					v := kv[1]
					userValues[k] = v
				}
			}

			modelNames := make([]string, 0)
			modelParams, _ := c.Flags().GetString("model")
			if modelParams != "" {
				modelNames = strings.Split(modelParams, ",")
			}

			serviceTemplate := generator.CreateServiceTemplate().
				WithRpcProtocol(protoUrl, goModule, goVersion, serviceName, httpEnabled).
				WithDataModelNames(modelNames).
				WithUserValues(userValues)

			generator := generator.NewGrpcGoGenerator()
			contents := generator.CreateCodeContent(tpl, serviceTemplate)

			for path, c := range contents {
				outputPath := filepath.Join(output, path)
				io.SaveToFile(outputPath, []byte(c))
				console.Info(outputPath)
			}
		},
	}

	cmd.Flags().StringP("proto", "p", "", "protobuf file")
	cmd.Flags().StringP("module", "m", "", "go module name")
	cmd.Flags().StringP("goversion", "v", "1.16", "go version")
	cmd.Flags().StringP("service", "s", "", "service name")
	cmd.Flags().StringP("tpl", "t", "skema-mux", "template name or url")
	cmd.Flags().String("http", "true", "enable http or not")
	cmd.Flags().StringP("output", "o", "", "output path")
	cmd.Flags().String("value", "", "user defined tpl parameters: key1:value1,key2:value2...")
	cmd.Flags().String("model", "", "data models supported by skema-data")
	cmd.MarkFlagRequired("proto")

	return cmd
}
