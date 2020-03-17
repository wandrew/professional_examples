//	Go Package that implements streadway protocols.
package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"

	"github.com/deseretdigital/sl-colony-domain-catalog/export/types"
	"github.com/deseretdigital/sl-colony-domain-catalog/internal/endpoint"

	// For DEV set the pkg dir
	amqpdriver "github.com/deseretdigital/sl-colony-platform/drivers/amqp" // PROD
	// amqpdriver "github.com/deseretdigital/sl-colony-domain-catalog/pkg/sl-colony-platform/drivers/amqp" // DEV

	"github.com/go-kit/kit/log"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

type Subscribers map[string]*SubscriberQueue

type SubscriberProfile struct {
	Logger   log.Logger
	Endpoint endpoint.Endpoints
}

type SubscriberQueue struct {
	Subscriber *amqptransport.Subscriber
	Queue      *amqpdriver.Queue
}

// I am going to hate myself for this code in like... a month
func (p SubscriberProfile) Init(ex *amqpdriver.Exchange) *Subscribers {
	subscribers := make(Subscribers)

	options := []amqptransport.SubscriberOption{
		amqptransport.SubscriberErrorLogger(p.Logger),
	}

	subscribers["catalogItem"] = &SubscriberQueue{
		Subscriber: amqptransport.NewSubscriber(
			p.Endpoint.CatalogItemEndpoint,
			DecodeAMQPCatalogItemRequest,
			EncodeAMQPCatalogItemResponse,
			options...,
		),
		Queue: &amqpdriver.Queue{
			Name:       "domain-catalog-item",
			BindingKey: "catalog-item",
			Exchange:   *ex,
		}}

	subscribers["allItems"] = &SubscriberQueue{
		Subscriber: amqptransport.NewSubscriber(
			p.Endpoint.AllItemsEndpoint,
			DecodeAMQPCatalogItemRequest,
			EncodeAMQPAllItemsResponse,
			options...,
		),
		Queue: &amqpdriver.Queue{
			Name:       "domain-catalog-allitems",
			BindingKey: "catalog-allitems",
			Exchange:   *ex,
		}}

	subscribers["listingCategoryItems"] = &SubscriberQueue{
		Subscriber: amqptransport.NewSubscriber(
			p.Endpoint.ListingCategoriesEndpoint,
			DecodeAMQPCatalogItemRequest,
			EncodeAMQPListinCategoryItemsResponse,
			options...,
		),
		Queue: &amqpdriver.Queue{
			Name:       "domain-catalog-listingcategoryitems",
			BindingKey: "catalog-listingcategoryitems",
			Exchange:   *ex,
		}}

	subscribers["geoItem"] = &SubscriberQueue{
		Subscriber: amqptransport.NewSubscriber(
			p.Endpoint.GeoItemEndpoint,
			DecodeAMQPGeoItemRequest,
			EncodeAMQPGeoItemResponse,
			options...,
		),
		Queue: &amqpdriver.Queue{
			Name:       "domain-catalog-geoitem",
			BindingKey: "catalog-geoitem",
			Exchange:   *ex,
		}}

	return &subscribers
}

// Maybe the cleanest thing I've done in this whole thing. -- Not after adding the ch arg though....
func (s *SubscriberQueue) Deliver(ch *rabbitmq.Channel) {
	for d := range s.Queue.Deliver(ch) {
		s.Subscriber.ServeDelivery(ch)(&d)
	}
}

func DecodeAMQPCatalogItemRequest(_ context.Context, d *amqp.Delivery) (interface{}, error) {
	var obj *types.FilterParams
	err := json.Unmarshal(d.Body, &obj)
	if err != nil {
		fmt.Println(err)
	}

	return obj, err
}

func EncodeAMQPCatalogItemResponse(_ context.Context, p *amqp.Publishing, request interface{}) error {
	req := request.(*types.CatalogItem)
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.Body = b
	return nil
}

func EncodeAMQPAllItemsResponse(_ context.Context, p *amqp.Publishing, request interface{}) error {
	req := request.(*types.CatalogItems)
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.Body = b
	return nil
}

func EncodeAMQPListinCategoryItemsResponse(_ context.Context, p *amqp.Publishing, request interface{}) error {
	req := request.(*types.ListingCategoryItems)
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.Body = b
	return nil
}

func DecodeAMQPGeoItemRequest(_ context.Context, d *amqp.Delivery) (interface{}, error) {
	var obj *types.GeoItemParams
	err := json.Unmarshal(d.Body, &obj)
	if err != nil {
		fmt.Println(err)
	}

	return obj, err
}

func EncodeAMQPGeoItemResponse(_ context.Context, p *amqp.Publishing, request interface{}) error {
	req := request.(*types.GeoItem)
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.Body = b
	return nil
}
