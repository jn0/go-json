// Internal representation of JSON values:
//   Scalars: integers, floats, bools, strings
//   Structural: arrays and objects
//   null: do not need special implementation
package json

// https://golangbot.com/interfaces-part-2/#implementinginterfacesusingpointerreceiversvsvaluereceivers

import (
	"fmt"
	"strconv"
	"strings"
)

type JSONable interface {
	Json() string
}

type JsonValue interface {
	JSONable
	Value() interface{}
	Set(interface{})
	Parse(string) error
	Append(interface{})
	Insert(string, interface{})
	Equal(JsonValue) bool
}

/******************************************************************************/

type JsonInt int

func (self *JsonInt) Equal(v JsonValue) bool {
	switch v.(type) {
	case *JsonInt:
		return self.Value() == v.(*JsonInt).Value()
	}
	return false
}
func (self *JsonInt) Json() string {
	if self == nil {
		return "null"
	}
	return fmt.Sprintf("%d", *self)
}
func (self *JsonInt) Set(v interface{}) {
	var t int
	switch iv := v.(type) {
	case int:
		t = iv
	case int32:
		t = int(v.(int32))
	case int64:
		t = int(v.(int64))
	default:
		panic(v)
	}
	*self = (JsonInt)(t)
}
func (self *JsonInt) Parse(s string) error {
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return e
	}
	self.Set(v)
	return nil
}
func (self *JsonInt) Value() interface{}    { return int(*self) }
func (*JsonInt) Append(interface{})         { panic("Int is immutable") }
func (*JsonInt) Insert(string, interface{}) { panic("Int is immutable") }

func NewJsonInt(v interface{}) *JsonInt {
	r := new(JsonInt)
	switch v.(type) {
	case string:
		r.Parse(v.(string))
	default:
		r.Set(v)
	}
	return r
}

/*----------------------------------------------------------------------------*/
type JsonFloat float64

func (self *JsonFloat) Equal(v JsonValue) bool {
	switch v.(type) {
	case *JsonFloat:
		return self.Value() == v.(*JsonFloat).Value()
	}
	return false
}
func (self *JsonFloat) Json() string {
	if self == nil {
		return "null"
	}
	return fmt.Sprintf("%f", *self)
}
func (self *JsonFloat) Set(v interface{}) {
	var t float64
	switch iv := v.(type) {
	case float32:
		t = float64(v.(float32))
	case float64:
		t = iv
	default:
		panic(v)
	}
	*self = (JsonFloat)(t)
}
func (self *JsonFloat) Value() interface{} { return float64(*self) }
func (self *JsonFloat) Parse(s string) error {
	v, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return e
	}
	*self = (JsonFloat)(v)
	return nil
}
func (*JsonFloat) Append(interface{})         { panic("Float is immutable") }
func (*JsonFloat) Insert(string, interface{}) { panic("Float is immutable") }

func NewJsonFloat(v interface{}) *JsonFloat {
	r := new(JsonFloat)
	switch v.(type) {
	case string:
		r.Parse(v.(string))
	default:
		r.Set(v)
	}
	return r
}

/*----------------------------------------------------------------------------*/
type JsonBool bool

func (self *JsonBool) Equal(v JsonValue) bool {
	switch v.(type) {
	case *JsonBool:
		return self.Value() == v.(*JsonBool).Value()
	}
	return false
}
func (self *JsonBool) Json() string {
	if self == nil {
		return "null"
	}
	return fmt.Sprintf("%v", *self)
}
func (self *JsonBool) Set(v interface{})  { *self = (JsonBool)(v.(bool)) }
func (self *JsonBool) Value() interface{} { return *self }
func (self *JsonBool) Parse(s string) error {
	v, e := strconv.ParseBool(s)
	if e != nil {
		return e
	}
	self.Set(v)
	return nil
}
func (*JsonBool) Append(interface{})         { panic("Bool is immutable") }
func (*JsonBool) Insert(string, interface{}) { panic("Bool is immutable") }

func NewJsonBool(v interface{}) *JsonBool {
	r := new(JsonBool)
	switch v.(type) {
	case string:
		r.Parse(v.(string))
	default:
		r.Set(v)
	}
	return r
}

/*----------------------------------------------------------------------------*/
type JsonString string

func (self *JsonString) Equal(v JsonValue) bool {
	switch v.(type) {
	case *JsonString:
		return self.Value() == v.(*JsonString).Value()
	}
	return false
}
func (self *JsonString) Json() string {
	if self == nil || (*self) == "" {
		return "null"
	}
	return fmt.Sprintf("%q", *self)
}
func (self *JsonString) Set(v interface{})  {
	switch v.(type) {
	case string:
		*self = (JsonString)(v.(string))
	case *JsonString:
		oth := v.(*JsonString)
		self.Set(fmt.Sprintf("%s", *oth)) // not the best conversion...
	default:
		panic(fmt.Sprintf("cannot %T.Set(%T)", self, v))
	}
}
func (self *JsonString) Value() interface{} { return *self }
func (self *JsonString) Parse(s string) error {
	obj, tail, err := parseString(s)
	if err != nil {
		return err
	}
	if tail != "" {
		return SyntaxError(fmt.Errorf("Bad string %+q", s))
	}
	self.Set(obj)
	return nil
}
func (self *JsonString) Append(interface{})         { panic("String is immutable") }
func (self *JsonString) Insert(string, interface{}) { panic("String is immutable") }

