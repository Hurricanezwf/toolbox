package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CloudSecurityGroupCreateReq struct {
	DriverTres *string `json:"driverTres,omitempty" driver:"driverTres"`
	//DriverCuatro *string `json:"driverCuatro,omitempty"`
	//SGName       *string `json:"sgName,omitempty"`
	//SGDesc       *string `json:"sgDesc,omitempty"`
	VPCID     []int  `json:"vpcId,omitempty" driver:"vpcId"`
	ProjectID *int64 `json:"projectId,omitempty"`
}

func main() {
	r := CloudSecurityGroupCreateReq{
		DriverTres: StringPtr("xxx"),
		//VPCID:      StringPtr("yyy"),
		VPCID:     []int{1, 5, 2},
		ProjectID: Int64Ptr(22),
	}

	p := NewFormParser("driver")
	m, err := p.Parse(&r)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v\n", m)
}

type FormParser struct {
	tag string
}

func NewFormParser(tag string) *FormParser {
	tag = strings.Trim(tag, " ")
	if len(tag) <= 0 {
		panic("tag could be empty")
	}
	return &FormParser{tag}
}

func (p FormParser) Parse(v interface{}) (map[string]string, error) {
	rv := reflect.ValueOf(v)
	for {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		} else if rv.Kind() == reflect.Struct {
			break
		} else {
			return nil, fmt.Errorf("v must be struct or struct ptr, %v", reflect.TypeOf(rv))
		}
	}

	m := make(map[string]string)
	for i := 0; i < rv.NumField(); i++ {
		// 获取标签名
		tagK := strings.Trim(rv.Type().Field(i).Tag.Get(p.tag), " ")
		if len(tagK) <= 0 {
			tagK = rv.Type().Field(i).Name
		} else if tagK == "-" {
			continue
		}

		// 获取字段值
		field := rv.Field(i)
		if field.IsNil() {
			continue
		}
		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		tagV := p.String(field)

		m[tagK] = tagV
	}
	return m, nil
}

func (p FormParser) String(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.Interface().(string)
	case reflect.Int:
		return strconv.Itoa(v.Interface().(int))
	case reflect.Int64:
		return strconv.FormatInt(v.Interface().(int64), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Interface().(bool))
	case reflect.Float64:
		return fmt.Sprintf("%v", v.Interface().(float64))
	case reflect.Slice, reflect.Array:
		arr := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			arr[i] = p.String(v.Index(i))
		}
		return strings.Join(arr, ",")
	case reflect.Int8:
		return strconv.FormatInt(int64(v.Interface().(int8)), 10)
	case reflect.Int16:
		return strconv.FormatInt(int64(v.Interface().(int16)), 10)
	case reflect.Int32:
		return strconv.FormatInt(int64(v.Interface().(int32)), 10)
	case reflect.Uint:
		return strconv.FormatUint(uint64(v.Interface().(uint)), 10)
	case reflect.Uint8:
		return strconv.FormatUint(uint64(v.Interface().(uint8)), 10)
	case reflect.Uint16:
		return strconv.FormatUint(uint64(v.Interface().(uint16)), 10)
	case reflect.Uint32:
		return strconv.FormatUint(uint64(v.Interface().(uint32)), 10)
	case reflect.Uint64:
		return strconv.FormatUint(v.Interface().(uint64), 10)
	case reflect.Float32:
		return fmt.Sprintf("%v", v.Interface().(float32))
	case reflect.Complex64:
		return fmt.Sprintf("%v", v.Interface().(complex64))
	case reflect.Complex128:
		return fmt.Sprintf("%v", v.Interface().(complex128))
	default:
		panic(fmt.Sprintf("Unknown type, %v", v.Kind()))
	}
}

func StringPtr(v string) *string {
	return &v
}

func Int64Ptr(v int64) *int64 {
	return &v
}
