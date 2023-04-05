package analyze

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	"github.com/k8sgpt-ai/k8sgpt/pkg/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	explain   bool
	backend   string
	output    string
	filters   []string
	language  string
	nocache   bool
	namespace string
)

// AnalyzeCmd represents the problems command
var AnalyzeCmd = &cobra.Command{
	Use:     "analyze",
	Aliases: []string{"analyse"},
	Short:   "This command will find problems within your Kubernetes cluster",
	Long: `This command will find problems within your Kubernetes cluster and
	provide you with a list of issues that need to be resolved`,
	Run: func(cmd *cobra.Command, args []string) {

		// get ai configuration
		var configAI ai.AIConfiguration
		err := viper.UnmarshalKey("ai", &configAI)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if len(configAI.Providers) == 0 {
			color.Red("Error: AI provider not specified in configuration. Please run k8sgpt auth")
			os.Exit(1)
		}

		var aiProvider ai.AIProvider
		for _, provider := range configAI.Providers {
			if backend == provider.Name {
				aiProvider = provider
				break
			}
		}

		if aiProvider.Name == "" {
			color.Red("Error: AI provider %s not specified in configuration. Please run k8sgpt auth", backend)
			os.Exit(1)
		}

		aiClient := ai.NewClient(aiProvider.Name)
		if err := aiClient.Configure(aiProvider.Password, aiProvider.Model, language); err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		ctx := context.Background()
		// Get kubernetes client from viper
		client := viper.Get("kubernetesClient").(*kubernetes.Client)
		// AnalysisResult configuration
		config := &analysis.Analysis{
			Namespace: namespace,
			NoCache:   nocache,
			Filters:   filters,
			Explain:   explain,
			AIClient:  aiClient,
			Client:    client,
			Context:   ctx,
		}

		err = config.RunAnalysis()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if explain {
			err := config.GetAIResults(output)
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
		}

		// print results
		switch output {
		case "json":
			output, err := config.JsonOutput()
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
		default:
			config.PrintOutput()
		}
	},
}

func init() {

	// namespace flag
	AnalyzeCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to analyze")
	// no cache flag
	AnalyzeCmd.Flags().BoolVarP(&nocache, "no-cache", "c", false, "Do not use cached data")
	// array of strings flag
	AnalyzeCmd.Flags().StringSliceVarP(&filters, "filter", "f", []string{}, "Filter for these analyzers (e.g. Pod, PersistentVolumeClaim, Service, ReplicaSet)")
	// explain flag
	AnalyzeCmd.Flags().BoolVarP(&explain, "explain", "e", false, "Explain the problem to me")
	// add flag for backend
	AnalyzeCmd.Flags().StringVarP(&backend, "backend", "b", "openai", "Backend AI provider")
	// output as json
	AnalyzeCmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text, json)")
	// add language options for output
	AnalyzeCmd.Flags().StringVarP(&language, "language", "l", "english", "Languages to use for AI (e.g. 'English', 'Spanish', 'French', 'German', 'Italian', 'Portuguese', 'Dutch', 'Russian', 'Chinese', 'Japanese', 'Korean')")
}
