package utils

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func Confirm(prompt string) {
	confirm := promptui.Prompt{
		Label:     prompt,
		IsConfirm: true,
	}
	_, err := confirm.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Operation Canceled!")
		os.Exit(1)
	}
}
