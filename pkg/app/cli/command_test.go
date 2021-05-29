package cli_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"testing"
)

func TestCommand(t *testing.T) {

	t.Run("should add child command and found the child", func(t *testing.T) {
		barCmd := newBarCommand()
		bazCmd := newBazCommand()
		fooCmd := newFooCommand(barCmd, bazCmd)
		fooCmd.SetName("foo")
		assert.Equal(t, "foo", fooCmd.FullName())

		fooCmd.SetFullName("foo command")

		assert.Equal(t, "foo command", fooCmd.FullName())

		assert.Equal(t, true, fooCmd.HasChild())
		assert.Equal(t, barCmd, fooCmd.Children()[0])

		assert.Equal(t, fooCmd.GetName(), barCmd.Parent().GetName())

		foundCmd, err := fooCmd.Find("bar")
		assert.Equal(t, nil, err)
		assert.Equal(t, foundCmd.GetName(), barCmd.GetName())
	})

	t.Run("should found the command directly", func(t *testing.T) {
		fooCmd := new(fooCommand)
		fooCmd.SetName("foo")
		_, err := fooCmd.Find("foo")
		assert.Equal(t, nil, err)
	})

	t.Run("should not found the non-exist command", func(t *testing.T) {
		fooCmd := new(fooCommand)
		fooCmd.SetName("foo")
		_, err := fooCmd.Find("bar")
		assert.Equal(t, cli.ErrCommandNotFound, err)
	})

	t.Run("should run command handler", func(t *testing.T) {
		cmd := new(cli.SubCommand)
		err := cmd.Run(nil)
		assert.Equal(t, nil, err)
	})

	t.Run("should run secondary command handler", func(t *testing.T) {
		cmd := new(fooCommand)

		res, err := cli.Dispatch(cmd, []string{"daz"})
		assert.Equal(t, nil, err)
		assert.Equal(t, false, res.(bool))

		res, err = cli.Dispatch(cmd, []string{"buzz"})
		assert.Equal(t, nil, err)
		assert.Equal(t, true, res.(bool))
	})
}
