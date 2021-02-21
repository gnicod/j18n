/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gnicod/j18n/config"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var lang string
var force bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new translation key",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jpath := args[0]
		config := config.NewConfig()
		for lang, path := range config.Langs {
			fullPath := fmt.Sprintf("%s/%s", config.BasePath, path)
			jsonContent := getJsonContent(fullPath)
			if !force {
				currentValue := gjson.Get(jsonContent, jpath).String()
				if currentValue != "" {
					fmt.Printf("Value for key %s and lang %s (%s) already exists, passing. Use -f to force change the value\n", jpath, lang, currentValue)
					continue
				}
			}
			value := promptTranslation(lang)
			newJsonValue, _ := sjson.Set(jsonContent, jpath, value)
			newJsonValue = jsonPrettyPrint(newJsonValue)
			_ = ioutil.WriteFile(fullPath, []byte(newJsonValue), 0644)
		}
	},
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func getJsonContent(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	return text
}

func promptTranslation(lang string) string {
	fmt.Printf("Value for %s: \n", lang)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return scanner.Text()

}

func init() {
	addCmd.Flags().StringVarP(&lang, "lang", "l", "", "Lang to set the value")
	addCmd.Flags().BoolVarP(&force, "force", "f", false, "Change the value of the translation even if the value already exist")
	rootCmd.AddCommand(addCmd)
}
