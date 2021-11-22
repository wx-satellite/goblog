package main

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestInt64ToString(t *testing.T) {
	v := viper.New()
	v.Set("map", map[string]interface{}{
		"name": "age",
	})

	fmt.Println(v.Get("map.name"))
}
