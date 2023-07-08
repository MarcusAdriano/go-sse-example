package httpext

import (
	"net/http"

	"github.com/google/uuid"
)

type SseEmitter interface {
	Id() string
	Emit(data []byte) error
}

type sseEmitter struct {
	id string
	w  http.ResponseWriter
}

func newSseEmitter(w http.ResponseWriter) SseEmitter {
	return &sseEmitter{
		id: uuid.New().String(),
		w:  w,
	}
}

func (s *sseEmitter) Id() string {
	return s.id
}

func (s *sseEmitter) Emit(data []byte) error {
	_, err := s.w.Write(data)
	if err != nil {
		return err
	}

	_, err = s.w.Write([]byte("\n\n"))
	if err != nil {
		return err
	}

	flusher, ok := s.w.(http.Flusher)
	if !ok {
		return err
	}

	flusher.Flush()

	return nil
}

type SseSubscriber interface {
	Subscribe() error
}
