package supabase

import (
	"mime/multipart"
	"os"

	supabase "github.com/adityarizkyramadhan/supabase-storage-uploader"
)

type SupabaseItf interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(link string) error
}

type storageSupabase struct {
	client *supabase.Client
}

func Init() SupabaseItf {
	spbClient := supabase.New(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_TOKEN"),
		os.Getenv("SUPABASE_BUCKET"),
	)

	return &storageSupabase{
		client: spbClient,
	}
}

func (s *storageSupabase) Upload(file *multipart.FileHeader) (string, error) {
	url, err := s.client.Upload(file)
	if err != nil {
		return url, err
	}

	return url, nil
}

func (s *storageSupabase) Delete(link string) error {
	err := s.client.Delete(link)
	if err != nil {
		return err
	}

	return nil
}
