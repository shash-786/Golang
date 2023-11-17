package hello

import "testing"

func Test(t *testing.T) {
	subtests := []struct {
		items  []string
		result string
	}{
		{
			result: "Hello World!",
		},

		{
			items:  []string{"matt", "sarah"},
			result: "Hello matt, sarah!",
		},
	}

	for _, value := range subtests {
		if s := Say(value.items); s != value.result {
			t.Errorf("want:%s got:%s", value.result, s)
		}
	}
}
