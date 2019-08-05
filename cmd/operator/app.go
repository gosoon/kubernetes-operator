/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operator

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientset "github.com/gosoon/kubernetes-operator/pkg/client/clientset/versioned"
	informers "github.com/gosoon/kubernetes-operator/pkg/client/informers/externalversions"
	ecsv1controller "github.com/gosoon/kubernetes-operator/pkg/controller"
	"github.com/gosoon/kubernetes-operator/pkg/kuberesource"
	"github.com/gosoon/kubernetes-operator/pkg/server"
	ctrl "github.com/gosoon/kubernetes-operator/pkg/server/controller"

	"github.com/gosoon/glog"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/resouer/k8s-controller-custom-resource/pkg/signals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/kubernetes/pkg/controller"
)

var (
	cfgFile    string
	masterURL  string
	kubeconfig string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubernetes-operator",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		// set up signals so we handle the first shutdown signal gracefully
		stopCh := signals.SetupSignalHandler()

		var cfg *rest.Config
		var err error
		if kubeconfig == "" {
			cfg, err = rest.InClusterConfig()
		} else {
			cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
			if err != nil {
				glog.Fatalf("Error building kubeconfig: %s", err.Error())
			}
		}

		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
		}

		kubernetesClusterClient, err := clientset.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building kubernetesCluster clientset: %s", err.Error())
		}

		// start http server
		go func() {
			opt := &ctrl.Options{KubernetesClusterClientset: kubernetesClusterClient, KubeClientset: kubeClient}
			server := server.New(server.Options{CtrlOptions: opt, ListenAddr: ":8080"})
			if err := server.ListenAndServe(); err != nil {
				glog.Fatalf("Failed to listen and serve admission webhook server: %v", err)
			}
			glog.Info("Server started")
		}()

		// new resourcelock
		rl, err := kuberesource.NewResourceLock(kubeClient)
		if err != nil {
			glog.Fatalf("new resourcelock failed with:%s", err.Error())
		}

		// add leader elector
		run := func(ctx context.Context) {
			kubernetesClusterInformerFactory := informers.NewSharedInformerFactory(kubernetesClusterClient, controller.NoResyncPeriodFunc())
			ecsv1Controller := ecsv1controller.NewController(kubeClient, kubernetesClusterClient,
				kubernetesClusterInformerFactory.Ecs().V1().KubernetesClusters())

			go kubernetesClusterInformerFactory.Start(stopCh)
			go func() {
				if err = ecsv1Controller.Run(2, stopCh); err != nil {
					glog.Fatalf("Error running controller: %s", err.Error())
				}
			}()

			// listening OS shutdown singal
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			<-signalChan

			glog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
			//server.Shutdown(context.Background())
		}

		leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
			Lock:          rl,
			LeaseDuration: 15 * time.Second,
			RenewDeadline: 10 * time.Second,
			RetryPeriod:   2 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: run,
				OnStoppedLeading: func() {
					glog.Info("leaderelection lost")
				},
			},
			Name: ecsv1controller.ComponentName,
		})

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubernetes-operator.yaml)")

	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	rootCmd.PersistentFlags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kubernetes-operator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubernetes-operator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
