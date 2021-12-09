package todosvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct{
	PostTodoEndpoint endpoint.Endpoint
	GetTodoEndpoint endpoint.Endpoint
	PutTodoEndpoint endpoint.Endpoint
	DeleteTodoEndpoint endpoint.Endpoint
}

//server side
func MakeServerEndpoint(s Service) Endpoints{
	return Endpoints{
		PostTodoEndpoint: MakePostTodoEndpoint(s),
		GetTodoEndpoint: MakeGetTodoEndpoint(s),
		PutTodoEndpoint: MakePutTodoEndpoint(s),
		DeleteTodoEndpoint: MakeDeleteTodoEndpoint(s),
	}
}

//client side
func MakeClientEndpoint(instance string) (Endpoints,error){
	if !strings.HasPrefix(instance,"Http"){
		instance = "http://"+ instance
	}
	tgt , err :=url.Parse(instance)
	if err != nil{
		return Endpoints{},err
	}
	tgt.Path =""
	options := []httptransport.ClientOption{}

	return Endpoints{
		PostTodoEndpoint: httptransport.NewClient("POST",tgt,encodePostTodoRequest,decodePostTodoResponse,options...).Endpoint(),
		GetTodoEndpoint: httptransport.NewClient("GET",tgt,encodeGetTodoRequest,decodeGetTodoResponse,options...).Endpoint(),
		PutTodoEndpoint: httptransport.NewClient("PUT",tgt,encodePutTodoRequest,decodePutTodoResponse,options...).Endpoint(),
		DeleteTodoEndpoint: httptransport.NewClient("DELETE",tgt,encodeDeleteTodoRequest,decodeDeleteTodoResponse,options...).Endpoint(),
	},nil
}


func (e Endpoints) PostTodo(ctx context.Context,t Todo)error{
	request := PostTodoRequest{Todo: t}
	response , err := e.PostTodoEndpoint(ctx, request)
	if err != nil{
		return err
	}
	resp := response.(PostTodoResponse)
	return resp.Err
}

func (e Endpoints) GetTodo(ctx context.Context,id string)(Todo, error){
	request := GetTodoRequest{ID: id}
	response , err := e.GetTodoEndpoint(ctx, request)
	if err != nil{
		return Todo{}, err
	}
	resp :=response.(GetTodoResponse)
	return resp.Todo,resp.Err
}

func (e Endpoints) PutTodo(ctx context.Context,id string,t Todo)error{
	request := PutTodoRequest{ID : id}
	response, err := e.PutTodoEndpoint(ctx, request)
	if err != nil{
		return err
	}
	resp :=  response.(PutTodoResponse)
	return resp.Err

}

func (e Endpoints) DeleteTodo(ctx context.Context,id string)error{
	request := DeleteTodoRequest{ID: id}
	response, err := e.DeleteTodoEndpoint(ctx, request)
	if err != nil{
		return err
	}
	resp := response.(DeleteTodoResponse)
	return resp.Err
}


//endpoint
func MakePostTodoEndpoint(s Service)endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostTodoRequest)
		e := s.PostTodo(ctx, req.Todo)
		return PostTodoResponse{Err: e},nil
	}
}

func MakeGetTodoEndpoint(s Service)endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetTodoRequest)
		t,e := s.GetTodo(ctx, req.ID)
		return GetTodoResponse{Todo: t,Err: e},nil
	}
}

func MakePutTodoEndpoint(s Service)endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostTodoRequest)
		e := s.PutTodo(ctx, req.Todo.ID,req.Todo)
		return PutTodoResponse{Err:e}, nil
	}
}

func MakeDeleteTodoEndpoint(s Service)endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTodoRequest)
		e := s.DeleteTodo(ctx , req.ID)
		return DeleteTodoResponse{Err: e}, nil
	}
}

//requests & response Post
type PostTodoRequest struct{
	Todo Todo
}

type PostTodoResponse struct{
	Err error `json:"err,omitempty"`
}

func (r PostTodoResponse) error() error{
	return r.Err
}

//requests & response Get
type GetTodoRequest struct{
	ID string
}

type GetTodoResponse struct{
	Todo Todo `json:"todo,omitempty"`
	Err error `json:"err,omitempty"`
}

func (r GetTodoResponse) error() error{
	return r.Err
}

//requests & response Put
type PutTodoRequest struct{
	ID string
	Todo Todo
}

type PutTodoResponse struct{
	Err error `json:"err,omitempty"`
}

func (r PutTodoResponse) error() error{
	return nil
}

//requests & response Delete
type DeleteTodoRequest struct{
	ID string
}

type DeleteTodoResponse struct{
	Err error `json:"err,omitempty"`
}

func (r DeleteTodoResponse) error() error{
	return r.Err
}

