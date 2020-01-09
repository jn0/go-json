// Internal representation of JSON values:
//   Scalars: integers, floats, bools, strings
//   Structural: arrays and objects
//   null: do not need special implementation
package json

// https://golangbot.com/interfaces-part-2/#implementinginterfacesusingpointerreceiversvsvaluereceivers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// any JSONable must have .Json() method to fetch a string representation of JSON
type JSONable interface {
	Json() string
}

// the very basic item of any JSON document, almost never being used "as is"
type JsonValue interface {
	JSONable                    // it has a mean to fetch JSON as a string
	Value() interface{}         // returns the value beneath this JsonValue
	Set(interface{}) JsonValue  // replaces the value of this item
	Parse(string) error         // parses the string and then .Set() result
	Append(interface{})         // extends a JsonArray
	Insert(string, interface{}) // updates a JsonObject
	Equal(JsonValue) bool       // compares two JsonValue to be equal
	IsNull() bool               // compares this JsonValue to be zero
}

/******************************************************************************/

type JsonInt int // one of the two JSON numerals, the integral one

// It is not meant to compare its value to 0
func (self *JsonInt) IsNull() bool { return self == nil || self == (*JsonInt)(nil) }

// nulls are equal, ints are sometimes equal, others aren't equal to int
func (self *JsonInt) Equal(v JsonValue) bool {
	switch v.(type) {
	case nil:
		return self.IsNull()
	case *JsonInt:
		if v.(*JsonInt).IsNull() {
			return self.IsNull()
		}
		return !self.IsNull() && self.Value() == v.(*JsonInt).Value()
	}
	return false
}

// the decimal representation is used
func (self *JsonInt) Json() string {
	if self.IsNull() {
		return "null"
	}
	return fmt.Sprintf("%d", *self)
}

// one can .Set() JsonInt from any Go int (not an uint) or from a string
func (self *JsonInt) Set(v interface{}) JsonValue {
	var t int
	switch iv := v.(type) {
	case int:
		t = iv
	case int8:
		t = int(v.(int8))
	case int16:
		t = int(v.(int16))
	case int32:
		t = int(v.(int32))
	case int64:
		t = int(v.(int64))
	case string:
		self.Parse(v.(string))
		return self
	default:
		panic(v)
	}
	*self = (JsonInt)(t)
	return self
}

// parse the string as decimal int64 and replace the current value
func (self *JsonInt) Parse(s string) error {
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return e
	}
	self.Set(v)
	return nil
}

// return the value as Go's int
func (self *JsonInt) Value() interface{}    { return int(*self) }
func (*JsonInt) Append(interface{})         { panic("Int is immutable") }
func (*JsonInt) Insert(string, interface{}) { panic("Int is immutable") }

// creates a new JsonInt from any compatible value (see the .Set() method)
func NewJsonInt(v interface{}) *JsonInt { return new(JsonInt).Set(v).(*JsonInt) }

/*----------------------------------------------------------------------------*/

type JsonFloat float64 // one of the two JSON numerals, the non-integral one

func (self *JsonFloat) IsNull() bool { return self == nil }

// nulls are equal, floats are sometimes equal, others aren't equal to float
func (self *JsonFloat) Equal(v JsonValue) bool {
	switch v.(type) {
	case nil:
		return self.IsNull()
	case *JsonFloat:
		if v.(*JsonFloat).IsNull() {
			return self.IsNull()
		}
		return !self.IsNull() && self.Value() == v.(*JsonFloat).Value()
	}
	return false
}

// the decimal fixed point representation is used
func (self *JsonFloat) Json() string {
	if self.IsNull() {
		return "null"
	}
	return fmt.Sprintf("%f", *self)
}

// one can .Set() JsonFloat from any Go float (float32 or float64) or from a string
func (self *JsonFloat) Set(v interface{}) JsonValue {
	var t float64
	switch iv := v.(type) {
	case float32:
		t = float64(v.(float32))
	case float64:
		t = iv
	case string:
		self.Parse(v.(string))
		return self
	default:
		panic(v)
	}
	*self = (JsonFloat)(t)
	return self
}

// return the value as Go's float64
func (self *JsonFloat) Value() interface{} { return float64(*self) }

// parse the string as decimal float64 and replace the current value
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

// creates a new JsonFloat from any compatible value (see the .Set() method)
func NewJsonFloat(v interface{}) *JsonFloat { return new(JsonFloat).Set(v).(*JsonFloat) }

/*----------------------------------------------------------------------------*/
type JsonBool bool

var boolStringValues = map[string]bool{
	"true":  true,
	"false": false,
}

