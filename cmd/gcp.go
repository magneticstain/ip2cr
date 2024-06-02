/*
Copyright © 2023 Josh Carlson <magneticstain@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// gcpCmd represents the gcp command
var (
	// FLAGS
	ProjectID string

	// CMDS
	gcpCmd = &cobra.Command{
		Use:   "gcp",
		Short: "Search for IP in GCP",
		Long:  `Search for a given IP address within one or more Google Cloud Platform subscriptions.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gcp called")
		},
	}
)

func init() {
	gcpCmd.PersistentFlags().StringVarP(
		&ProjectID,
		"project-id",
		"",
		"",
		"The ID of the GCP project to target",
	)
	gcpCmd.MarkFlagRequired("project-id")

	rootCmd.AddCommand(gcpCmd)
}
