// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"strings"

	"github.com/google/gapid/core/log"
)

// CurrentBranch returns the current active git branch.
// If git is in a detached head state, it may return HEAD
func (g Git) CurrentBranch(ctx log.Context) (string, error) {
	str, _, err := g.run(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(str), err
}

// Rebase performs a `git rebase` on to the target branch.
func (g Git) Rebase(ctx log.Context, targetBranch string) error {
	_, _, err := g.run(ctx, "rebase", targetBranch)
	return err
}
