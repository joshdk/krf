// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"time"

	shlex "github.com/carapace-sh/carapace-shlex"

	"github.com/joshdk/krf/resources"
)

// NewExecMatcher matches resources.Resource instances based on the results of
// executing an external program.
func NewExecMatcher(command string) (Matcher, error) {
	tokens, err := shlex.Split(command)
	if err != nil {
		return nil, err
	}

	args := tokens.Strings()

	return execMatcher{args[0], args[1:]}, nil
}

type execMatcher struct {
	command   string
	arguments []string
}

func (m execMatcher) Matches(item resources.Resource) bool {
	// Marshal the current item into a single line of JSON. Scripts can then
	// make a single `read` call to consume this from stdin.
	data, err := item.MarshalJSON()
	if err != nil {
		return false
	}

	// Restrict command execution to a conservative 2 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, m.command, m.arguments...)
	cmd.Stdin = bytes.NewBuffer(data)

	// Configure several helpful environment variable values into executed
	// process. Can be used for e.g. short-circuiting to avoid parsing the
	// input JSON.
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"RESOURCE_APIVERSION="+item.GetAPIVersion(),
		"RESOURCE_KIND="+item.GetKind(),
		"RESOURCE_NAME="+item.GetName(),
		"RESOURCE_NAMESPACE="+item.GetNamespace(),
	)

	// Run the program and match the input resource only if the program
	// finishes with a 0 exit status. Any non-0 exit status or errors will
	// prevent matching of the input resource.
	return cmd.Run() == nil
}
