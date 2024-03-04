package ffmpeg

// JobModel

type JobModel struct {
	Name     string
	FileName string
}

// AudioIndexTitleModel

type AudioIndexTitleModel struct {
	Index int16
	Title string
}

// WorkerResult

type WorkerResultAudios struct {
	Title string
	Index int16
}

type WorkerResult struct {
	Name   string
	Audios []WorkerResultAudios
}

// FFprobeModel

type FFprobeModel struct {
	Stream []stream `json:"streams"`
	Format format   `json:"format"`
}

type tag struct {
	Language string `json:"language"`
	Title    string `json:"title,omitempty"`
}

type stream struct {
	Index      int16  `json:"index"`
	Codec_name string `json:"codec_name"`
	Codec_type string `json:"codec_type"`
	Tags       tag    `json:"tags"`
}

type format struct {
	Format_name string `json:"format_name"`
	Duration    string `json:"DURATION"`
	Size        string `json:"size"`
}
