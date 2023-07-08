package httpext

import (
	"fmt"
	"net/http"
)

var (
	errEmitterNotFound = fmt.Errorf("emitter not found")
)

type SseManager interface {
	NewEmitter(w http.ResponseWriter) SseEmitter
	EmitTo(id string, data []byte) error
	Emitter(id string) (*SseEmitter, bool)
	EmitToAll(data []byte) error
	Close(id string)
}

type sseManager struct {
	emitters map[string]SseEmitter
}

func NewSseManager() SseManager {
	return &sseManager{
		emitters: make(map[string]SseEmitter),
	}
}

func (s *sseManager) Emitter(id string) (*SseEmitter, bool) {
	emitter, ok := s.emitters[id]
	return &emitter, ok
}

func (s *sseManager) NewEmitter(w http.ResponseWriter) SseEmitter {
	emitter := newSseEmitter(w)
	s.emitters[emitter.Id()] = emitter
	return emitter
}

func (s *sseManager) EmitTo(id string, data []byte) error {
	emitter, ok := s.emitters[id]
	if !ok {
		return errEmitterNotFound
	}
	return emitter.Emit(data)
}

func (s *sseManager) EmitToAll(data []byte) error {
	for _, emitter := range s.emitters {
		err := emitter.Emit(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sseManager) Close(id string) {
	delete(s.emitters, id)
}
