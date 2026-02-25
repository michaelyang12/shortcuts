package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultFrameCount = 6

// ExtractFrames downloads a video from url and extracts evenly-spaced frames.
// Returns paths to the extracted JPEG frames and a cleanup function.
func ExtractFrames(url string) ([]string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "shortcuts-video-*")
	if err != nil {
		return nil, nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }

	videoPath := filepath.Join(tmpDir, "video.mp4")

	cmd := exec.Command("yt-dlp",
		"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
		"--merge-output-format", "mp4",
		"-o", videoPath,
		url,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("yt-dlp failed: %w\n%s", err, string(out))
	}

	totalFrames, err := getFrameCount(videoPath)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("getting frame count: %w", err)
	}

	step := totalFrames / defaultFrameCount
	if step < 1 {
		step = 1
	}

	framesDir := filepath.Join(tmpDir, "frames")
	if err := os.MkdirAll(framesDir, 0755); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("creating frames dir: %w", err)
	}

	selectFilter := fmt.Sprintf("select='not(mod(n\\,%d))',scale=1280:-1", step)
	cmd = exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", selectFilter,
		"-vsync", "vfr",
		"-frames:v", strconv.Itoa(defaultFrameCount),
		"-q:v", "2",
		filepath.Join(framesDir, "frame_%04d.jpg"),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("ffmpeg failed: %w\n%s", err, string(out))
	}

	entries, err := os.ReadDir(framesDir)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("reading frames dir: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jpg") {
			paths = append(paths, filepath.Join(framesDir, e.Name()))
		}
	}

	if len(paths) == 0 {
		cleanup()
		return nil, nil, fmt.Errorf("no frames extracted")
	}

	return paths, cleanup, nil
}

func getFrameCount(videoPath string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-count_frames",
		"-show_entries", "stream=nb_read_frames",
		"-of", "default=nokey=1:noprint_wrappers=1",
		videoPath,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(out)))
}
