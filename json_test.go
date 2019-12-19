package json

import "testing"

func assert(t *testing.T, e error, s string, a ...interface{}) {
	if e != nil {
		a = append(a, e)
		t.Fatalf(s+": %v", a...)
	}
}

func TestValues(t *testing.T) {

	var is = 12345
	i1 := NewJsonInt(is)
	i2 := new(JsonInt)
	i2.Parse(i1.Json())
	if !i1.Equal(i2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", i1.Json(), i2.Json())
	}
	if i1.Json() != i2.Json() {
		t.Errorf("array: parse != json: %q != %q", i1.Json(), i2.Json())
	}

	var fs = -3.1
	f1 := NewJsonFloat(fs)
	f2 := new(JsonFloat)
	f2.Parse(f1.Json())
	if !f1.Equal(f2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", f1.Json(), f2.Json())
	}
	if f1.Json() != f2.Json() {
		t.Errorf("array: parse != json: %q != %q", f1.Json(), f2.Json())
	}

	var bs = true
	b1 := NewJsonBool(bs)
	b2 := new(JsonBool)
	b2.Parse(b1.Json())
	if !b1.Equal(b2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", b1.Json(), b2.Json())
	}
	if b1.Json() != b2.Json() {
		t.Errorf("array: parse != json: %q != %q", b1.Json(), b2.Json())
	}

	var ss = "some string"
	s1 := NewJsonString(ss)
	s2 := new(JsonString)
	s2.Parse(s1.Json())
	if !s1.Equal(s2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", s1.Json(), s2.Json())
	}
	if s1.Json() != s2.Json() {
		t.Errorf("array: parse != json: %q != %q", s1.Json(), s2.Json())
	}


	var as = []JsonValue{
		NewJsonInt(123),
		NewJsonString("asdasdasd"),
		nil,
	}
	a1 := NewJsonArray(as)
	a1.Append(NewJsonFloat(-3.5))
	a2 := new(JsonArray)
	a2.Parse(a1.Json())
	if !a1.Equal(a2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", a1.Json(), a2.Json())
	}
	if a1.Json() != a2.Json() {
		t.Errorf("array: parse != json: %q != %q", a1.Json(), a2.Json())
	}

	var os = map[string]JsonValue{
		"one": NewJsonInt(123),
		"two": NewJsonString("asdasdasd"),
		"tree": nil,
	}
	o1 := NewJsonObject(os)
	o1.Insert("new", NewJsonFloat(-3.5))
	o2 := new(JsonObject)
	o2.Parse(o1.Json())
	if !o1.Equal(o2) {
		t.Errorf("array: !parse.Equal(json): %q != %q", o1.Json(), o2.Json())
	}

}

func TestParsers(t *testing.T) {

	test := func(
		s string,
		f func (string) (JsonValue, string, error),
		eval func(JsonValue, string, error) bool,
	) {
		v, tail, err := f(s)
		ok := eval(v, tail, err)
		x := t.Logf
		if !ok { x = t.Errorf }
		o := ""
		if v != nil { o = v.Json() }
		x("s=%q (v=%#v=%s tail=%#v err=%#v) ok=%v", s, v, o, tail, err, ok)
	}

	noerr := func(v JsonValue, t string, e error) bool { return e == nil }
	erratic := func(v JsonValue, t string, e error) bool { return e != nil }
	notail := func(v JsonValue, t string, e error) bool { return t == "" }
	taily := func(v JsonValue, t string, e error) bool { return t != "" }
	clean := func(v JsonValue, t string, e error) bool { return noerr(v,t,e) && notail(v,t,e) }

	test(`{ "simple": "object" }`, parseObject, clean)
	test(`{ "simple": "object" }, 1`, parseObject, taily)
	test(`{ "better": "object", "one": 1, "null": null, "neg": -2.5 }`, parseObject, clean)
	test(`{ "simple" = "bad object"}`, parseObject, erratic)
	test(`{ "simple": "bad object"`, parseObject, erratic)
	test(`{ "simple": "bad object }`, parseObject, erratic)
	test(`{ "simple": "bad object", }`, parseObject, erratic)
	test(`{ wrong: "object" }`, parseObject, erratic)

	test(`[ 1, "simple", true, "list" ]`, parseArray, clean)
	test(`  [ 1, "simple", true, "list" ], "tail"`, parseArray, taily)
	test(`[ null, "simple", false, "list", ]`, parseArray, erratic)
	test(`[ null, "simple", false, "list",,,0 ]`, parseArray, erratic)
	test(`[ null, "simple", false, "list"`, parseArray, erratic)

	test(`"simple\nstring"`, parseString, clean)
	test(`"simple\nstring", 123`, parseString, taily)
	test(`"simple wrong string`, parseString, erratic)
	test(`"simple wrong string\"`, parseString, erratic)

	test(`123`, parseNumber, clean)
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
	test(`false`, parseBool, clean)
	test(`true,xxx`, parseBool, taily)
	test(`trUe`, parseBool, erratic)

	test(`null`, parseNull, clean)
	test(`null,`, parseNull, taily)
	test(`nUll`, parseNull, erratic)
	
	test(`null`, ParseValue, clean)
	test(`{ "not": [{ "so": "simple" }, "object"] }`, ParseValue, clean)
	test(`[ 1, "simple", [ true, "list" ], null, -2.5 ]`, ParseValue, clean)
	test(`"simple\nstring"`, ParseValue, clean)
	test(`+0123.03210`, ParseValue, clean)
	test(`false`, ParseValue, clean)
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
	assert(t, err, "GetValue(%+q)", source)
	if tail != "" {
		t.Fatalf("%+q tail", tail)
	}
	t.Logf("%s", json.Json())
}

func BenchmarkAll(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < 1000; i++ {
		ParseValue(source)
	}
}

