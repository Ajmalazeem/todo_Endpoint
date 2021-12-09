package todosvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler(Progranmmer error)")
)

func MakeHTTPHandler(s Service)http.Handler{
	r:= mux.NewRouter()
	e :=MakeServerEndpoint(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	//post
	r.Methods("POST").Path("/todo/").Handler(httptransport.NewServer(
		e.PostTodoEndpoint,
		decodePostTodoRequest,
		encodeResponse,
		options...
	))

	//get
	r.Methods("GET").Path("/todo/{id}").Handler(httptransport.NewServer(
		e.GetTodoEndpoint,
		decodeGetTodoRequest,
		encodeResponse,
		options...
	))
	
	//put
	r.Methods("PUT").Path("/todo/{id}").Handler(httptransport.NewServer(
		e.PutTodoEndpoint,
		decodePutTodoRequest,
		encodeResponse,
		options...
	))

	//Delete
	r.Methods("DELETE").Path("/todo/{id}").Handler(httptransport.NewServer(
		e.DeleteTodoEndpoint,
		decodeDeleteTodoRequest,
		encodeResponse,
		options...
	))

	return r
}


//DECODE request

func decodePostTodoRequest(_ context.Context,r *http.Request)(request interface{},err error){
	var req PostTodoRequest
	if e :=json.NewDecoder(r.Body).Decode(&req.Todo); e!= nil{
		return nil, e
	}
	return req, nil

}

func decodeGetTodoRequest(_ context.Context, r *http.Request)(request interface{}, err error){
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return GetTodoRequest{ID: id},nil
}

func decodePutTodoRequest(_ context.Context, r *http.Request)(request interface{}, err error){
	vars := mux.Vars(r)
	id , ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var Todo Todo
	if err := json.NewDecoder(r.Body).Decode(&Todo); err !=nil{
		return nil, err
	}
	return PutTodoRequest{
		ID: id,
		Todo: Todo,

	},nil
}

func decodeDeleteTodoRequest(_ context.Context,r *http.Request)(request interface{},err error){
	vars := mux.Vars(r)
	id, ok :=vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return DeleteTodoRequest{ID: id}, nil
}

//DECODE Response

func decodePostTodoResponse(_ context.Context,resp *http.Response)(interface{},error){
	var response PostTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetTodoResponse(_ context.Context,resp *http.Response)(interface{},error){
	var response GetTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutTodoResponse(_ context.Context, resp *http.Response)(interface{},error){
	var response PutTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteTodoResponse(_ context.Context, resp *http.Response)(interface{},error){
	var response DeleteTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response , err
}

//ENCODE request

func encodeRequest(ctx context.Context,req *http.Request,request interface{})error{
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil{
		return err
	}
	req.Body=ioutil.NopCloser(&buf)
	return nil
}

func encodePostTodoRequest(ctx context.Context,req *http.Request,request interface{})error{
	req.URL.Path = "/todo/"
	return encodeRequest(ctx, req, request) 
}

func encodeGetTodoRequest(ctx context.Context,req *http.Request, request interface{})error{
	r := request.(GetTodoRequest)
	todoID := url.QueryEscape(r.ID)
	req.URL.Path = "/todo/"+todoID
	return encodeRequest(ctx, req, request)
}

func encodePutTodoRequest(ctx context.Context, req *http.Request,request interface{})error{
	r :=request.(PutTodoRequest)
	todoID := url.QueryEscape(r.ID)
	req.URL.Path = "/todo/" + todoID
	return encodeRequest(ctx,req, request)
}

func encodeDeleteTodoRequest(ctx context.Context,req *http.Request,request interface{})error{
	r := request.(DeleteTodoRequest)
	todoID := url.QueryEscape(r.ID)
	req.URL.Path = "/todo/"+ todoID
	return encodeRequest(ctx, req, request)
}

//ENCODE Response
type errorer interface{
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{})error{
	if e, ok := response.(errorer); ok && e.error != nil {
			encodeError(ctx, e.error(), w)
			return nil
	}
	w.Header().Set("Content-Type","application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter){
	if err == nil{
		panic("encode error with nil error")
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int{
	switch  err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExits,ErrInconsistentId:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError	
	}
}


