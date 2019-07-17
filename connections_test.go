package connections

import (
	"testing"
)

func Test_connections_connect(t *testing.T) {

	tests := []struct {
		name   string
		region string
		match  bool
	}{
		{
			name:   "connect to region foo",
			region: "foo",
			match:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connections{}
			c.Connect(tt.region)
			if (c.Region == tt.region) != tt.match {
				t.Errorf("Connections.connect() c.region = %v, expected %v",
					c.Region, tt.region)
			}
		})
	}
}
