# cli application examples

## Integrate with promptui 

It's easy to integrate promptui into Hiboot cli application

```go
// declare main package
package main

// import cli starter and fmt
import (
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/manifoldco/promptui"
	"fmt"
)

// define the command
type PromptuiCommand struct {
	cli.RootCommand
}

// Init constructor
func (c *PromptuiCommand) Init() {
	c.Use = "promptui"
	c.Short = "promptui command"
	c.Long = "run promptui command for getting started"

}

// Run run the command
func (c *PromptuiCommand) Run(args []string) error {
	// define selections
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
	}

    // call prompt.Run() to get users selection
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return fmt.Errorf("Prompt failed %v\n", err)
	}

	fmt.Printf("You choose %q\n", result)

	return nil
}

// main function
func main() {
	// create new cli application and run it
	cli.NewApplication(new(PromptuiCommand)).Run()
}

```