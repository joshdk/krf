// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"errors"
	"strconv"
	"time"

	"github.com/joshdk/krf/resources"
)

// parseTime parses the provided time spec and returns a function that can be
// used to compare against that point in time. The provided time spec can
// reference either an absolute point in time, or a duration in the past
// relative to the compared time.
func parseTime(spec string) (func(time.Time) int, error) {
	// Attempt to parse as a duration.
	if duration, err := time.ParseDuration(spec); err == nil {
		return func(t time.Time) int {
			return time.Now().Add(-1 * duration).Compare(t)
		}, nil
	}

	// Attempt to parse as an integer, and treat it as a Unix time.
	if unix, err := strconv.ParseInt(spec, 10, 64); err == nil {
		return time.Unix(unix, 0).Compare, nil
	}

	layouts := []string{
		time.RFC3339,
		time.UnixDate,
	}

	// Attempt to parse as any of the above time layouts.
	for _, layout := range layouts {
		if timestamp, err := time.Parse(layout, spec); err == nil {
			return timestamp.Compare, nil
		}
	}

	return nil, errors.New("could not parse timestamp or duration")
}

// NewCreatedAfterMatcher matches resources.Resource instances which have a
// creationTimestamp which is newer than (was created after) the point in time
// described by the provided spec.
func NewCreatedAfterMatcher(spec string) (Matcher, error) {
	compareFn, err := parseTime(spec)
	if err != nil {
		return nil, err
	}

	return &creationTimestampMatcher{
		compareFn:   compareFn,
		checkBefore: false,
	}, nil
}

// NewCreatedBeforeMatcher matches resources.Resource instances which have a
// creationTimestamp which is older than (was created before) the point in time
// described by the provided spec.
func NewCreatedBeforeMatcher(spec string) (Matcher, error) {
	compareFn, err := parseTime(spec)
	if err != nil {
		return nil, err
	}

	return &creationTimestampMatcher{
		compareFn:   compareFn,
		checkBefore: true,
	}, nil
}

type creationTimestampMatcher struct {
	compareFn   func(time.Time) int
	checkBefore bool
}

func (m creationTimestampMatcher) Matches(item resources.Resource) bool {
	creationTimestamp := item.GetCreationTimestamp().Time

	// Objects that do not have a creationTimestamp are never matched.
	if creationTimestamp.IsZero() {
		return false
	}

	if m.checkBefore {
		// Return true if the item's creationTimestamp is older than the matcher time.
		return m.compareFn(creationTimestamp) > 0
	}

	// Return true if the item's creationTimestamp is newer than the matcher time.
	return m.compareFn(creationTimestamp) < 0
}
