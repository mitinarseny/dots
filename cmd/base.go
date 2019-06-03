package cmd

import "fmt"

type SymlinkTargetExistsError struct {
	string
}

func (err SymlinkTargetExistsError) Error() string {
	return fmt.Sprintf("Symlink target already exists: '%v'", err.string)
}