func NewJsonString(v string) *JsonString {
	r := new(JsonString)
	r.Set(v)
	return r
}

/******************************************************************************/

type JsonArray []JsonValue

func (self *JsonArray) Equal(v JsonValue) bool {
	switch v.(type) {
	case *JsonArray:
		var other *JsonArray = v.(*JsonArray)
		if self == nil && other == nil {
			return true
		}
		if self == nil || other == nil {
			if self != other {
				return false
			}
		}
		if len(*self) != len(*other) {
			return false
		}
		for i, v := range *self {
			o := (*other)[i]
			if v == nil || o == nil {
				if v != o {
					return false
				}
				continue
			}
			if !v.Equal(o) {
				return false
			}
		}
		return true
	}
	return false
}
func (self *JsonArray) Json() string {
	if self == nil || *self == nil {
		return "null"
	}
	var r []string
	for _, o := range *self {
		if o == nil {
			r = append(r, "null")
		} else {
			r = append(r, (o.(JsonValue)).Json())
		}
	}
	return "[ " + strings.Join(r, ", ") + " ]"
}
func (self *JsonArray) Set(v interface{}) {
	switch v.(type) {
	case *JsonArray:
		*self = *(v.(*JsonArray))
	case []JsonValue:
		for _, x := range v.([]JsonValue) {
			self.Append(x)
		}
	default:
		panic("cannot")
	}
}
func (self *JsonArray) Value() interface{} { return *self }
func (self *JsonArray) Parse(s string) error {
	obj, tail, err := parseArray(s)
	if err != nil {
		return err
	}
	if tail != "" {
		return SyntaxError(fmt.Errorf("Bad array %+q", s))
	}
	self.Set(obj)
	return nil
}
func (self *JsonArray) Append(v interface{}) {
	if v == nil {
		*self = append(*self, nil)
	} else {
		*self = append(*self, v.(JsonValue))
	}
}
func (*JsonArray) Insert(string, interface{}) { panic("arrays are not insertable") }

func NewJsonArray(o []JsonValue) *JsonArray {
	r := new(JsonArray)
	r.Set(o)
	return r
}

/*----------------------------------------------------------------------------*/
type JsonObject map[string]JsonValue

func (self *JsonObject) Equal(v JsonValue) bool {
	switch v.(type) {
	case *JsonObject:
		var other *JsonObject = v.(*JsonObject)
		if self == nil && other == nil {
			return true
		}
		if self == nil || other == nil {
			if self != other {
				return false
			}
		}
		for k, v := range *self {
			o, ok := (*other)[k]
			if v == nil || o == nil {
				if v != o {
					return false
				}
				continue
			}
			ok = ok && v.Equal(o)
			if !ok {
				return false
			}
		}
		for k, v := range *other {
			o, ok := (*self)[k]
			if v == nil || o == nil {
				if v != o {
					return false
				}
				continue
			}
			ok = ok && v.Equal(o)
			if !ok {
				return false
			}
		}
		return true
	}
	return false
}
func (self *JsonObject) Json() string {
	if self == nil || *self == nil {
		return "null"
	}
	var r []string
	var v string
	for k, o := range *self {
		if o == nil {
			v = fmt.Sprintf("%q: null", k)
		} else {
			v = fmt.Sprintf("%q: %s", k, (o.(JsonValue)).Json())
		}
		r = append(r, v)
	}
	return "{ " + strings.Join(r, ", ") + " }"
}
func (self *JsonObject) Set(v interface{}) {
	*self = *(v.(*JsonObject))
}
func (self *JsonObject) Value() interface{} { return *self }
func (self *JsonObject) Parse(s string) error {
	obj, tail, err := parseObject(s)
	if err != nil {
		return err
	}
	if tail != "" {
		return SyntaxError(fmt.Errorf("Bad object %+q", s))
	}
	self.Set(obj)
	return nil
}
func (JsonObject) Append(v interface{}) { panic("objects are not appendable") }
func (self *JsonObject) Insert(n string, v interface{}) {
	if *self == nil {
		*self = make(map[string]JsonValue)
	}
	if v == nil {
		(*self)[n] = nil
	} else {
		(*self)[n] = v.(JsonValue)
	}
}

func NewJsonObject(o map[string]JsonValue) *JsonObject {
	r := new(JsonObject)
	for k, v := range o {
		r.Insert(k, v)
	}
	return r
}

/* EOF */
