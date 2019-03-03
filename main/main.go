package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"log"
	"net/http"
	"strings"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	svc := receiptService{}
	backoutReceiptsHandler := httptransport.NewServer(
		makeBackoutReceiptsEndpoint(svc),
		decodeBackoutReceiptsRequest,
		encodeResponse)
	getReceivedQtySummaryByPOForDeliveryHandler := httptransport.NewServer(
		makeGetReceivedQtySummaryByPOForDeliveryEndpoint(svc),
		decodeGetReceivedQtySummaryByPOForDeliveryRequest,
		encodeResponse)
	heartBeatHandler := httptransport.NewServer(
		makeHeartBeatEndpoint(svc),
		decodeHeartbeatRequest,
		encodeResponse)
	http.Handle("/receipts/backout", backoutReceiptsHandler)
	http.Handle("/receipts/delivery/:deliveryNumber/summary", getReceivedQtySummaryByPOForDeliveryHandler)
	http.Handle("/heartbeat", heartBeatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decodeBackoutReceiptsRequest(_ context.Context, req *http.Request) (interface{}, error){
	var request backoutReceiptsRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func decodeHeartbeatRequest(_ context.Context, req *http.Request) (interface{}, error){
	var request heartBeatRequest
	if req.Method != "GET" {
		return request, nil
	}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetReceivedQtySummaryByPOForDeliveryRequest(_ context.Context, req *http.Request) (interface{}, error){
	var request getReceivedQtySummaryByPOForDeliveryRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, writer http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(writer).Encode(response)
}


type ReceiptService interface {
	GetReceivedQtySummaryByPOForDelivery(string) (string, error)
	BackoutReceipts(string) (string, error)
	HeartBeat() string

}



type receiptService struct {}

func (receiptService) GetReceivedQtySummaryByPOForDelivery(body string) (string, error){
	return "", nil
}

func (receiptService) BackoutReceipts(body string) (string, error){
	return strings.ToUpper(body), nil
}

func (receiptService) HeartBeat() string {
	return "Receipt service is alive :) "
}

type getReceivedQtySummaryByPOForDeliveryRequest struct {
	S string `json:"s"`
}
type getReceivedQtySummaryByPOForDeliveryResponse struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}
type backoutReceiptsRequest struct {
	DeliveryNumber string `json:"deliveryNumber"`
}
type backoutReceiptsResponse struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}
type heartBeatRequest struct {
}
type heartBeatResponse struct {
	V string `json:"v"`
}

func makeGetReceivedQtySummaryByPOForDeliveryEndpoint(svc ReceiptService) endpoint.Endpoint{
	return func(_ context.Context, request interface{}) (interface{},error){
			req := request.(getReceivedQtySummaryByPOForDeliveryRequest)
			v, err := svc.GetReceivedQtySummaryByPOForDelivery(req.S)
			if err != nil {
				return getReceivedQtySummaryByPOForDeliveryResponse{v, err.Error()}, nil
			}
			return getReceivedQtySummaryByPOForDeliveryResponse{v, ""}, nil
	}
}

func makeBackoutReceiptsEndpoint(svc ReceiptService) endpoint.Endpoint{
	return func(_ context.Context, request interface{}) (interface{},error){
		req := request.(backoutReceiptsRequest)
		v, err := svc.BackoutReceipts(req.DeliveryNumber)
		if err != nil {
			return backoutReceiptsResponse{v, err.Error()}, nil
		}
		return backoutReceiptsResponse{v, ""}, nil
	}
}

func makeHeartBeatEndpoint(svc ReceiptService) endpoint.Endpoint{
	return func(_ context.Context, request interface{}) (interface{},error){
		v := svc.HeartBeat()
		return heartBeatResponse{v}, nil
	}
}