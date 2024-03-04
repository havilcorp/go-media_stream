package miniostore

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FilmModel struct {
	Name string
}

type Minio struct {
	MinioClient *minio.Client
}

func NewMinioStore() *Minio {
	endpoint := "192.168.0.118:9000"
	accessKeyID := "yorAXzMbHw6Uhy5mwQBV"
	secretAccessKey := "5NwpgqbIPVUea2t2moGDzRhGo6YZREM9yYv1fbHR"
	useSSL := false
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// listCh := minioClient.ListObjects(ctx, "films", minio.ListObjectsOptions{})
	// for object := range listCh {
	// 	if object.Err != nil {
	// 		fmt.Println(object.Err)
	// 		return nil
	// 	}
	// 	fmt.Println(object.ETag, object.Key)
	// }

	// Генерация ссылки на чтение файла
	// objectName := "films/кунг-фу панда 3 low.mp4"
	// expiry := 24 * time.Hour // Время жизни ссылки
	// presignedURL, err := minioClient.PresignedGetObject(context.Background(), "films", objectName, expiry, nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("Ссылка на чтение файла:", presignedURL)
	return &Minio{
		MinioClient: minioClient,
	}
}

func (s *Minio) GetFilms(ctx context.Context, prefix string) ([]string, []string, error) {
	options := minio.ListObjectsOptions{}
	if prefix != "" {
		options.Recursive = true
		if prefix[len(prefix)-1:] != "/" {
			options.Prefix = fmt.Sprintf("%s/", prefix)
		} else {
			options.Prefix = prefix
		}
		fmt.Println(options.Prefix)
	}
	listCh := s.MinioClient.ListObjects(ctx, "films", options)
	folders := make([]string, 0)
	files := make([]string, 0)
	for object := range listCh {
		if object.Err != nil {
			return nil, nil, object.Err
		}
		if object.ETag == "" {
			folders = append(folders, object.Key)
		} else {
			key := object.Key[len(options.Prefix):]
			split := strings.Split(key, "/")
			if len(split) != 1 {
				find := false
				for _, a := range folders {
					if a == split[0] {
						find = true
						break
					}
				}
				if !find {
					folders = append(folders, split[0])
				}
			} else {
				files = append(files, object.Key)
			}
		}
	}
	return folders, files, nil
}

func (s *Minio) GetObjects() {
	object, err := s.MinioClient.GetObject(context.Background(), "films", "кунг-фу панда 3 low.mp4", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer object.Close()
	o, err := object.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(o)
}
