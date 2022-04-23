package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()
	query, err := gojq.Parse(".foo , .foo.[1]")
	if err != nil {
		exitCode = 1
		fmt.Print(err.Error())
		return
	}
	input := map[string]interface{}{"foo": []interface{}{1, 2, 3}}
	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			exitCode = 1
			fmt.Print(err.Error())
			return
		}
		fmt.Printf("%#v\n", v)
		var js []byte
		if js, err = json.Marshal(v); err != nil {
			fmt.Print(err.Error())
			return
		}
		fmt.Println(string(js))
	}
}
