package image

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/tiagorlampert/CHAOS/internal"
)

// WritePNG write a given content to a PNG file
func WritePNG(content []byte) (string, error) {
	filename := fmt.Sprint(uuid.New().String(), ".png")
	file, err := os.OpenFile(
		fmt.Sprint(internal.TempDirectory, string(os.PathSeparator), filename),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		os.ModePerm,
	)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return "", err
	}
	return filename, nil
}
