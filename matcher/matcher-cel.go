// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"errors"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"

	"github.com/joshdk/krf/resources"
)

// NewCELMatcher matches resources.Resource instances based on the results of
// evaluating a boolean CEL expression.
func NewCELMatcher(expression string) (Matcher, error) {
	// Prepare the CEL environment to expect a single variable called "object".
	env, err := cel.NewEnv(
		cel.Variable("object",
			cel.MapType(cel.StringType, cel.AnyType),
		),
	)
	if err != nil {
		return nil, err
	}

	// Compile the user-supplied CEL expression.
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Require that the output type of the CEL expression is a boolean.
	if !ast.OutputType().IsExactType(types.BoolType) {
		return nil, errors.New("cel expression does not have a bool output type")
	}

	// Prepare a bundled CEL program for evaluation.
	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	return celMatcher{program: program}, nil
}

type celMatcher struct {
	program cel.Program
}

func (m celMatcher) Matches(item resources.Resource) bool {
	// Evaluate the CEL program against the current resource.
	result, _, err := m.program.Eval(map[string]any{
		"object": item.Object,
	})
	if err != nil {
		return false
	}

	if boolResult, ok := result.Value().(bool); ok {
		return boolResult
	}

	return false
}
