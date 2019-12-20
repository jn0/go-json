package json

import "testing"
import "github.com/stretchr/testify/assert"

import ( // for Example*
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func BenchmarkAll(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < 1000; i++ {
		ParseValue(source)
	}
}

func ExampleUptime() {
	// Suppose you have to feed some data to a monitor.
	// The API makes you to use JSON.
	// One of the monitored values is system uptime.

	// Read the data from `procfs(5)` (or run `uptime(1)` and catch its output)
	const uptime = "/proc/uptime"
	data, err := ioutil.ReadFile(uptime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ioutil.ReadFile(%+q): %v", uptime, err)
		return // handle file access/read error here
	}
	// It should arrive as something like "3294591.50 12242861.60\n",
	// where the first token is uptime (the 2nd one we wont use).

	word := strings.Split(strings.TrimSpace(string(data)), " ")
	if len(word) != 2 {
		fmt.Fprintf(os.Stderr, "Garbage read from %+q: %v", uptime, data)
		return // not much one can do if there is a "syntax error"...
	}

	// Use these values instead of "stable" fakes below
	now := time.Now().UTC().Unix() // an integer
	val := word[0]                 // a string (float will be parsed)

	now = 1576839878    // fake
	val = "3295164.960" // fake; still use string here, a "real" float would be ok too

	report := new(JsonObject)                  // The report.
	report.Insert("time", NewJsonInt(now))     // Add time stamp.
	report.Insert("uptime", NewJsonFloat(val)) // Add the value reported.

	fmt.Println(report.Json()) // Produce JSON text...
	// Output:
	// { "time": 1576839878, "uptime": 3295164.960000 }
}

func TestPanics(t *testing.T) {
	// The only parser that panics is parseString for wrong \uXXXX things.

	i1 := new(JsonInt)
	assert.Panics(t, func() { i1.Set(123.123) }, "Int.Set(Float)")
	assert.Panics(t, func() { i1.Append(123) }, "Int.Append()")
	assert.Panics(t, func() { i1.Insert("xyz", 123) }, "Int.Insert()")

	f1 := new(JsonFloat)
	assert.Panics(t, func() { f1.Set(123) }, "Float.Set(Int)")
	assert.Panics(t, func() { f1.Append(123.123) }, "Float.Append()")
	assert.Panics(t, func() { f1.Insert("xyz", 123.123) }, "Float.Insert()")

	b1 := new(JsonBool)
	assert.Panics(t, func() { b1.Set(123.123) }, "Bool.Set(Float)")
	assert.Panics(t, func() { b1.Parse("never") }, "Bool.Parse(garbage)")
	assert.Panics(t, func() { b1.Append(true) }, "Bool.Append()")
	assert.Panics(t, func() { b1.Insert("xyz", true) }, "Bool.Insert()")

	s1 := new(JsonString)
	assert.Panics(t, func() { s1.Set(123.123) }, "String.Set(Float)")
	assert.Panics(t, func() { s1.Append("123") }, "String.Append()")
	assert.Panics(t, func() { s1.Insert("xyz", "123") }, "String.Insert()")
	assert.Panics(t, func() { s1.Parse(`"zzz\u123zzz"`) }, `String.Parse("\u123z")`)

	a1 := new(JsonArray)
	assert.Panics(t, func() { a1.Set(123.123) }, "Array.Set(Float)")
	assert.Panics(t, func() { a1.Insert("xyz", 123) }, "Array.Insert()")

	o1 := new(JsonObject)
	assert.Panics(t, func() { o1.Set(123.123) }, "Object.Set(Float)")
	assert.Panics(t, func() { o1.Append(123) }, "Object.Append()")
	assert.Panics(t, func() { o1.Insert("xyz", o1) }, "Object looped")

	t.Logf("All panics performed")
}

func TestValues(t *testing.T) {

	var is = 12345
	i1 := NewJsonInt(is)
	i2 := new(JsonInt)
	if e := i2.Parse(i1.Json()); e != nil {
		t.Errorf("string: Parse(%+q): %v", i1.Json(), e)
	}
	if !i1.Equal(i2) {
		t.Errorf("int: !parse.Equal(json): %q != %q", i1.Json(), i2.Json())
	}
	if i1.Json() != i2.Json() {
		t.Errorf("int: parse != json: %q != %q", i1.Json(), i2.Json())
	}
	i1.Set(int8(123))
	i2.Set(int32(123))
	if !i1.Equal(i2) {
		t.Errorf("int: !Equal(): %q != %q", i1.Json(), i2.Json())
	}
	i1.Set(int(12388))
	i2.Set(int64(12388))
	if !i1.Equal(i2) {
		t.Errorf("int: !Equal(): %q != %q", i1.Json(), i2.Json())
	}
	i1.Set(int(12399))
	i2.Set("12399")
	if !i1.Equal(i2) {
		t.Errorf("int: !Equal(): %q != %q", i1.Json(), i2.Json())
	}

	var fs = -3.1
	f1 := NewJsonFloat(fs)
	f2 := new(JsonFloat)
	if e := f2.Parse(f1.Json()); e != nil {
		t.Errorf("string: Parse(%+q): %v", f1.Json(), e)
	}
	if !f1.Equal(f2) {
		t.Errorf("float: !parse.Equal(json): %q != %q", f1.Json(), f2.Json())
	}
	if f1.Json() != f2.Json() {
		t.Errorf("float: parse != json: %q != %q", f1.Json(), f2.Json())
	}
	f1.Set(123.2)
	f2.Set("123.2")
	if !f1.Equal(f2) ||
		f1.Value() != 123.2 || f2.Value() != 123.2 ||
		f1.Value() != f2.Value() {
		t.Errorf("float: !Equal(): %q != %q", f1.Json(), f2.Json())
	}
	f1.Set(float32(123.4))
	f2.Set(float64(123.4))
	if !f1.Equal(f2) {
		t.Logf("float: !Equal(): %q=%f != %q=%f (IT HAPPENS)",
			f1.Json(), f1.Value(), f2.Json(), f2.Value())
	}

	var bs = true
	b1 := NewJsonBool(bs)
	b2 := new(JsonBool)
	if e := b2.Parse(b1.Json()); e != nil {
		t.Errorf("string: Parse(%+q): %v", b1.Json(), e)
	}
	if !b1.Equal(b2) || !(b1.Value().(bool)) || !(b2.Value().(bool)) {
		t.Errorf("bool: !parse.Equal(json): %q != %q", b1.Json(), b2.Json())
	}
	if b1.Json() != b2.Json() {
		t.Errorf("bool: parse != json: %q != %q", b1.Json(), b2.Json())
	}
	b1.Set(true)
	b2.Set(!false)
	if !b1.Equal(b2) || !(b1.Value().(bool)) || !(b2.Value().(bool)) {
		t.Errorf("bool: !Equal(): %q != %q", b1.Json(), b2.Json())
	}
	b1.Set("true")
	b2.Set("false")
	if b1.Equal(b2) || !(b1.Value().(bool)) || (b2.Value().(bool)) {
		t.Errorf("bool: Equal(): %q != %q", b1.Json(), b2.Json())
	}

	var ss = "some string"
	s1 := NewJsonString(ss)
	s2 := new(JsonString)
	if e := s2.Parse(s1.Json()); e != nil {
		t.Errorf("string: Parse(%+q): %v", s1.Json(), e)
	}
	if !s1.Equal(s2) {
		t.Errorf("string: !parse.Equal(json): %q != %q", s1.Json(), s2.Json())
	}
	if s1.Json() != s2.Json() {
		t.Errorf("string: parse != json: %q != %q", s1.Json(), s2.Json())
	}
	if s1.Value() != ss {
		t.Errorf("string: %q != %q", s1.Json(), ss)
	}
	s1.Set("string one")
	s2.Set("string one")
	if !s1.Equal(s2) {
		t.Errorf("string: !Equal(): %q != %q", s1.Json(), s2.Json())
	}
	s2.Parse("\"string one\"")
	if !s1.Equal(s2) {
		t.Errorf("string: !Equal(): %q != %q", s1.Json(), s2.Json())
	}

	var as = []JsonValue{i1, b1, s1, nil} // don't use floats! they aren't equal
	a1 := NewJsonArray(as)
	a1.Append(NewJsonInt(-35))
	a2 := new(JsonArray)
	if e := a2.Parse(a1.Json()); e != nil {
		t.Errorf("string: Parse(%+q): %v", a1.Json(), e)
	}
	if !a1.Equal(a2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", a1.Json(), a2.Json())
	}
	if a1.Json() != a2.Json() {
		t.Errorf("array: parse != json: %q != %q", a1.Json(), a2.Json())
	}
	a1.Append(NewJsonInt(999))
	if a1.Equal(a2) {
		t.Errorf("array: Equal(): %q != %q", a1.Json(), a2.Json())
	}
	if a1.Json() == a2.Json() {
		t.Errorf("array: %q == %q", a1.Json(), a2.Json())
	}

	var os = map[string]JsonValue{
		"one":  i1,
		"two":  s1,
		"tree": nil,
		"list": a1,
		"bool": b1,
	}
	o1 := NewJsonObject(os)
	o1.Insert("new", NewJsonInt(-35))
	// o1.Insert("self", o1) // he-he...
	o2 := new(JsonObject)
	if e := o2.Parse(o1.Json()); e != nil {
		t.Errorf("string: Parse(%+q): %v", o1.Json(), e)
	}
	if !o1.Equal(o2) {
		t.Errorf("object: !parse.Equal(json): %q != %q", o1.Json(), o2.Json())
	}
	o2.Insert("other", o1)
	if o1.Equal(o2) {
		t.Errorf("object: Equal(): %q != %q", o1.Json(), o2.Json())
	}

}

func TestParsers(t *testing.T) {

	test := func(
		s string,
		f func(string) (JsonValue, string, error),
		eval func(JsonValue, string, error) bool,
	) {
		v, tail, err := f(s)
		ok := eval(v, tail, err)
		x := t.Logf
		if !ok {
			x = t.Errorf
		}
		o := ""
		if v != nil {
			o = v.Json()
		}
		x("s=%q (v=%#v=%s tail=%#v err=%#v) ok=%v", s, v, o, tail, err, ok)
	}

	noerr := func(v JsonValue, t string, e error) bool { return e == nil }
	erratic := func(v JsonValue, t string, e error) bool { return e != nil }
	notail := func(v JsonValue, t string, e error) bool { return t == "" }
	taily := func(v JsonValue, t string, e error) bool { return t != "" }
	clean := func(v JsonValue, t string, e error) bool { return noerr(v, t, e) && notail(v, t, e) }

	test(`bad`, parseObject, erratic)
	test(`	bad	`, parseObject, erratic)
	test(`	`, parseObject, erratic)
	test(`{ "simple": "object" }`, parseObject, clean)
	test(` 	{	"simple"	:	"object"	}	`, parseObject, clean)
	test(`{ "simple": "object" }, 1`, parseObject, taily)
	test(`{ "": "object" }, 1`, parseObject, taily)
	test(`{ "simple": "" }, 1`, parseObject, taily)
	test(`{ "better": "object", "one": 1, "null": null, "neg": -2.5 }`, parseObject, clean)
	test(`{ "simple" = "bad object"}`, parseObject, erratic)
	test(`{ "simple": "bad object"`, parseObject, erratic)
	test(`{ "simple": }`, parseObject, erratic)
	test(`{ "simple": `, parseObject, erratic)
	test(`{ "simple" `, parseObject, erratic)
	test(`{ "simple": "bad object }`, parseObject, erratic)
	test(`{ "simple": "bad object", }`, parseObject, erratic)
	test(`{ wrong: "object" }`, parseObject, erratic)

	test(`	`, parseArray, erratic)
	test(`	[	1,	"simple",    true,    "list"	]	`, parseArray, clean)
	test(`[ 1, "simple", true, "list" ]`, parseArray, clean)
	test(`  [ 1, "simple", true, "list" ], "tail"`, parseArray, taily)
	test(`[ null, "simple", false, "list", ]`, parseArray, erratic)
	test(`[ null, "simple", false, "list",,,0 ]`, parseArray, erratic)
	test(`[ null, "simple", false, "list"`, parseArray, erratic)
	test(`xxx`, parseArray, erratic)

	test(`"simple\nstring"`, parseString, clean)
	test(`"\nstring\twith\rescapes\u005c\u002Fyepp"`, parseString, clean)
	test(`	`, parseString, erratic)
	test(`	"simple\nstring"	`, parseString, clean)
	test(`	""	`, parseString, clean)
	test(`	" 		 "  	`, parseString, clean)
	test(`"simple\nstring", 123`, parseString, taily)
	test(`"simple wrong string`, parseString, erratic)
	test(`"simple wrong string\"`, parseString, erratic)

	test(``, parseNumber, erratic)
	test(`		`, parseNumber, erratic)
	test(`123`, parseNumber, clean)
	test(`	  123  	`, parseNumber, clean)
	test(`0123`, parseNumber, clean)
	test(`0123.03210`, parseNumber, clean)
	test(`+0123.03210`, parseNumber, clean)
	test(`-0123.03210`, parseNumber, clean)
	test(`x123`, parseNumber, erratic)
	test(`0x123`, parseNumber, taily)
	test(`123x`, parseNumber, taily)
	test(`123.321`, parseNumber, clean)
	test(`+987`, parseNumber, clean)
	test(`-321.`, parseNumber, clean)

	test(`true`, parseBool, clean)
	test(`	`, parseBool, erratic)
	test(` xxxx	`, parseBool, erratic)
	test(`false`, parseBool, clean)
	test(`true,xxx`, parseBool, taily)
	test(`trUe`, parseBool, erratic)

	test(`null`, parseNull, clean)
	test(`	`, parseNull, erratic)
	test(`  xxx  `, parseNull, erratic)
	test(`null,`, parseNull, taily)
	test(`nUll`, parseNull, erratic)

	test(`null`, ParseValue, clean)
	test(`{ "not": [{ "so": "simple" }, "object"] }`, ParseValue, clean)
	test(`[ 1, "simple", [ true, "list" ], null, -2.5 ]`, ParseValue, clean)
	test(`"simple\nstring"`, ParseValue, clean)
	test(`+0123.03210`, ParseValue, clean)
	test(`false`, ParseValue, clean)
	test(`  xxx  `, ParseValue, erratic)
	test(`	false	`, ParseValue, clean)
	test(`	`, ParseValue, erratic)
}

const source = `
{
  "update": true,
  "time": 1576745142824,
  "uptime": 2749251.96,
  "hostname": "jet-one",
  "loadavg": {
    "kernel": {
      "runnable": 5,
      "total": 297
    },
    "last_pid": 6810,
    "average": [
      1.01,
      1.03,
      1.11
    ]
  },
  "mounts": [
    {
      "spec": "/dev/root",
      "file": "/",
      "vfstype": "ext4",
      "mntops": [
        "rw",
        "relatime",
        "data=ordered"
      ],
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 27736252,
        "used_kb": 13165476,
        "available_kb": 13138816
      }
    },
    {
      "file": "/dev",
      "vfstype": "devtmpfs",
      "mntops": [
        "rw",
        "relatime",
        "size=7976116k",
        "nr_inodes=1994029",
        "mode=755"
      ],
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 7976116,
        "used_kb": 0,
        "available_kb": 7976116
      },
      "spec": "devtmpfs"
    },
    {
      "vfstype": "sysfs",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "sysfs",
      "file": "/sys"
    },
    {
      "size": null,
      "spec": "proc",
      "file": "/proc",
      "vfstype": "proc",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime"
      ],
      "freq": 0,
      "passno": 0
    },
    {
      "passno": 0,
      "size": {
        "used_kb": 4612,
        "available_kb": 8038300,
        "total_kb": 8042912
      },
      "spec": "tmpfs",
      "file": "/dev/shm",
      "vfstype": "tmpfs",
      "mntops": [
        "rw",
        "nosuid",
        "nodev"
      ],
      "freq": 0
    },
    {
      "spec": "devpts",
      "file": "/dev/pts",
      "vfstype": "devpts",
      "mntops": [
        "rw",
        "nosuid",
        "noexec",
        "relatime",
        "gid=5",
        "mode=620",
        "ptmxmode=000"
      ],
      "freq": 0,
      "passno": 0,
      "size": null
    },
    {
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 8042912,
        "used_kb": 833212,
        "available_kb": 7209700
      },
      "spec": "tmpfs",
      "file": "/run",
      "vfstype": "tmpfs",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "mode=755"
      ]
    },
    {
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "size=5120k"
      ],
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 5120,
        "used_kb": 12,
        "available_kb": 5108
      },
      "spec": "tmpfs",
      "file": "/run/lock",
      "vfstype": "tmpfs"
    },
    {
      "vfstype": "tmpfs",
      "mntops": [
        "ro",
        "nosuid",
        "nodev",
        "noexec",
        "mode=755"
      ],
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 8042912,
        "used_kb": 0,
        "available_kb": 8042912
      },
      "spec": "tmpfs",
      "file": "/sys/fs/cgroup"
    },
    {
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "xattr",
        "release_agent=/lib/systemd/systemd-cgroups-agent",
        "name=systemd"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/systemd"
    },
    {
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "pstore",
      "file": "/sys/fs/pstore",
      "vfstype": "pstore",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime"
      ]
    },
    {
      "passno": 0,
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/cpu,cpuacct",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "cpu",
        "cpuacct"
      ],
      "freq": 0
    },
    {
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "freezer"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/freezer",
      "vfstype": "cgroup"
    },
    {
      "passno": 0,
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/perf_event",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "perf_event"
      ],
      "freq": 0
    },
    {
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "debug"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/debug",
      "vfstype": "cgroup"
    },
    {
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "memory"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/memory"
    },
    {
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/pids",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "pids"
      ],
      "freq": 0,
      "passno": 0
    },
    {
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/net_cls,net_prio",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "net_cls",
        "net_prio"
      ],
      "freq": 0,
      "passno": 0,
      "size": null
    },
    {
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/cpuset",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "cpuset"
      ],
      "freq": 0,
      "passno": 0
    },
    {
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/blkio",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "blkio"
      ],
      "freq": 0,
      "passno": 0,
      "size": null
    },
    {
      "size": null,
      "spec": "cgroup",
      "file": "/sys/fs/cgroup/devices",
      "vfstype": "cgroup",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "noexec",
        "relatime",
        "devices"
      ],
      "freq": 0,
      "passno": 0
    },
    {
      "vfstype": "debugfs",
      "mntops": [
        "rw",
        "relatime"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "debugfs",
      "file": "/sys/kernel/debug"
    },
    {
      "vfstype": "mqueue",
      "mntops": [
        "rw",
        "relatime"
      ],
      "freq": 0,
      "passno": 0,
      "size": null,
      "spec": "mqueue",
      "file": "/dev/mqueue"
    },
    {
      "spec": "configfs",
      "file": "/sys/kernel/config",
      "vfstype": "configfs",
      "mntops": [
        "rw",
        "relatime"
      ],
      "freq": 0,
      "passno": 0,
      "size": null
    },
    {
      "vfstype": "tmpfs",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "relatime",
        "size=804292k",
        "mode=700",
        "uid=1001",
        "gid=1001"
      ],
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 804292,
        "used_kb": 16,
        "available_kb": 804276
      },
      "spec": "tmpfs",
      "file": "/run/user/1001"
    },
    {
      "freq": 0,
      "passno": 0,
      "size": {
        "total_kb": 804292,
        "used_kb": 0,
        "available_kb": 804292
      },
      "spec": "tmpfs",
      "file": "/run/user/1002",
      "vfstype": "tmpfs",
      "mntops": [
        "rw",
        "nosuid",
        "nodev",
        "relatime",
        "size=804292k",
        "mode=700",
        "uid=1002",
        "gid=1002"
      ]
    }
  ],
  "meminfo": {
    "MemFree_kB": 3776048,
    "MemAvailable_kB": 5085380,
    "Buffers_kB": 199708,
    "Cached_kB": 1904324,
    "SwapCached_kB": 0,
    "SwapTotal_kB": 0,
    "SwapFree_kB": 0,
    "MemTotal_kB": 8042912
  },
  "spools": [
    {
      "path": "/var/cache/jetson",
      "exists": true,
      "count": 0
    }
  ],
  "sensors": [
    {
      "adapter": "Virtual device",
      "values": [
        {
          "crit": 101,
          "name": "temp1",
          "input": 20.5
        }
      ],
      "sensor": "BCPU-therm-virtual-0"
    },
    {
      "adapter": "Virtual device",
      "values": [
        {
          "name": "temp1",
          "input": 20.5,
          "crit": 101
        }
      ],
      "sensor": "MCPU-therm-virtual-0"
    },
    {
      "sensor": "GPU-therm-virtual-0",
      "adapter": "Virtual device",
      "values": [
        {
          "crit": -40,
          "name": "temp1",
          "input": 19
        }
      ]
    },
    {
      "sensor": "Tboard_tegra-virtual-0",
      "adapter": "Virtual device",
      "values": [
        {
          "name": "temp1",
          "input": 18,
          "crit": 107
        }
      ]
    },
    {
      "sensor": "Tdiode_tegra-virtual-0",
      "adapter": "Virtual device",
      "values": [
        {
          "name": "temp1",
          "input": 17.75,
          "crit": 107
        }
      ]
    },
    {
      "sensor": "thermal-fan-est-virtual-0",
      "adapter": "Virtual device",
      "values": [
        {
          "name": "temp1",
          "input": 19.6
        }
      ]
    }
  ],
  "net": [
    {
      "response": {
        "status": 200,
        "content_type": "text/plain; charset=utf-8",
        "body": "31.173.81.31"
      },
      "url": "https://ifconfig.co/",
      "scheme": "https",
      "error": null
    }
  ]
}
`

func TestAll(t *testing.T) {
	json, tail, err := ParseValue(source)
	if err != nil {
		t.Errorf("GetValue(%+q): %v", source, err)
	}
	if tail != "" {
		t.Errorf("%+q tail", tail)
	}
	t.Logf("%s", json.Json())
}