func (self *JsonBool) IsNull() bool { return self == nil }
func (self *JsonBool) Equal(v JsonValue) bool {
	switch v.(type) {
	case nil:
		return self.IsNull()
	case *JsonBool:
		if v.(*JsonBool).IsNull() {
			return self.IsNull()
		}
		return !self.IsNull() && self.Value() == v.(*JsonBool).Value()
	}
	return false
}
func (self *JsonBool) Json() string {
	if self.IsNull() {
		return "null"
	}
	return fmt.Sprintf("%v", *self)
}
func (self *JsonBool) Set(v interface{}) JsonValue {
	switch v.(type) {
	case bool:
		*self = (JsonBool)(v.(bool))
	case string:
		self.Parse(v.(string))
	default:
		panic(v)
	}
	return self
}
func (self *JsonBool) Value() interface{} { return bool(*self) }
func (self *JsonBool) Parse(s string) error {
	v, found := boolStringValues[strings.ToLower(strings.TrimSpace(s))]
	if !found {
		panic(fmt.Sprintf("Bool: bad literal %+q", s))
	}
	self.Set(v)
	return nil
}
func (*JsonBool) Append(interface{})         { panic("Bool is immutable") }
func (*JsonBool) Insert(string, interface{}) { panic("Bool is immutable") }

func NewJsonBool(v interface{}) *JsonBool { return new(JsonBool).Set(v).(*JsonBool) }

/*----------------------------------------------------------------------------*/
type JsonString string

func (self *JsonString) IsNull() bool { return self == nil || (*self) == "" }
func (self *JsonString) Equal(v JsonValue) bool {
	switch v.(type) {
	case nil:
		return self.IsNull()
	case *JsonString:
		if v.(*JsonString).IsNull() {
			return self.IsNull()
		}
		return !self.IsNull() && self.Value() == v.(*JsonString).Value()
	}
	return false
}
func (self *JsonString) Json() string {
	if self.IsNull() {
		return "null"
	}
	return fmt.Sprintf("%q", *self)
}
func (self *JsonString) Set(v interface{}) JsonValue {
	switch v.(type) {
	case string:
		*self = (JsonString)(v.(string))
	case *JsonString:
		oth := v.(*JsonString)
		self.Set(fmt.Sprintf("%s", *oth)) // not the best conversion...
	default:
		panic(fmt.Sprintf("cannot %T.Set(%T)", self, v))
	}
	return self
}
func (self *JsonString) Value() interface{} { return string(*self) }
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

func NewJsonString(v interface{}) *JsonString { return new(JsonString).Set(v).(*JsonString) }

/******************************************************************************/

type JsonArray []JsonValue

func (self *JsonArray) IsNull() bool { return self == nil || (*self) == nil }
func (self *JsonArray) Equal(v JsonValue) bool {
	switch v.(type) {
	case nil:
		return self.IsNull()
	case *JsonArray:
		var other *JsonArray = v.(*JsonArray)
		if other.IsNull() {
			return self.IsNull()
		}
		if self.IsNull() {
			return false
		}
		if len(*self) != len(*other) {
			return false
		}
		for i, v := range *self {
			o := (*other)[i]
			if (v == nil || v.IsNull()) && (o == nil || o.IsNull()) {
				continue
			}
			if v == nil || o == nil {
				return false
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
	if self.IsNull() {
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
func (self *JsonArray) Set(v interface{}) JsonValue {
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
	return self
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

func NewJsonArray(v interface{}) *JsonArray { return new(JsonArray).Set(v).(*JsonArray) }

/*----------------------------------------------------------------------------*/
type JsonObject map[string]JsonValue

func (self *JsonObject) IsNull() bool { return self == nil || (*self) == nil }

func cmpMap(m1, m2 map[string]JsonValue) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		o, ok := m2[k]
		if !ok {
			return false
		}
		if (v == nil || v.IsNull()) && (o == nil || o.IsNull()) {
			continue
		}
		if !v.Equal(o) {
			return false
		}
	}
	return true
}

func (self *JsonObject) Equal(v JsonValue) bool {
	switch v.(type) {
	case nil:
		return self.IsNull()
	case *JsonObject:
		var other *JsonObject = v.(*JsonObject)
		if other.IsNull() {
			return self.IsNull()
		}
		if self.IsNull() {
			return false
		}
		if cmpMap(*self, *other) && cmpMap(*other, *self) {
			return true
		}
	}
	return false
}
func (self *JsonObject) Json() string {
	if self == nil || *self == nil {
		return "null"
	}
	var r, keys []string
	for k, _ := range *self {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		o := (*self)[k]
		var v string
		if o == nil {
			v = fmt.Sprintf("%q: null", k)
		} else {
			v = fmt.Sprintf("%q: %s", k, (o.(JsonValue)).Json())
		}
		r = append(r, v)
	}
	return "{ " + strings.Join(r, ", ") + " }"
}
func (self *JsonObject) Set(v interface{}) JsonValue {
	switch v.(type) {
	case *JsonObject:
		*self = *(v.(*JsonObject))
	case JsonObject:
		*self = v.(JsonObject)
	case map[string]JsonValue:
		*self = nil
		*self = make(map[string]JsonValue)
		for k, v := range v.(map[string]JsonValue) {
			self.Insert(k, v)
		}
	default:
		panic(v)
	}
	return self
}
func (self *JsonObject) Value() interface{} { return map[string]JsonValue(*self) }
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
		switch v.(type) {
		case *JsonObject:
			if self == v.(*JsonObject) {
				panic("Ooops!")
			}
		}
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
