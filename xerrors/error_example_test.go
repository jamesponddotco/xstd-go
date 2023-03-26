package xerrors_test

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

func Example() {
	const errOpenFile xerrors.Error = "failed to open file"

	if _, err := os.Open("file-does-not-exist.txt"); err != nil {
		fmt.Println(fmt.Errorf("%w: %w", errOpenFile, err))
	}
	// Output: failed to open file: open file-does-not-exist.txt: no such file or directory
}
