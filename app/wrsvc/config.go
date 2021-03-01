package wrsvc

import (
	"fmt"
	"os"
)

type WriteAPIConfg struct {
	token  string
	bucket string
	org    string
	url    string
}

func GetWriteAPIConfig() WriteAPIConfg {
	mustString := func(s string) string {
		if v, ok := os.LookupEnv(s); ok {
			return v
		}

		panic(fmt.Sprintf("environment value for %q not found", s))
	}

	return WriteAPIConfg{
		token:  mustString("WRCLI_TOKEN"),
		bucket: mustString("WRCLI_BUCKET"),
		org:    mustString("WRCLI_ORG"),
		url:    mustString("WRCLI_URL"),
	}
}
