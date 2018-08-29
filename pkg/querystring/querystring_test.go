package querystring

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name      string `json:"name,omitempty"`
	Age       int    `json:"age,omitempty"`
	IsMarried bool   `json:"is_married,omitempty"`
}

func TestEncode(t *testing.T) {
	assert := assert.New(t)
	in := &testStruct{
		Name:      "john",
		Age:       10,
		IsMarried: true,
	}
	o := Encode(in)
	assert.Equal("age=10&is_married=true&name=john", o.Encode(), "should encode the struct into querystring")
}

func TestDecode(t *testing.T) {
	assert := assert.New(t)
	u := url.Values{}
	u.Add("name", "john")
	u.Add("age", "10")
	u.Add("is_married", "true")

	var o testStruct
	err := Decode(&o, u)
	assert.Nil(err)
	assert.Equal("john", o.Name, "should decode name")
	assert.Equal(10, o.Age, "should decode age")
	assert.Equal(true, o.IsMarried, "should decode is_married")
}
