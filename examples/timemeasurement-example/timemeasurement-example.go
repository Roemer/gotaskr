package main

import (
	"os"
	"time"

	"github.com/roemer/gotaskr"
)

func init() {
	gotaskr.Task("Time-Measurement-Task", func() error {
		// Using with a func
		if err := gotaskr.MeasureTime("Sub-Item-1", func() error {
			time.Sleep(time.Second * 1)
			return nil
		}); err != nil {
			return err
		}

		// Using with an object
		measurement := gotaskr.StartTimeMeasurement("Sub-Item-2")
		time.Sleep(time.Millisecond * 500)
		measurement.Finish()

		// Nesting
		measurement1 := gotaskr.StartTimeMeasurement("Sub-Item-3")
		time.Sleep(time.Millisecond * 500)
		measurement2 := gotaskr.StartTimeMeasurement("Sub-Item-4")
		time.Sleep(time.Millisecond * 500)
		measurement2.Finish()
		time.Sleep(time.Millisecond * 500)
		measurement1.Finish()

		return nil
	})
}

func main() {
	os.Exit(gotaskr.Execute())
}
