package handler

import (
	cfg "biscotti/config"
	b64 "encoding/base64"
	"encoding/json"

	"fmt"
	"github.com/dangkaka/go-kafka-avro"
	"github.com/linkedin/goavro/v2"
	"github.com/valyala/fasthttp"
	"math"
	"time"
)

type Message struct {
	Value     string  `json:"value"`
	Timestamp float64 `json:"timestamp"`
}

type Schema struct {
	id int
}

var schema *Schema

var codec *goavro.Codec

func GetCodec() *goavro.Codec {
	if codec == nil {
		codec, _ = goavro.NewCodec(
			fmt.Sprintf(`
        {
          "type": "record",
          "name": "%s",
          "fields" : [
            {"name": "user_id", "type": "string"},
            {"name": "network_user_id", "type": "string"},
            {"name": "network_id", "type": "string"}
          ]
        }`, cfg.GetAppConfig().KafkaTopic),
		)
	}
	return codec

}

func GetSchemaID() *Schema {
	if schema == nil {
		schemaRegistry := kafka.NewCachedSchemaRegistryClientWithRetries([]string{cfg.GetAppConfig().SchemaRegistryURL}, 5)
		schemaID, schemaRegistryErr := schemaRegistry.CreateSubject(cfg.GetAppConfig().KafkaTopic, GetCodec())
		if schemaRegistryErr != nil {
			fmt.Println(schemaRegistryErr.Error())
		}
		schema = &Schema{
			schemaID,
		}
	}
	return schema
}

func GetConfluentAvroBinary(data map[string]interface{}) ([]byte, error) {
	binary, err := GetCodec().BinaryFromNative(nil, data)
	if err != nil {
		fmt.Println("Encoding Error:", err.Error())
	}

	avroEncoder := kafka.AvroEncoder{
		SchemaID: GetSchemaID().id,
		Content:  binary,
	}
	return avroEncoder.Encode()
}

func addMatchedCookies(networkID, networkUserID, userID string) {
	data := map[string]interface{}{
		"user_id":         userID,
		"network_user_id": networkUserID,
		"network_id":      networkID}
	encodedData, _ := GetConfluentAvroBinary(data)

	loggerStorage := GetLoggerStorage()
	sEnc := b64.StdEncoding.EncodeToString(encodedData)

	message := Message{Value: sEnc, Timestamp: float64(time.Now().UnixNano()) / math.Pow(10, 9)}
	marshaledData, _ := json.Marshal(message)
	err := loggerStorage.Store(marshaledData)
	if err != nil {
		return
	}
}

func MatchHandler(ctx *fasthttp.RequestCtx) {
	networkID := string(ctx.QueryArgs().Peek("id"))
	networkUserID := string(ctx.QueryArgs().Peek("user_id"))
	userID := string(ctx.Request.Header.Cookie(cfg.GetAppConfig().CookieKey))
	go addMatchedCookies(networkID, networkUserID, userID)
}
