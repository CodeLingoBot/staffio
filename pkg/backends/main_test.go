package backends

import (
	"log"
	"testing"

	. "lcgc/platform/staffio/pkg/settings"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	Settings.Parse()
	m.Run()
}