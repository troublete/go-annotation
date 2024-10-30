package auxiliary

import (
	"golang.org/x/tools/imports"
	"os"
)

func WriteGoFile(path string, content []byte) error {
	err := os.WriteFile(path, content, os.ModeTemporary)
	if err != nil {
		return err
	}

	out, err := imports.Process(path, content, &imports.Options{
		Fragment:  true,
		AllErrors: true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(path, out, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
