// Copyright © 2019 suquiya
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/suquiya/swimmy"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewSwimmyCmd return the base command when called without any subcommands
func NewSwimmyCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "swimmy input <output>",
		Short: "Swimmy is a tool that pull meta info from url and write info to html or json",
		Long: `Swimmy is a tool that pull meta info from url and write info to html or json.
It is a package that fetch URL info and process it. It is for embedding external site information as card or outputting json.
Usage: swimmy url outputPath.
More details, Please type "swimmy --help" and enter.
`,
		RunE: func(cmd *cobra.Command, args []string) error {

			l, err := cmd.Flags().GetBool("list")
			if err != nil {
				return err
			}

			i, err := cmd.Flags().GetString("IfOutputExist")

			if err != nil {
				return err
			}

			formatStr, err := cmd.Flags().GetString("format")

			if err != nil {
				return err
			}

			tempSrc, err := cmd.Flags().GetString("template")

			if err != nil {
				return err
			}

			tojson := formatStr == "json" || strings.HasPrefix(formatStr, "json,") || strings.HasSuffix(formatStr, ",json")

			tohtml := formatStr == "html" || strings.HasPrefix(formatStr, "html,") || strings.HasSuffix(formatStr, ",html")

			ExistUserTemp := false
			var userTemp *template.Template
			if tohtml {
				if e, _ := swimmy.IsExistFilePath(tempSrc); e {
					tmpl, err := template.ParseFiles(tempSrc)
					if err == nil {
						ExistUserTemp = true
						userTemp = tmpl
					} else {
						ExistUserTemp = false
					}
				} else {
					tempSrc = ""
				}
			}
			fmt.Println(tojson, tohtml)

			i = strings.ToLower(i)

			argNum := len(cmd.Flags().Args())
			if argNum < 1 {
				return fmt.Errorf("swimmy requires at least one arguments: ex. swimmy url <output>")
			}

			input := cmd.Flags().Arg(0)

			output := ""

			var ow *bufio.Writer
			var owf *os.File
			owf = nil

			/*
				if output == "" {
					ow = cmd.OutOrStdout()
				} else {
					is, err := swimmy.IsFilePath(output)
					if is {

					} else {
						fmt.Println("Not filepath: %s, %s", output, err)
					}
				}
			*/

			if argNum > 1 {
				output = cmd.Flags().Arg(1)
				isfp, err := swimmy.IsFilePath(output)
				if !isfp {
					return err
				}
				if err == nil { //output file is already exist.

					if i == "o" || i == "overwrite" {
						owf, err = os.Create(output)
						defer owf.Close()
						if err != nil {
							return err
						}
						ow = bufio.NewWriter(owf)

					} else if i == "a" || i == "append" {
						owf, err = os.OpenFile(output, os.O_WRONLY|os.O_APPEND, 0666)
						defer owf.Close()

						if err != nil {
							return err
						}
						ow = bufio.NewWriter(owf)
					} else {
						ow = bufio.NewWriter(cmd.OutOrStdout())
					}
				} else { //output file does not exist.
					owf, err = os.Create(output)
					defer owf.Close()
					if err != nil {
						return err
					}
					ow = bufio.NewWriter(owf)
				}
			} else {
				ow = bufio.NewWriter(cmd.OutOrStdout())
			}

			swimmy.Init()

			cardBuilder := swimmy.DefaultCardBuilder
			if ExistUserTemp {
				cardBuilder = swimmy.NewCardBuilder(userTemp, swimmy.DefaultClasses())
			}
			if l {
				isfp, err := swimmy.IsExistFilePath(input)
				if !isfp {
					return err
				}
				f, err := os.Open(input)
				defer f.Close()
				if err != nil {
					return err
				}

				if tojson {
					ow.WriteString("[")
				}

				scanner := bufio.NewScanner(f)

				count := 0
				for scanner.Scan() {
					line := scanner.Text()
					if govalidator.IsRequestURL(line) {
						/*if count > 1 {
							ow.WriteString(",")
						}*/
						pd, err := swimmy.FetchAndBuildPageData(line, cmd.OutOrStdout())
						if err == nil {
							if tojson {
								err = swimmy.WriteJSON(pd, ow, cmd.OutOrStdout(), count > 1)
								if err != nil {
									cmd.Println(err)
								} else {
									count++
								}
							}
							if tohtml {
								err = swimmy.WriteHTML(pd, cardBuilder, ow, cmd.OutOrStdout(), count > 1)
								if err != nil {
									cmd.Println(err)
								} else {
									count++
								}
							}
						} else {
							cmd.Println(err)
						}

					} else {
						cmd.Println("This line is not url: ", line)
					}
				}

				if tojson {
					ow.WriteString("]")
				}

				if err := scanner.Err(); err != nil {
					panic(err)
				}

				return ow.Flush()
			}
			if govalidator.IsRequestURL(input) {

				pd, err := swimmy.FetchAndBuildPageData(input, cmd.OutOrStdout())
				if err == nil {
					if tojson {
						err = swimmy.WriteJSON(pd, ow, cmd.OutOrStdout(), false)

						if err != nil {
							cmd.Println(err)
						}
					}
					if tohtml {

						err = swimmy.WriteHTML(pd, cardBuilder, ow, cmd.OutOrStdout(), false)

						if err != nil {
							cmd.Println(err)
						}

					}
				} else {
					cmd.Println(err)
				}

				return ow.Flush()
			}
			cmd.Println("inputted url is not url: ", input)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swimmy.yaml)")

	rootCmd.Flags().BoolP("list", "l", false, "use list, you can specify urls list in txt format (separated by newline)")

	rootCmd.Flags().StringP("IfOutputExist", "i", "S", "this flag define behavior in case that specified output file is already exist: [A]ppend,[O]verwritte or [S]tdout, default is S.")

	rootCmd.Flags().StringP("format", "f", "json", "this flag decide format for outputting url information. You can specify \"json\" or \"html\". If you want to get two type formats, you should specify \"json,html\". Default is \"json\".")

	rootCmd.Flags().StringP("template", "t", "", "User can specify template file path by using this flag.")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := NewSwimmyCmd()
	rootCmd.SetOutput(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOutput(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swimmy.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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

		// Search config in home directory with name ".swimmy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".swimmy")

		//fmt.Println(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
