package linkedin

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func TestHelloName(t *testing.T) {
	browser := rod.New().ControlURL(launcher.New().Leakless(false).MustLaunch()).Trace(true).MustConnect()
	defer browser.MustClose()
	in := New(browser)
	profile := in.RetrieveProfile()
	json.NewEncoder(log.Writer()).Encode(profile)
}
