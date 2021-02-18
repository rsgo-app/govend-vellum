//  Copyright (c) 2017 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/spf13/cobra"
)

var (
	expvarBind   string
	memprofile   string
	cpuprofile   string
	traceprofile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vellum",
	Short: "A utility to work with vellum FST files",
	Long:  `A utility to work with vellum FST files.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if expvarBind != "" {
			go http.ListenAndServe(expvarBind, nil)
		}
		if cpuprofile != "" {
			f, err := os.Create(cpuprofile)
			if err != nil {
				return err
			}
			pprof.StartCPUProfile(f)
		}
		if traceprofile != "" {
			f, err := os.Create(traceprofile)
			if err != nil {
				return err
			}
			if err := trace.Start(f); err != nil {
				return err
			}
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if cpuprofile != "" {
			pprof.StopCPUProfile()
		}
		if memprofile != "" {
			f, err := os.Create(memprofile)
			if err != nil {
				return err
			}
			if err := pprof.WriteHeapProfile(f); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
		if traceprofile != "" {
			trace.Stop()
		}
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&expvarBind, "expvar", "", "bind address for expvar, default none")
	RootCmd.PersistentFlags().StringVar(&cpuprofile, "cpuprofile", "", "cpuprofile output file, default none")
	RootCmd.PersistentFlags().StringVar(&memprofile, "memprofile", "", "memprofile output file, default none")
	RootCmd.PersistentFlags().StringVar(&traceprofile, "traceprofile", "", "traceprofile output file, default none")
}
