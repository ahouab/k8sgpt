/*
Copyright 2023 The K8sGPT Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package analyze

import (
	"fmt"
	"os"

	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	"github.com/spf13/cobra"
)

var (
	explain        bool
	backend        string
	output         string
	filters        []string
	language       string
	nocache        bool
	namespace      string
	anonymize      bool
	maxConcurrency int
	withDoc        bool

	example = templates.Examples(i18n.T(`

	# analyze all supported resources across all namespaces
	k8sgpt analyze

	# explain with AI power
	k8sgpt analyze --explain

	# filter by namespaces
	k8sgpt analyze --namespace=default

	# filter by resource type
	k8sgpt analyze --explain --filter=Service -n default

	# filter by multiple resource types
	k8sgpt analyze --explain --filter=Service,Pod

	# analyze specific resources(different kinds)
	k8sgpt analyze pod/nginx-7b6fc59cd6-shz7h service/nginx   -n default --explain

	# analyze specific resources with filter(resource kind can be implicit when single filter)
	k8sgpt analyze nginx-7b6fc59cd6-shz7h nginx-7db4cf8964-b26wf -n default --filter=Pod

	# output to json
	k8sgpt analyze --explain --filter=Service --output=json

	# anonymize during explai
	k8sgpt analyze --explain --filter=Service --output=json --anonymize`))
)

// AnalyzeCmd represents the problems command
var AnalyzeCmd = &cobra.Command{
	Use:     "analyze <specific-resources>",
	Aliases: []string{"analyse"},
	Example: example,
	Short:   "This command will find problems within your Kubernetes cluster",
	Long: `This command will find problems within your Kubernetes cluster and
	provide you with a list of issues that need to be resolved`,
	Run: func(cmd *cobra.Command, args []string) {

		// Create analysis configuration first.
		config, err := analysis.NewAnalysis(
			backend,
			language,
			filters,
			namespace,
			nocache,
			explain,
			maxConcurrency,
			withDoc,
			args,
		)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		defer config.Close()

		config.RunAnalysis()

		if explain {
			if err := config.GetAIResults(output, anonymize); err != nil {
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
	// add max concurrency
	AnalyzeCmd.Flags().IntVarP(&maxConcurrency, "max-concurrency", "m", 10, "Maximum number of concurrent requests to the Kubernetes API server")
	// kubernetes doc flag
	AnalyzeCmd.Flags().BoolVarP(&withDoc, "with-doc", "d", false, "Give me the official documentation of the involved field")
}
