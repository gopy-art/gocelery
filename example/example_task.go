/*
This is an example of celery worker for rabbitmq and redis implementation
in this example we just declare requirements for celery workers
*/

package example

import "fmt"

// exampleAddTask is integer addition task
// with named arguments
type ExampleAddTask struct {
	TaskID string
	a      int
	b      int
}

// this function is for reading the argument that has been passed to the message from 'Kwargs'
func (a *ExampleAddTask) ParseKwargs(kwargs map[string]interface{}) error {
	kwargA, ok := kwargs["a"]
	if !ok {
		return fmt.Errorf("undefined kwarg a")
	}
	kwargAFloat, ok := kwargA.(float64)
	if !ok {
		return fmt.Errorf("malformed kwarg a")
	}
	a.a = int(kwargAFloat)
	kwargB, ok := kwargs["b"]
	if !ok {
		return fmt.Errorf("undefined kwarg b")
	}
	kwargBFloat, ok := kwargB.(float64)
	if !ok {
		return fmt.Errorf("malformed kwarg b")
	}
	a.b = int(kwargBFloat)
	return nil
}

func (a *ExampleAddTask) ParseId(id string) error {
	a.TaskID = id
	return nil
}

// The main function that will be execute
func (a *ExampleAddTask) RunTask() (interface{}, error) {
	result := a.a + a.b
	fmt.Printf("Task with uuid %v has result %v \n", a.TaskID, result)
	return result, nil
}
