// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"context"
	"os"

	"github.com/open-policy-agent/opa/v1/rego"

	"github.com/joshdk/krf/resources"
)

var query = rego.Query(`data.krf.joshdk.github.com.matched`)

// NewRegoMatcher matches resources.Resource instances based on an evaluation
// of the given rego file.
//
// The rego file must have a package value of `krf.joshdk.github.com` and
// provide a rule named `matches` in order to match a resource.
func NewRegoMatcher(filename string) (Matcher, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	prepared, err := rego.New(
		query,
		rego.Module(filename, string(data)),
	).PrepareForEval(context.Background())
	if err != nil {
		return nil, err
	}

	return regoMatcher{preparedQuery: prepared}, nil
}

type regoMatcher struct {
	preparedQuery rego.PreparedEvalQuery
}

func (m regoMatcher) Matches(item resources.Resource) bool {
	result, err := m.preparedQuery.Eval(context.Background(), rego.EvalInput(item.Object))
	if err != nil {
		return false
	}

	return result.Allowed()
}
