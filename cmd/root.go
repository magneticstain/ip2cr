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
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ip-2-cloudresource",
	Short: "a CLI tool for correlating a cloud IP address with its associated resources, with a focus on speed and ease-of-use.",
	Long: `IP-2-CloudResource (IP2CR) is a tool used for correlating a cloud IP address with its associated resources. It focuses on providing as much context to the user as possible, as fast as possible.
	
	Basic Usage:
	ip2cr -ipaddr=1.2.3.4
	
	Org Search
	ip2cr -ipaddr=1.2.3.4 -org-search -org-search-role-name=ip2cr-xaccount-role -org-search-role-name=arn:aws:iam::123456789012:role/org-manage -org-search-ou-id=ou-abcd-12345`,
	Version: "2.1.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
