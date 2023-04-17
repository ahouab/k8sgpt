package analyze

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	"github.com/spf13/cobra"
)

var (
	explain   bool
	backend   string
	output    string
	filters   []string
	language  string
	nocache   bool
	namespace string
	anonymize bool
)

// AnalyzeCmd represents the problems command
var AnalyzeCmd = &cobra.Command{
	Use:     "analyze",
	Aliases: []string{"analyse"},
	Short:   "This command will find problems within your Kubernetes cluster",
	Long: `This command will find problems within your Kubernetes cluster and
	provide you with a list of issues that need to be resolved`,
	Run: func(cmd *cobra.Command, args []string) {

		// AnalysisResult configuration
		config, err := analysis.NewAnalysis(backend, language, filters, namespace, nocache, explain)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		err = config.RunAnalysis()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if explain {
			err := config.GetAIResults(output, anonymize)
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
		}

		// print results
		output, err := config.PrintOutput(output)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func init() {

	// namespace flag
	AnalyzeCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to analyze")
	// no cache flag
	AnalyzeCmd.Flags().BoolVarP(&nocache, "no-cache", "c", false, "Do not use cached data")
	// anonymize flag
	AnalyzeCmd.Flags().BoolVarP(&anonymize, "anonymize", "a", false, "Anonymize data before sending it to the AI backend. This flag masks sensitive data, such as Kubernetes object names and labels, by replacing it with a key. However, please note that this flag does not currently apply to events.")
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
