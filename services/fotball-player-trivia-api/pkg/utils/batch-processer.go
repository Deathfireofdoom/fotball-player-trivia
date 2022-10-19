package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func StartBatchProcessExample() {
	filePath := "/Users/oskarelvkull/Documents/big-corp/fotball-player-trivia/services/fotball-trivia-api/database/data/player-data-set.csv"
	BatchProcessFile(filePath, ExampleBatchProcess, 1, 10)

}

// BatchProcessFile takes a file and a function to apply on each line.
func BatchProcessFile(filePath string, batchProcess func([]string), concurrency, batchSize int) {
	// Cancel chanel, used to communicate that the file has been read successfuly.
	// The channel needs to have 1 buffer since last signal will not be read.
	cancelCh := make(chan bool, 1)
	// Chancel to send the batches to other go-routines. The channel will have 5 messages in the buffer.
	batchCh := make(chan []string, concurrency+5)

	// Using a wait group so we can gracefully exit without loosing any batches, or atleast not loosing any batch due to premature exit heh.
	wg := new(sync.WaitGroup)

	// Starts the concurrency.
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		go consumer(wg, batchProcess, batchCh, cancelCh)
	}

	// Creates a file reader that later will be buff read.
	file, err := os.Open(filePath)
	if err != nil {
		panic("Could not read file.")
	}
	defer file.Close()

	// Buffering of batch processing.
	scanner := bufio.NewScanner(file)
	batch := []string{}

	for scanner.Scan() {
		// Checks if batch-size is met. If so, batch it publish to chanel where the go-routines listens to.
		if len(batch) >= batchSize {
			batchCh <- batch

			// Resets batch.
			batch = []string{}
		}
		// Add line to batch.
		batch = append(batch, scanner.Text())
	}
	// Publish last batch even though max-size is not met.
	batchCh <- batch
	// Publishing a nil to communicate that the last message has been read.
	batchCh <- nil

	// Gotta figure out why we need to call wg.Done() one extra time..
	fmt.Println("LOG - Waiting for processes to finish.")
	//wg.Done()
	wg.Wait()
	// Maybe should delete the cancelCancel? Not sure if that is needed, a true value will be stuck there.
	fmt.Println("LOG - All processes finished, successfully processed the full file-.")

}

// consumer is the Go rountine waiting for batches in the batchCh, when a "nil" is sent in the batchCh
// the last line has been consumed(or is being consumed.). A exit signal in sent in cancelCh and the
// consumer exits.
func consumer(wg *sync.WaitGroup, batchProcess func([]string), batchCh chan []string, cancelCh chan bool) {
	defer wg.Done()
	for {
		select {
		case batch := <-batchCh:
			fmt.Println("Processing")
			if batch != nil {
				batchProcess(batch)
			} else {
				cancelCh <- true
				return
			}

		case signal := <-cancelCh:
			cancelCh <- signal
			return
		}
	}
}

// ExampleBatchProcess is just an example how a batchProcess function could look like.
func ExampleBatchProcess(batch []string) {
	for _, line := range batch {
		values := ParseLine(line, ",")
		fmt.Println(values)
	}
}

func ParseLine(line, delimiter string) []string {
	return strings.Split(line, delimiter)
}
