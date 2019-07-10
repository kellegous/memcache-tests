package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"runtests/api"
	"runtests/memcache"
)

func dump(w io.Writer, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if _, err = w.Write(b); err != nil {
		return err
	}

	_, err = w.Write([]byte{'\n'})
	return err
}

func keysNotHavingValueOf(gr *api.GetResult, exp string) []string {
	var keys []string
	for key, val := range gr.Result {
		if val != exp {
			keys = append(keys, key)
		}
	}
	return keys
}

func runAll(app *api.Client, mcs []*memcache.Client) {
	keys := []string{
		"00", "01", "02", "03", "04",
		"05", "06", "07", "08", "09",
	}

	// First set a value
	sr, err := app.Set(keys, "1")
	if err != nil {
		log.Panic(err)
	}

	dump(os.Stdout, sr)

	if !sr.Result {
		log.Panic("unable to set first value")
	}

	if err := mcs[0].SetState(false); err != nil {
		log.Panic(err)
	}

	// issue sets until server is marked dead, this
	// will result in the server being de-pooled and
	// the keys will be re-distributed.
	for {
		time.Sleep(2 * time.Second)
		sr, err := app.Set(keys, "1")
		if err != nil {
			log.Panic(err)
		}

		dump(os.Stdout, sr)
		if sr.ResultCode == 35 {
			break
		}
	}

	// at this point, we should get a value for all our
	// keys.
	gr, err := app.Get(keys)
	if err != nil {
		log.Panic(err)
	}
	dump(os.Stdout, gr)

	// now we're going to set the keys to a new value
	// with the server still failed.
	sr, err = app.Set(keys, "2")
	if err != nil {
		log.Panic(err)
	}
	dump(os.Stdout, sr)

	// heal the network partition
	if err := mcs[0].SetState(true); err != nil {
		log.Panic(err)
	}

	for {
		time.Sleep(1 * time.Second)
		gr, err := app.Get(keys)
		if err != nil {
			log.Panic(err)
		}

		sk := keysNotHavingValueOf(gr, "2")
		if len(sk) != 0 {
			fmt.Printf("stale values return for keys %s\n",
				strings.Join(sk, ", "))
			dump(os.Stdout, gr)
			break
		}
	}
}
