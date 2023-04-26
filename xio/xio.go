// Package xio implements helper and utility functions for I/O.
package xio

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies the contents of the file named src to the file named by dst.
// If the destination file exists, all it's contents will be replaced by the
// contents of the source file.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("error opening destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	if err = dstFile.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %w", err)
	}

	return nil
}
