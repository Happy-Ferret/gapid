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

// +build linux darwin windows

package stub_test

import (
	"testing"
	"time"

	"github.com/google/gapid/core/context/jot"
	"github.com/google/gapid/core/event/task"
	"github.com/google/gapid/core/fault/cause"
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/os/shell"
	"github.com/google/gapid/core/os/shell/stub"
)

// TestNothing exists because some tools don't like packages that don't have any tests in them.
func TestNothing(t *testing.T) {}

// This example shows how to use a stub command target to fake an starting error.
func ExampleAlways() {
	ctx := log.Background().PreFilter(log.NoLimit).Filter(log.Pass).Handler(log.Stdout(log.Normal)).Enter("Example")
	s, err := shell.Command("echo", "Hello from the shell").On(stub.Respond("Hello")).Call(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say hello")
	}
	ctx.Print(s)
	ctx.Print("Done")
	// Output:
	//Info:Example:Hello
	//Info:Example:Done
}

// This example shows how to use simple exact matches with a fallback to the command echo.
func ExampleHello() {
	ctx := log.Background().PreFilter(log.NoLimit).Filter(log.Pass).Handler(log.Stdout(log.Normal)).Enter("Example")
	target := stub.OneOf(
		stub.RespondTo(`echo "Hello from the shell"`, "Nice to meet you"),
		stub.Echo{},
	)
	s, err := shell.Command("echo", "Hello from the shell").On(target).Call(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say hello")
	}
	ctx.Print(s)
	s, err = shell.Command("echo", "Goodbye now").On(target).Call(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say goodbye")
	}
	ctx.Print(s)
	ctx.Print("Done")
	// Output:
	//Info:Example:Nice to meet you
	//Info:Example:Goodbye now
	//Info:Example:Done
}

// This example shows how to use simple exact matches with a fallback to the command echo.
func ExampleRegexp() {
	ctx := log.Background().PreFilter(log.NoLimit).Filter(log.Pass).Handler(log.Stdout(log.Normal)).Enter("Example")
	target := stub.OneOf(
		stub.Regex(`smalltalk`, stub.Respond("Nice to meet you")),
		stub.Echo{},
	)
	s, _ := shell.Command("echo", "Hello").On(target).Call(ctx)
	ctx.Print(s)
	s, _ = shell.Command("echo", "Insert smalltalk here").On(target).Call(ctx)
	ctx.Print(s)
	s, _ = shell.Command("echo", "Goodbye").On(target).Call(ctx)
	ctx.Print(s)
	ctx.Print("Done")
	// Output:
	//Info:Example:Hello
	//Info:Example:Nice to meet you
	//Info:Example:Goodbye
	//Info:Example:Done
}

// This example shows how to use a stub command target to fake an starting error.
func ExampleStartError() {
	ctx := log.Background().PreFilter(log.NoLimit).Filter(log.Pass).Handler(log.Stdout(log.Normal))
	ctx = ctx.Enter("Example")
	target := stub.OneOf(
		stub.Match(`echo Goodbye`, &stub.Response{StartErr: cause.Explain(ctx, nil, "bad command")}),
		stub.Echo{},
	)
	output, err := shell.Command("echo", "Hello").On(target).Call(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say hello")
	}
	ctx.Info().Log(output)
	err = shell.Command("echo", "Goodbye").On(target).Run(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say goodbye")
	}
	ctx.Print("Done")
	// Output:
	//Info:Example:Hello
	//Error:Example:Unable to say goodbye:{Example:Start:{Example:⦕bad command⦖}:Command=echo Goodbye,On=stub}
	//Info:Example:Done
}

// This example shows how to use a stub command target to fake an blocking process, and cancelling it.
func ExampleBlockingCancel() {
	ctx := log.Background().PreFilter(log.NoLimit).Filter(log.Pass).Handler(log.Stdout(log.Normal))
	ctx, cancel := task.WithCancel(ctx.Enter("Example"))
	go func() {
		time.Sleep(time.Millisecond)
		cancel()
	}()
	response := &stub.Response{WaitErr: cause.Explain(ctx, nil, "Cancelled")}
	response.WaitSignal, response.KillTask = task.NewSignal()
	target := stub.OneOf(stub.Match(`echo "Hello from the shell"`, response))
	err := shell.Command("echo", "Hello from the shell").On(target).Run(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say hello")
	}
	ctx.Print("Done")
	// Output:
	//Error:Example:Unable to say hello:{Example:Wait:{Example:⦕Cancelled⦖}:Command=echo "Hello from the shell",On=stub}
	//Info:Example:Done
}

// This example shows what happens if you attempt a command that has no response.
func ExampleUnmatchedCommand() {
	ctx := log.Background().PreFilter(log.NoLimit).Filter(log.Pass).Handler(log.Stdout(log.Normal))
	ctx = ctx.Enter("Example")
	target := stub.OneOf(
		stub.RespondTo(`echo "Hello from the shell"`, "Nice to meet you"),
	)
	err := shell.Command("echo", "Not hello from the shell").On(target).Run(ctx)
	if err != nil {
		jot.Fail(ctx, err, "Unable to say hello")
	}
	ctx.Print("Done")
	// Output:
	//Error:Example:Unable to say hello:{Example:Start:⦕unmatched:echo "Not hello from the shell"⦖:Command=echo "Not hello from the shell",On=stub}
	//Info:Example:Done
}
