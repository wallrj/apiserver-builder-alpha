/*
Copyright 2017 The Kubernetes Authors.

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

package create

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/markbates/inflect"
	"github.com/spf13/cobra"

	"sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot/boot/util"
)

func ValidateResourceFlags() {
	util.GetDomain()
	if len(groupName) == 0 {
		log.Fatalf("Must specify --group")
	}
	if len(versionName) == 0 {
		log.Fatalf("Must specify --version")
	}
	if len(kindName) == 0 {
		log.Fatal("Must specify --kind")
	}
	if len(resourceName) == 0 {
		resourceName = inflect.NewDefaultRuleset().Pluralize(strings.ToLower(kindName))
	}

	groupMatch := regexp.MustCompile("^[a-z]+$")
	if !groupMatch.MatchString(groupName) {
		log.Fatalf("--group must match regex ^[a-z]+$ but was (%s)", groupName)
	}
	versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)*$")
	if !versionMatch.MatchString(versionName) {
		log.Fatalf(
			"--version has bad format. must match ^v\\d+(alpha\\d+|beta\\d+)*$.  "+
				"e.g. v1alpha1,v1beta1,v1 but was (%s)", versionName)
	}

	kindMatch := regexp.MustCompile("^[A-Z]+[A-Za-z0-9]*$")
	if !kindMatch.MatchString(kindName) {
		log.Fatalf("--kind must match regex ^[A-Z]+[A-Za-z0-9]*$ but was (%s)", kindName)
	}
}

func RegisterResourceFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&groupName, "group", "", "name of the API group excluding its domain name. i.e. package name"+
		"  **Must be single lowercase word (match ^[a-z]+$)**.")
	cmd.Flags().StringVar(&versionName, "version", "", "name of the API version.  **must match regex v\\d+(alpha\\d+|beta\\d+)** e.g. v1, v1beta1, v1alpha1")
	cmd.Flags().StringVar(&kindName, "kind", "", "name of the API kind.  **Must be CamelCased (match ^[A-Z]+[A-Za-z0-9]*$)**")
	cmd.Flags().StringVar(&resourceName, "resource", "", "optional name of the API resource, defaults to the plural name of the lowercase kind")
}

// Yesno reads from stdin looking for one of "y", "yes", "n", "no" and returns
// true for "y" and false for "n"
func Yesno(reader *bufio.Reader) bool {
	for {
		text := readstdin(reader)
		switch text {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Printf("invalid input %q, should be [y/n]", text)
		}
	}
}

// Readstdin reads a line from stdin trimming spaces, and returns the value.
// log.Fatal's if there is an error.
func readstdin(reader *bufio.Reader) string {
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error when reading input: %v", err)
	}
	return strings.TrimSpace(text)
}