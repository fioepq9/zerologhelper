package main

import (
	"github.com/cockroachdb/errors"
	"github.com/fioepq9/zerologhelper"
	"github.com/rs/zerolog"
)

func foo() error {
	return errors.New("foo")
}

func bar() error {
	return errors.Wrap(foo(), "bar")
}

func baz() error {
	return errors.Wrap(bar(), "baz")
}

func main() {
	h := zerologhelper.New()

	h.SetInterfaceMarshalFunc().SetErrorStackMarshaler()

	log := zerolog.New(h.ConsoleWriter()).With().Stack().Caller().Timestamp().Logger()

	log.Info().Str("foo", "bar").Msg("hello world")

	log.Err(baz()).Msg("error")
}
