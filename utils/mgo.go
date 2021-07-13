package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
)

func GetMongoPartialUpdateMap(aMap map[string]interface{}) bson.M {
	finalMap := make(map[string]interface{})
	parseMap("", aMap, &finalMap)
	return finalMap
}

func parseMap(k string, aMap map[string]interface{}, finalMap *map[string]interface{}) {
	if len(aMap) == 0 {
		(*finalMap)[k] = nil
		return
	}

	for key, val := range aMap {
		if val != nil {
			switch concreteVal := val.(type) {
			case map[string]interface{}:
				parseMap(getKey(k, key), val.(map[string]interface{}), finalMap)
			case []interface{}:
				(*finalMap)[getKey(k, key)] = val.([]interface{})
			default:
				concreteValType := reflect.TypeOf(concreteVal)
				if concreteValType.Kind() == reflect.Map {
					parseMap(getKey(k, key), concreteVal.(primitive.M), finalMap)
				} else {
					(*finalMap)[getKey(k, key)] = concreteVal
				}
			}
		} else {
			(*finalMap)[getKey(k, key)] = nil
		}
	}
}

func getKey(k string, key string) string {
	if k == "" {
		return key
	}
	return k + "." + key
}

func ToBsonMDoc(v interface{}) (*bson.M, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	//log.Println("data:", string(data))

	var doc bson.M
	err = bson.Unmarshal(data, &doc)
	return &doc, err
}

func GetStringFromObjectID(id primitive.ObjectID) string {
	return strings.ReplaceAll(fmt.Sprintf("%q", id.Hex()), `"`, "")
}

func AppendSearchPattern(filter *bson.D, key, value string, caseInsensitive bool) *bson.D {
	if value == "" {
		return filter
	}
	pattern := primitive.Regex{
		Pattern: "^" + value,
	}

	if caseInsensitive {
		pattern.Options = "i"
	}

	*filter = append(*filter, bson.E{Key: key, Value: pattern})

	return filter
}

func AppendStringValue(filter *bson.D, key, value string) *bson.D {
	if value != "" {
		*filter = append(*filter, bson.E{Key: key, Value: value})
	}
	return filter
}

func AppendTimeValue(filter *bson.D, key string, value time.Time, isGreater bool) *bson.D {
	if !value.IsZero() {
		dateCompare := "$gte"
		if !isGreater {
			dateCompare = "$lte"
		}

		*filter = append(*filter, bson.E{Key: key, Value: bson.M{dateCompare: value}})
	}
	return filter
}

func GetSearchPattern(key, value string, caseInsensitive bool) interface{} {
	pattern := primitive.Regex{
		Pattern: value,
	}

	if caseInsensitive {
		pattern.Options = "i"
	}

	filter := bson.E{Key: key, Value: pattern}

	return filter
}
