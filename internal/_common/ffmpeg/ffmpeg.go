package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"sync"

	"go-media-stream/internal/config"

	"github.com/sirupsen/logrus"
)

type FFMPEG struct {
	jobs chan JobModel
	wg   *sync.WaitGroup
}

func NewFFMPEG(conf *config.Config) *FFMPEG {
	jobs := make(chan JobModel, conf.FilmDecWrCount)
	var wg sync.WaitGroup
	return &FFMPEG{
		jobs: jobs,
		wg:   &wg,
	}
}

func (f *FFMPEG) Worker(result chan<- WorkerResult) {
	for j := range f.jobs {
		logrus.Info("Worker start")
		logrus.Info("ExecVideo")
		err := f.ExecVideo(j.Name, j.FileName)
		if err != nil {
			logrus.Error(err)
			break
		}
		logrus.Info("GetAudioList")
		audioList, err := f.GetAudioList(j.Name, j.FileName)
		if err != nil {
			logrus.Error(err)
			break
		}
		logrus.Info(len(audioList))
		audioListGood := make([]WorkerResultAudios, 0)
		for _, al := range audioList {
			fmt.Println(al.Index, al.Title)
			logrus.Info("ExecAudio")
			err := f.ExecAudio(j.Name, j.FileName, al.Index)
			if err != nil {
				logrus.Error(err)
				break
			}
			audioListGood = append(audioListGood, WorkerResultAudios{
				Title: al.Title,
				Index: al.Index,
			})
		}

		logrus.Info("Send result")
		result <- WorkerResult{
			Name:   j.Name,
			Audios: audioListGood,
		}
	}
	logrus.Info("Worker Done")
	f.wg.Done()
}

func (f *FFMPEG) CreateWorkers(count int, result chan<- WorkerResult) {
	for i := 0; i < count; i++ {
		f.wg.Add(1)
		logrus.Info("Create worker")
		go f.Worker(result)
	}
}

func (f *FFMPEG) Wait() {
	logrus.Info("Workers Wait")
	f.wg.Wait()
}

func (f *FFMPEG) Close() {
	close(f.jobs)
}

func (f *FFMPEG) Add(j JobModel) {
	f.jobs <- j
}

func (f *FFMPEG) GetAudioList(folder string, fileName string) ([]AudioIndexTitleModel, error) {
	jsonData, err := exec.Command(
		"ffprobe", "-v", "quiet", "-print_format", "json", "-show_format",
		"-show_streams", path.Join("uploads", path.Join(folder, fileName)),
	).CombinedOutput()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	ffprobeModel := FFprobeModel{}
	err = json.Unmarshal(jsonData, &ffprobeModel)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	out := make([]AudioIndexTitleModel, 0)
	for _, stream := range ffprobeModel.Stream {
		if stream.Codec_type == "audio" {
			title := stream.Tags.Title
			if title == "" {
				title = "default"
			}
			out = append(out, AudioIndexTitleModel{
				Index: stream.Index,
				Title: fmt.Sprintf("%s %s", title, stream.Tags.Language),
			})
		}
	}
	return out, nil
}

func (f *FFMPEG) ExecAudio(folder string, name string, index int16) error {
	input := path.Join("uploads", folder, name)
	output := path.Join("uploads", folder, fmt.Sprintf("%d.mp3", index))
	cmd := exec.Command("ffmpeg", "-y", "-i", input, "-map", fmt.Sprintf("0:%d", index), output)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (f *FFMPEG) ExecVideo(folder string, name string) error {
	input := path.Join("uploads", folder, name)
	output := path.Join("uploads", folder, "video.mp4")
	// ffmpeg -i TED.avi -c:v libx264 -crf 19 -preset slow -c:a aac -b:a 192k -ac 2 out.mp4
	// ffmpeg -i TED.avi -c:v libx264 -crf 18 -preset slow -c:a aac -b:a 192k output_video.mp4
	cmd := exec.Command("ffmpeg", "-i", input, "-c:v", "libx264", "-crf", "18", "-preset", "slow", output)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
