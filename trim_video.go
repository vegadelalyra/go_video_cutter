package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func main() {
	// Input file and desired output
	inputFile := "2025-01-22 14-37-32.mkv"
	outputFile := "go 2025-01-22 14-37-32.mkv"

	// Get the number of CPU cores
	numChunks := runtime.NumCPU()
	log.Printf("Detected %d CPU cores. Processing video in %d chunks...\n", numChunks, numChunks)

	// Get video duration
	duration, err := getVideoDuration(inputFile)
	if err != nil {
		log.Fatalf("Failed to get video duration: %v", err)
	}

	// Calculate chunk size
	chunkSize := int(duration / float64(numChunks))
	var wg sync.WaitGroup

	// Process each chunk in parallel
	for i := 0; i < numChunks; i++ {
		wg.Add(1)

		go func(chunkIndex int) {
			defer wg.Done()

			startTime := chunkIndex * chunkSize
			outputChunk := fmt.Sprintf("chunk_%d.mp4", chunkIndex)

			err := trimVideo(inputFile, outputChunk, startTime, chunkSize)
			if err != nil {
				log.Printf("Failed to process chunk %d: %v", chunkIndex, err)
			}
		}(i)
	}

	// Wait for all chunks to finish processing
	wg.Wait()

	// Merge chunks into a single output video
	err = mergeChunks(numChunks, "merged_output.mp4")
	if err != nil {
		log.Fatalf("Failed to merge chunks: %v", err)
	}

	// Set the final duration you want for the output video (in seconds)
	finalDuration := float64(6720) // Change this value to your desired final duration in seconds

	// Trim the final merged video based on finalDuration
	err = trimFinalVideo("merged_output.mp4", outputFile, finalDuration)
	if err != nil {
		log.Fatalf("Failed to trim final video: %v", err)
	}

	log.Println("Video processed and trimmed successfully!")
}

// getVideoDuration retrieves the total duration of the video in seconds
func getVideoDuration(inputFile string) (float64, error) {
	// Convert backslashes to forward slashes
	inputFile = strings.ReplaceAll(inputFile, "\\", "/")
	
	cmd := exec.Command("ffprobe", "-i", inputFile, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running ffprobe: %v\n", err)
		log.Printf("ffprobe output:\n%s", output)
		return 0, err
	}

	var duration float64
	_, err = fmt.Sscanf(string(output), "%f", &duration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %v", err)
	}
	return duration, nil
}

// trimVideo cuts a portion of the video
func trimVideo(inputFile, outputFile string, startTime, duration int) error {
	start := fmt.Sprintf("%02d:%02d:%02d", startTime/3600, (startTime%3600)/60, startTime%60)
	dur := fmt.Sprintf("%d", duration)

	// Run the ffmpeg command to trim the video
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", start, "-t", dur, "-c", "copy", outputFile)
	return cmd.Run()
}

// mergeChunks merges all processed chunks into a single video
func mergeChunks(numChunks int, outputFile string) error {
	// Create a text file with the list of chunk files
	fileList := "file_list.txt"
	file, err := os.Create(fileList)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < numChunks; i++ {
		_, err := file.WriteString(fmt.Sprintf("file 'chunk_%d.mp4'\n", i))
		if err != nil {
			return err
		}
	}

	// Run FFmpeg to concatenate the chunks
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", fileList, "-c", "copy", outputFile)
	return cmd.Run()
}

// trimFinalVideo trims the final merged video to the desired final duration
func trimFinalVideo(inputFile, outputFile string, finalDuration float64) error {
	// Format the final duration to a string (in seconds)
	duration := fmt.Sprintf("%f", finalDuration)

	// Run the ffmpeg command to trim the video
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-t", duration, "-c", "copy", outputFile)
	return cmd.Run()
}
