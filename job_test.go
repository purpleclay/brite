/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package brite_test

import (
	"context"
	"errors"
	"testing"

	"github.com/purpleclay/brite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestString(t *testing.T) {
	name := "this is a test"
	job := brite.NewJob(name)

	assert.Equal(t, name, job.String())
}

func TestRunSynchronousTask(t *testing.T) {
	m := mockTask(t)
	m.On("Skip", mock.Anything).Return(false)
	m.On("Run", mock.Anything).Return(nil)

	job := brite.NewJob("job 1")
	job.Register(m)
	err := job.Run(context.Background())

	assert.NoError(t, err)
}

func TestRunSynchronousTaskError(t *testing.T) {
	m := mockTask(t)
	m.On("Skip", mock.Anything).Return(false)
	m.On("Run", mock.Anything).Return(errors.New("task error"))

	job := brite.NewJob("job 2")
	job.Register(m)
	err := job.Run(context.Background())

	assert.EqualError(t, err, "task error")
}

func TestSkip(t *testing.T) {
	m := mockTask(t)
	m.On("Skip", mock.Anything).Return(true)

	job := brite.NewJob("job 3")
	job.Register(m)
	job.Run(context.Background())

	m.AssertNotCalled(t, "Run")
}

type runnerMock struct {
	mock.Mock
}

func (m *runnerMock) Run(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *runnerMock) Skip(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *runnerMock) String() string {
	return "mock"
}

func mockTask(t testing.TB) *runnerMock {
	t.Helper()

	mock := &runnerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() {
		mock.AssertExpectations(t)
	})

	return mock
}
