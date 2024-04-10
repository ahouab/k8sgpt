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
package cache

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/cache"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List cache providers",
	Long:  `The list command displays a list of available cache providers with their status.`,
	Run: func(cmd *cobra.Command, args []string) {

		// load remote cache if it is configured
		c, err := cache.GetCacheConfiguration()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		fmt.Print(color.YellowString("Active: \n"))
		fmt.Printf("> %s\n", color.GreenString("%s", c.GetName()))
		fmt.Print(color.YellowString("Unused: \n"))
		for _, cache := range cache.GetAllCacheProviders() {
			if cache != c.GetName() {
				fmt.Printf("> %s\n", color.RedString("%s", cache))
			}
		}
	},
}

