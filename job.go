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

package brite

import "context"

// Job defines a way of registering a series of tasks and executing them
// in registration order.
type Job struct {
	name  string
	tasks []Runner
}

// NewJob creates a new named job.
func NewJob(name string) *Job {
	return &Job{name: name}
}

// Register a new synchronous task with this job. Registration order
// is guaranteed to be preserved.
func (j *Job) Register(task Runner) {
	j.tasks = append(j.tasks, task)
}

// Run the job. A job can consist of any number of tasks, each being
// executed in the order they are registered. Execution of the job
// is halted upon the first error encountered.
func (j *Job) Run(ctx context.Context) error {
	for _, task := range j.tasks {
		if task.Skip(ctx) {
			continue
		}

		if err := task.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

// String returns the name of this job.
func (j *Job) String() string {
	return j.name
}
