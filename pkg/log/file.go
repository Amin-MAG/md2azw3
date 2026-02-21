package ravandlog

import (
	"fmt"
	"math"
)

// OutputFileConfig represents details about output logs in files
type OutputFileConfig struct {
	Path             string
	Name             string
	MaxSizeInBytes   int64
	MaxNumberOfFiles int
}

func (ofc OutputFileConfig) FullPath(fileIndex int) string {
	if ofc.Path == "" {
		ofc.Path = "."
	}
	if ofc.Name == "" {
		ofc.Name = "vcloud"
	}

	// Calculate number of digits
	numberOfDigits := 1
	if ofc.MaxNumberOfFiles > 1 {
		numberOfDigits = int(math.Log10(float64(ofc.MaxNumberOfFiles)-1) + 1)
	}

	return fmt.Sprintf("%s/%s-%0*d.log", ofc.Path, ofc.Name, numberOfDigits, fileIndex)
}
