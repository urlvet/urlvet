package domaininfo

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestSetAge(t *testing.T) {
	// .de and .eu publish no creation date.
	unknown := &RegistrationData{}
	if err := setAge(unknown); err != nil {
		t.Fatal(err)
	}
	if _, ok := unknown.Age(); ok || unknown.AgeKnown {
		t.Errorf("missing creation date read as a known age")
	}
	b, _ := json.Marshal(unknown)
	for _, want := range []string{`"age_known":false`, `"age_days":null`, `"age_human":""`} {
		if !strings.Contains(string(b), want) {
			t.Errorf("JSON %s missing %s", b, want)
		}
	}

	known := &RegistrationData{CreatedDate: time.Now().AddDate(0, 0, -10)}
	if err := setAge(known); err != nil {
		t.Fatal(err)
	}
	if days, ok := known.Age(); !ok || !known.AgeKnown || days != 10 {
		t.Errorf("Age() = %d, %v (known=%v), want 10, true", days, ok, known.AgeKnown)
	}
}
