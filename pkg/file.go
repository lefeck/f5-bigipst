package pkg

import (
	"bufio"
	"os"
)

func WriteFile(data string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	buf.WriteString(data)
	if err = buf.Flush(); err != nil {
		return err
	}
	return nil
}
