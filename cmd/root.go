package cmd

import (
	"os"

	"github.com/k8sgpt-ai/k8sgpt/cmd/generate"

	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/cmd/analyze"
	"github.com/k8sgpt-ai/k8sgpt/cmd/auth"
	"github.com/k8sgpt-ai/k8sgpt/pkg/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	masterURL  string
	kubeconfig string
	version    string
	k8sContext string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8sgpt",
	Short: "Kubernetes debugging powered by AI",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    //Initialise the kubeconfig
    kubernetesClient, err := kubernetes.NewClient(masterURL, kubeconfig, k8sContext)
    if err != nil {
      color.Red("Error initialising kubernetes client: %v", err)
      os.Exit(1)
    }

    viper.Set("kubernetesClient", kubernetesClient)
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	version = v
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(analyze.AnalyzeCmd)
	rootCmd.AddCommand(generate.GenerateCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8sgpt.yaml)")
	rootCmd.PersistentFlags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().StringVarP(&k8sContext, "kube-context", "", "", "Kubernetes context to be used")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".k8sgpt.git" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".k8sgpt")

		viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
