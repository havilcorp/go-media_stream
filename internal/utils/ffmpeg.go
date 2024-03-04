package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"sync"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/log"

	"github.com/sirupsen/logrus"
)

type (
	tag struct {
		Language string `json:"language"`
		Title    string `json:"title,omitempty"`
	}
	stream struct {
		Index      int16  `json:"index"`
		Codec_name string `json:"codec_name"`
		Codec_type string `json:"codec_type"`
		Tags       tag    `json:"tags"`
	}
	format struct {
		Format_name string `json:"format_name"`
		Duration    string `json:"DURATION"`
		Size        string `json:"size"`
	}
	ffmpegStreams struct {
		Stream []stream `json:"streams"`
		Format format   `json:"format"`
	}
	audioIndexTitleModel struct {
		Index int16
		Title string
	}
)

type JobFFMPEG struct {
	QueueId int
}

type FFMPEG struct {
	log  *log.Logger
	db   *sql.DB
	jobs chan JobFFMPEG
	wg   *sync.WaitGroup
}

func NewFFMPEG(log *log.Logger, db *sql.DB, countJobs int) *FFMPEG {
	jobs := make(chan JobFFMPEG, countJobs)
	var wg sync.WaitGroup
	return &FFMPEG{
		log:  log,
		db:   db,
		jobs: jobs,
		wg:   &wg,
	}
}

func (f *FFMPEG) Worker() {
	for j := range f.jobs {
		logrus.Info("Worker start")
		logrus.Info(j.QueueId)
		row := f.db.QueryRow(`
			SELECT id, user_id, video_id, folder, title, idx, type FROM queue WHERE id = ?
		`, j.QueueId)
		if err := row.Err(); err != nil {
			f.log.Error(err)
			break
		}
		queue := domain.Queue{}
		err := row.Scan(&queue.ID, &queue.UserID, &queue.VideoID, &queue.Folder, &queue.Title, &queue.Idx, &queue.Type)
		if err != nil {
			f.log.Error(err)
			break
		}
		if queue.Type == "video" {
			err = f.ExecVideo(context.Background(), queue.Folder, "original")
			if err != nil {
				logrus.Error(err)
				break
			}
			fmt.Println(queue.ID, "VIDEO DONE")
			_, err = f.db.Exec("UPDATE queue SET is_done = 1 WHERE id = ?;", queue.ID)
			if err != nil {
				f.log.Error(err)
				break
			}
			resultVideo, err := f.db.Exec(`
				INSERT INTO video (id, user_id, name) 
				VALUES (NULL, ?, ?);
			`, queue.UserID, queue.Folder)
			if err != nil {
				logrus.Error(err)
				break
			}
			videoId, err := resultVideo.LastInsertId()
			if err != nil {
				logrus.Error(err)
				break
			}

			audioList, err := f.GetAudioList(context.Background(), queue.Folder, "original")
			if err != nil {
				logrus.Error(err)
				break
			}
			for _, al := range audioList {
				_, err = f.db.Exec(`
					INSERT INTO queue (id, user_id, video_id, folder, title, idx, type) 
					VALUES (NULL, ?, ?, ?, ?, ?, "audio");
				`, queue.UserID, videoId, queue.Folder, al.Title, al.Index)
				if err != nil {
					logrus.Error(err)
					break
				}
			}
		} else if queue.Type == "audio" {
			if !queue.Idx.Valid {
				f.log.Error(errors.New("WRONG_IDX"))
				break
			}
			err = f.ExecAudio(queue.Folder, "original", queue.Idx.Int16)
			if err != nil {
				logrus.Error(err)
				break
			}
			fmt.Println(queue.ID, "AUDIO DONE")
			_, err = f.db.Exec("UPDATE queue SET is_done = 1 WHERE id = ?;", queue.ID)
			if err != nil {
				f.log.Error(err)
				break
			}
			_, err = f.db.Exec(`
				INSERT INTO audio (id, name, idx, video_id) 
				VALUES (NULL, ?, ?, ?);
			`, queue.Title, queue.Idx.Int16, queue.VideoID)
			if err != nil {
				logrus.Error(err)
				break
			}
		} else {
			f.log.Error(errors.New("WRONG_TYPE"))
			break
		}
	}
	logrus.Info("Worker Done")
	f.wg.Done()
}

func (f *FFMPEG) CreateWorkers(count int) {
	for i := 0; i < count; i++ {
		f.wg.Add(1)
		logrus.Info("Create worker")
		go f.Worker()
	}
}

func (f *FFMPEG) Wait() {
	logrus.Info("Workers Wait")
	f.wg.Wait()
}

func (f *FFMPEG) Close() {
	close(f.jobs)
}

func (f *FFMPEG) Add(j JobFFMPEG) {
	fmt.Println("Job add", j.QueueId)
	f.jobs <- j
}

func (f *FFMPEG) ExecVideo(ctx context.Context, folder string, name string) error {
	input := path.Join("uploads", folder, name)
	output := path.Join("uploads", folder, "video.mp4")
	// ffmpeg -i TED.avi -c:v libx264 -crf 19 -preset slow -c:a aac -b:a 192k -ac 2 out.mp4
	// ffmpeg -i TED.avi -c:v libx264 -crf 18 -preset slow -c:a aac -b:a 192k output_video.mp4
	// cmd := exec.CommandContext(ctx, "ffmpeg", "-i", input, "-c:v", "libx264", "-crf", "20", output)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", input, "-c:v", "copy", output)
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

func (f *FFMPEG) GetAudioList(ctx context.Context, folder string, fileName string) ([]audioIndexTitleModel, error) {
	jsonData, err := exec.CommandContext(ctx,
		"ffprobe", "-v", "quiet", "-print_format", "json", "-show_format",
		"-show_streams", path.Join("uploads", path.Join(folder, fileName)),
	).CombinedOutput()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	ffmpegStreams := ffmpegStreams{}
	err = json.Unmarshal(jsonData, &ffmpegStreams)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	out := make([]audioIndexTitleModel, 0)
	for _, stream := range ffmpegStreams.Stream {
		if stream.Codec_type == "audio" {
			title := stream.Tags.Title
			if title == "" {
				title = "default"
			}
			out = append(out, audioIndexTitleModel{
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
