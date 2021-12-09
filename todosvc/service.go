package todosvc


import (
	"context"
	"errors"
	"sync"
)

type Service interface{
	PostTodo(ctx context.Context,t Todo)error
	GetTodo(ctx context.Context,id string)(Todo, error)
	PutTodo(ctx context.Context,id string,t Todo)error
	DeleteTodo(ctx context.Context,id string)error
}

type Todo struct{
	ID 			string 	`json:"id"`
 	Todo 		string 	`json:"todo"`
 	Completed 	bool 	`json:"completed"`
}

var(
	ErrInconsistentId = errors.New("inconsistent id")
	ErrAlreadyExits = errors.New("already exits")
	ErrNotFound = errors.New("not found")
)

type inmemService struct{
	mtx sync.RWMutex
	m map[string]Todo
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Todo{},
	}
}

func(s *inmemService) PostTodo(ctx context.Context,t Todo)error{
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[t.ID]; ok {
		return ErrAlreadyExits //cant ovrwrite
	}
	s.m[t.ID] = t
	return nil
}

func(s *inmemService) GetTodo(ctx context.Context,id string)(Todo, error){
	s.mtx.Lock()
	defer s.mtx.Unlock()
	t, ok := s.m[id]
	if !ok{
		return Todo{},ErrNotFound//cant find created id
	}
	return t, nil
}

func(s *inmemService) PutTodo(ctx context.Context,id string,t Todo)error{
	if id != t.ID{
		return ErrInconsistentId//cant create or update
	}
	s.m[id] = t
	return nil
}
func(s *inmemService) DeleteTodo(ctx context.Context,id string)error{
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok{
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}
