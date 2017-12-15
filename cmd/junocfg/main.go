package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/gojuno/junocfg"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "jtools",
	Short: "automation tools",
	Long:  `jtools is a CLI tools suite for automation of development activity.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//initLogger()
	},
	Run: func(cmd *cobra.Command, args []string) {

		filenames := strings.Split(viper.GetString("input"), ",")

		correct := true

		if viper.GetBool("check-tmpl") {
			// fmt.Fprintf(os.Stderr, "mode: check-tmpl\n")
			if rawTmpl, err := junocfg.ReadData(viper.GetString("template")); err != nil {
				log.Fatal(err.Error())
				os.Exit(1)
			} else {
				if tmpl, err := template.New("template").Parse(rawTmpl); err != nil {
					fmt.Fprintf(os.Stderr, "template error: [%v]\n", err)
					os.Exit(1)
				} else {
					fmt.Fprintf(os.Stderr, "%s\n", tmpl.Root.String())
				}
			}
		} else if viper.GetBool("merge") {
			// fmt.Fprintf(os.Stderr, "mode: merge\n")
			if cfg, err := junocfg.GetConfigs(filenames); err == nil {
				d, err := yaml.Marshal(cfg)
				if err != nil {
					log.Fatal(fmt.Sprintf("Could not create YAML file: %s\n", err))
				}
				buffer := bytes.NewBuffer([]byte{})
				buffer.Write(d)
				junocfg.OutResult(viper.GetString("output"), buffer)
			} else {
				log.Fatalf("%s\n", err)
			}
		} else {
			// fmt.Fprintf(os.Stderr, "mode: default\n")
			// default: generate config file from input + template
			buffer := bytes.NewBuffer([]byte{})

			if viper.GetString("template") == "" {
				log.Fatal("Field [template(-t|--tmpl)] required")
			}
			tmpl, err := junocfg.GetTemplate(viper.GetString("template"))
			if err != nil {
				log.Fatal(err.Error())
			}

			cfg := make(map[string]interface{})
			if config, err := junocfg.GetConfigs(filenames); err != nil {
				log.Fatalf("%s\n", err)
			} else {
				cfg = config
			}

			if err := tmpl.Execute(buffer, cfg); err != nil {
				log.Fatalf("failed to render template [%s]\n[%s]\n", err, cfg)
			}

			buffer = junocfg.PreprocessYaml(buffer)

			// check yaml
			outYaml := map[string]interface{}{}
			if err = yaml.Unmarshal(buffer.Bytes(), &outYaml); err != nil {
				fmt.Fprintf(os.Stderr, "Not valid output yaml: %s\n", err.Error())
				correct = correct && false
			}

			// check variables
			if viper.GetBool("check") {
				strOut := strings.Split(buffer.String(), "\n")

				for posInFile, str := range strOut {
					if i := strings.Index(str, "<no value>"); i != -1 {
						fmt.Fprintf(os.Stderr, "<no value> at %s#%d:%s\n", viper.GetString("output"), posInFile, str)
						correct = correct && false
					}
				}

			}
			junocfg.OutResult(viper.GetString("output"), buffer)
		}

		if !correct {
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.junocfg.yaml)")

	RootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug mode")
	RootCmd.PersistentFlags().BoolP("dry-run", "", false, "Enable dry-run mode")

	RootCmd.PersistentFlags().BoolP("check", "", false, "check")
	RootCmd.PersistentFlags().BoolP("check-tmpl", "", false, "check tmpl")
	RootCmd.PersistentFlags().BoolP("merge", "", false, "merge")

	RootCmd.PersistentFlags().StringP("input", "i", "", "input")
	RootCmd.PersistentFlags().StringP("output", "o", "", "output")
	RootCmd.PersistentFlags().StringP("template", "t", "", "template")

	for _, f := range []string{"debug", "dry-run",
		"check", "check-tmpl", "merge",
		"input", "output", "template",
	} {
		viper.BindPFlag(f, RootCmd.PersistentFlags().Lookup(f))
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".junocfg")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	}

	if viper.GetBool("debug") {
		cfgJson, _ := json.Marshal(viper.AllSettings())
		fmt.Fprintf(os.Stderr, "args: %v\n", string(cfgJson))
	}
}

func main() {
	Execute()
}
