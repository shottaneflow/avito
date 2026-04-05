package getselleritems

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("ALLURE_OUTPUT_PATH", "../../..")
	os.Exit(m.Run())
}
