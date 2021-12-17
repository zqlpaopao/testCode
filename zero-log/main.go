package main

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main(){
	type field struct {
		Name string
		Age int
	}
	var  field2 []field
	field2 = append(field2,field{
		Name: "name",
		Age:  18,
	})
	var b []byte
	//var err error
	b,_ = json.Marshal(&field2)
	b =b
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("hello world")
	log.Debug().
		Str("name", "xy").
		Str("name", "xy").
		Float64("height", 100.0).Dict("dict", zerolog.Dict().
		Str("bar", "baz").
		Int("n", 1),
	).
		Str("bar", "baz").
		Int("n", 1).Msg("msg")

	log.Info().
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Msg("hello world")
}