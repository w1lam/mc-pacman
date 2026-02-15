package progress

import "fmt"

type CLIProgress struct{}

func (c *CLIProgress) Emit(event ProgressEvent) {
	switch event.Type {
	case ProgressStart:
		fmt.Printf("[%s] %s started\n", event.Context, event.PackageID)
	case ProgressUpdate:
		fmt.Printf("[%s] %s %.0f%% %s\n", event.Context, event.PackageID, event.Percentage, event.Message)
	case ProgressSuccess:
		fmt.Printf("[%s] %s finished\n", event.Context, event.PackageID)
	case ProgressFailure:
		fmt.Printf("[%s] %s failed\n", event.Context, event.PackageID)
	}
}
