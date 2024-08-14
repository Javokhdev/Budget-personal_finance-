package kafka

import (
	"encoding/json"
	"log"
	"budget-service/model"
)

func CreateNotification(kaf KafkaProducer, request *model.Send) error {
	response, err := json.Marshal(request)
	if err != nil {
		log.Println("cannot produce messages via kafka", err.Error())
		return err
	}
	err = kaf.ProduceMessages("create", response)
	if err != nil {
		log.Fatal("Error while ProduceMessages: ", err.Error())
		return err
	}
	return nil
}
