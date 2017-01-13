// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/markbates/gentronics"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:     "console",
	Aliases: []string{"c"},
	Short:   "Runs your Buffalo app in a REPL console",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := exec.LookPath("gore")
		if err != nil {
			return errors.New("we could not find \"gore\" in your path.You must first install \"gore\" in order to use the Buffalo console:\n\n$ go get -u github.com/motemen/gore")
		}
		rootPath, err := rootPath("")
		if err != nil {
			return err
		}
		packagePath := packagePath(rootPath)
		packages := []string{}
		for _, p := range []string{"models", "actions"} {
			s, _ := os.Stat(filepath.Join(rootPath, p))
			if s != nil {
				packages = append(packages, filepath.Join(packagePath, p))
			}
		}

		fname := inflect.Parameterize(packagePath) + "_loader.go"
		g := gentronics.New()
		g.Add(gentronics.NewFile(fname, cMain))
		err = g.Run(os.TempDir(), gentronics.Data{
			"packages": packages,
		})
		os.Chdir(rootPath)
		if err != nil {
			return err
		}

		cmd := exec.Command("gore", "-autoimport", "-context", filepath.Join(os.TempDir(), fname))
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()
	},
}

func init() {
	RootCmd.AddCommand(consoleCmd)
}

var cMain = `
package main

{{#each packages}}
import _ "{{.}}"
{{/each}}
`
