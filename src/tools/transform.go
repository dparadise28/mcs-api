//package main

package tools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	//"github.com/mailru/easyjson"
	"github.com/julienschmidt/httprouter"
	"github.com/mailru/easyjson/jlexer"
	"github.com/pquerna/ffjson/ffjson"
	//"github.com/mailru/easyjson/jwriter"
)

var ArrayId int

type TransformationTracer struct {
	/*
		this struct is meant to keep track of
	*/

	//trackers
	objects    map[string]interface{}   // map of {path: obj} of objects processed to be re-assembled
	newObjs    []map[string]interface{} // array that keeps track of {path: obj} to objects not yet processed
	pathToObjs []string                 // array of paths to objects in original dismantled obj

	arrays     map[string]interface{} // array that keeps track of {path: ary} to objects not yet processed
	newArrays  map[string]interface{} // array that keeps track of arrays not yet processed
	pathToArys []string

	//objects
	preprocessed *[]byte                // original data to be remodled
	processed    map[string]interface{} // result of remodling

	structureType string // ex, json (currently only support json)
}

func dismantleObj(remodler *TransformationTracer) {
	pathToObj := ""
	for len(remodler.newObjs) > 0 {
		for key, value := range remodler.newObjs[0] {
			if key == "-----_PATH_TO_OBJECT_-----" {
				continue
			}
			if val, ok := remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"]; ok {
				pathToObj = val.(string) + "." + key
			} else {
				pathToObj = key
			}
			switch value.(interface{}).(type) {
			case string:
				for _, keyPath := range strings.Split((value).(string), " |or| ") {
					if val, _, _, err := jsonparser.Get(*remodler.preprocessed, strings.Split(string(keyPath), ".")...); err == nil {
						remodler.newObjs[0][key] = string(val)
						break
					}
				}
			case map[string]interface{}:
				remodler.newObjs = append(remodler.newObjs, value.(map[string]interface{}))
				remodler.newObjs[len(remodler.newObjs)-1]["-----_PATH_TO_OBJECT_-----"] = pathToObj //val.(string) + "." + key
				delete(remodler.newObjs[0], key)
			case []interface{}:
				dismantleArray(remodler, &pathToObj, &value)
			default: // for any values that cant be a path or nested element keep whats set as default and move on
				continue
			}
		}
		pathToObj = remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"].(string)
		if pathToObj != "~~~~~-----%!%!%Root%!%!%-----~~~~~" {
			remodler.pathToObjs = append(remodler.pathToObjs, pathToObj)
		}
		remodler.objects[pathToObj] = remodler.newObjs[0]
		delete(remodler.objects[pathToObj].(map[string]interface{}), "-----_PATH_TO_OBJECT_-----")
		remodler.newObjs = remodler.newObjs[1:]
	}
}

func dismantleArray(remodler *TransformationTracer, pathToObj *string, value *interface{}) {
	for _, item := range (*value).([]interface{}) {
		currentKey := *pathToObj + "." + strconv.Itoa(ArrayId)
		remodler.newArrays[currentKey] = []interface{}{item} //temp
		switch item.(interface{}).(type) {
		case string:
			path := strings.Split(item.(string), ".")
			jsonparser.ArrayEach(*remodler.preprocessed, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				newVal, _, _, err := jsonparser.Get(value, path[len(path)-1])
				if err == nil {
					remodler.newArrays[currentKey] = append(remodler.newArrays[currentKey].([]interface{}), string(newVal))
				}
			}, path[:len(path)-1]...)
			remodler.arrays[currentKey] = remodler.newArrays[currentKey].([]interface{})[1:]
			delete(remodler.newArrays, currentKey)
			//case map[string]interface{}:

		}
		remodler.pathToArys = append(remodler.pathToArys, currentKey)
		ArrayId += 1
	}
}

func reassembleObj(remodler *TransformationTracer) {
	for len(remodler.pathToObjs) > 0 {
		index := len(remodler.pathToObjs) - 1
		path := strings.Split(remodler.pathToObjs[index], ".")
		parent := strings.Join(path[:len(path)-1], ".")

		remodler.objects[parent].(map[string]interface{})[path[len(path)-1]] = remodler.objects[remodler.pathToObjs[index]]
		delete(remodler.objects, remodler.pathToObjs[index])
		remodler.pathToObjs = remodler.pathToObjs[:index]
	}
	remodler.objects = remodler.objects["~~~~~-----%!%!%Root%!%!%-----~~~~~"].(map[string]interface{})
}

func Remodel(expected []byte, original []byte, action []byte) TransformationTracer {
	//map string interface representation of json input
	expectedJsonItr := jlexer.Lexer{Data: expected}
	expectedJsonMSI := expectedJsonItr.Interface().(map[string]interface{})
	remodler := TransformationTracer{
		make(map[string]interface{}, 0),   // obj
		make([]map[string]interface{}, 0), // obj
		make([]string, 0),                 // pathToObjs
		make(map[string]interface{}, 0),   // arrays
		make(map[string]interface{}, 0),   // arrays
		make([]string, 0),
		&original,                    // preprocessed
		make(map[string]interface{}), // processed
		"json", // string
	}

	remodler.newObjs = append(remodler.newObjs, expectedJsonMSI)
	remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"] = "~~~~~-----%!%!%Root%!%!%-----~~~~~"
	dismantleObj(&remodler)
	if bytes.Equal(action, []byte(`dismantle`)) {
		return remodler
	}
	reassembleObj(&remodler)
	return remodler
}

func printJ(JsonMSI map[string]interface{}) {
	buf, _ := ffjson.Marshal(&JsonMSI)
	log.Println(string(buf))
}

func printI(JsonMSI interface{}) {
	buf, _ := ffjson.Marshal(&JsonMSI)
	log.Println(string(buf))
}

func RemodelJ(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("could not read request body")
	} else {
		inputs := map[string][]byte{
			"action":             []byte(``),
			"response_structure": []byte(``),
			"original_structure": []byte(``),
		}

		for key, _ := range inputs {
			if val, _, _, err := jsonparser.Get(body, strings.Split(key, ".")...); err == nil {
				inputs[key] = val
			}
		}

		remodler := Remodel(inputs["response_structure"], inputs["original_structure"], inputs["action"])
		buf, _ := ffjson.Marshal(&remodler.objects) //&expectedJsonMSI)
		log.Println(string(buf))
		fmt.Fprintf(w, string(buf))
	}
}
