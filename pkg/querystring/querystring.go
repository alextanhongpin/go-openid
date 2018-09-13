package querystring

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// TODO: NewDecoder(u).Decode(interface{})
func Decode(in interface{}, u url.Values) error {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			tag := f.Tag.Get("json")
			tags := strings.Split(tag, ",")
			if len(tags) == 0 {
				continue
			}
			tagName := tags[0]
			if strings.TrimSpace(tagName) == "" {
				continue
			}
			field := v.FieldByName(f.Name)
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			val := u.Get(tagName)
			switch f.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				tmp, _ := strconv.Atoi(val)
				field.SetInt(int64(tmp))
			case reflect.String:
				field.SetString(val)
			case reflect.Bool:
				field.SetBool(val == "true")
			}
		}
		return nil
	default:
		return fmt.Errorf("%v is not of type pointer", in)
	}
}

// TODO: NewEncoder(url.Values).Encode(interface{})
func Encode(in interface{}) url.Values {
	u := url.Values{}
	t := reflect.TypeOf(in)
	vt := reflect.ValueOf(in)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		vt = vt.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			v := vt.Field(i)
			tag := f.Tag.Get("json")
			tags := strings.Split(tag, ",")
			if len(tags) == 0 {
				continue
			}
			name := tags[0]
			if strings.TrimSpace(name) == "" {
				continue
			}
			z := reflect.Zero(v.Type())
			isZero := z.Interface() == v.Interface()
			omitempty := strings.HasSuffix(tag, ",omitempty")

			if omitempty && isZero {
				continue
			}

			var val interface{}
			switch f.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val = v.Int()
			case reflect.String:
				val = v.String()
			case reflect.Bool:
				val = v.Bool()
			}
			u.Add(name, fmt.Sprint(val))
		}
		return u
	default:
		return u
	}

}
