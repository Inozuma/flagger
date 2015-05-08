package flagger

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

// FlagTest defines a test for a structure and expected flag settings
type FlagTest struct {
	Generate      func() interface{} // structure to test
	ExpectedError bool
	ExpectedFlags []*ExpectedFlag
}

type ExpectedFlag struct {
	Name     string
	Value    string
	DefValue string
	Usage    string
}

type TestValue struct {
	s string
}

func (v *TestValue) Get() interface{} {
	return v.s
}

func (v *TestValue) Set(s string) error {
	v.s = s
	return nil
}

func (v *TestValue) String() string {
	return v.s
}

func TestFlag(t *testing.T) {
	tests := []FlagTest{
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					B      bool
					I      int
					I64    int64
					Ui     uint
					Ui64   uint64
					F      float64
					S      string
					T      time.Duration
					Custom TestValue
				}{}
			},
			ExpectedFlags: []*ExpectedFlag{
				&ExpectedFlag{
					Name:     "B",
					DefValue: "false",
					Value:    "true",
				},
				&ExpectedFlag{
					Name:     "I",
					DefValue: "0",
					Value:    "-42",
				},
				&ExpectedFlag{
					Name:     "I64",
					DefValue: "0",
					Value:    "-101010",
				},
				&ExpectedFlag{
					Name:     "Ui",
					DefValue: "0",
					Value:    "42",
				},
				&ExpectedFlag{
					Name:     "Ui64",
					DefValue: "0",
					Value:    "101010",
				},
				&ExpectedFlag{
					Name:     "F",
					DefValue: "0",
					Value:    "1.2",
				},
				&ExpectedFlag{
					Name:  "S",
					Value: "lorem ipsum",
				},
				&ExpectedFlag{
					Name:     "T",
					DefValue: "0",
					Value:    "42s",
				},
				&ExpectedFlag{
					Name:  "Custom",
					Value: "test value",
				},
			},
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					B      bool          `flag:"bool"`
					I      int           `flag:"int"`
					I64    int64         `flag:"int64"`
					Ui     uint          `flag:"uint"`
					Ui64   uint64        `flag:"uint64"`
					F      float64       `flag:"float"`
					S      string        `flag:"string"`
					T      time.Duration `flag:"duration"`
					Custom TestValue     `flag:"custom"`
				}{}
			},
			ExpectedFlags: []*ExpectedFlag{
				&ExpectedFlag{
					Name:     "bool",
					DefValue: "false",
					Value:    "true",
				},
				&ExpectedFlag{
					Name:     "int",
					DefValue: "0",
					Value:    "-42",
				},
				&ExpectedFlag{
					Name:     "int64",
					DefValue: "0",
					Value:    "-101010",
				},
				&ExpectedFlag{
					Name:     "uint",
					DefValue: "0",
					Value:    "42",
				},
				&ExpectedFlag{
					Name:     "uint64",
					DefValue: "0",
					Value:    "101010",
				},
				&ExpectedFlag{
					Name:     "float",
					DefValue: "0",
					Value:    "1.2",
				},
				&ExpectedFlag{
					Name:  "string",
					Value: "lorem ipsum",
				},
				&ExpectedFlag{
					Name:     "duration",
					DefValue: "0",
					Value:    "42s",
				},
				&ExpectedFlag{
					Name:  "custom",
					Value: "test value",
				},
			},
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					B      bool          `flag:"bool,false,boolean value"`
					I      int           `flag:"int,-42,integer value"`
					I64    int64         `flag:"int64,-101010,64-bit integer value"`
					Ui     uint          `flag:"uint,42,unsigned integer value"`
					Ui64   uint64        `flag:"uint64,101010,64-bit unsigned integer value"`
					F      float64       `flag:"float,1.2,float value"`
					S      string        `flag:"string,lorem ipsum,string value"`
					T      time.Duration `flag:"duration,42s,duration value"`
					Custom TestValue     `flag:"custom,,custom value"`
				}{}
			},
			ExpectedFlags: []*ExpectedFlag{
				&ExpectedFlag{
					Name:     "bool",
					DefValue: "false",
					Value:    "true",
					Usage:    "boolean value",
				},
				&ExpectedFlag{
					Name:     "int",
					DefValue: "-42",
					Value:    "-21",
					Usage:    "integer value",
				},
				&ExpectedFlag{
					Name:     "int64",
					DefValue: "-101010",
					Value:    "-100100100",
					Usage:    "64-bit integer value",
				},
				&ExpectedFlag{
					Name:     "uint",
					DefValue: "42",
					Value:    "21",
					Usage:    "unsigned integer value",
				},
				&ExpectedFlag{
					Name:     "uint64",
					DefValue: "101010",
					Value:    "100100100",
					Usage:    "64-bit unsigned integer value",
				},
				&ExpectedFlag{
					Name:     "float",
					DefValue: "1.2",
					Value:    "42.21",
					Usage:    "float value",
				},
				&ExpectedFlag{
					Name:     "string",
					DefValue: "lorem ipsum",
					Value:    "ipsum lorem",
					Usage:    "string value",
				},
				&ExpectedFlag{
					Name:     "duration",
					DefValue: "42s",
					Value:    "21h0m0s",
					Usage:    "duration value",
				},
				&ExpectedFlag{
					Name:  "custom",
					Value: "value test",
					Usage: "custom value",
				},
			},
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					B bool `flag:"bool,tru,boolean value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					I int `flag:"int,-42i,integer value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					I64 int64 `flag:"int64,-101010i,64-bit integer value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					Ui uint `flag:"uint,42x,unsigned integer value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					Ui64 uint64 `flag:"uint64,101010x,64-bit unsigned integer value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					F float64 `flag:"float,...,float value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					T time.Duration `flag:"duration,42ss,duration value"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					FieldToIgnore string `flag:"-"`
				}{}
			},
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					unexportedField string
				}{}
			},
		},
		FlagTest{
			Generate: func() interface{} {
				return &struct {
					UnsupportedType *string `flag:"unsupported"`
				}{}
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return nil
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				return 0
			},
			ExpectedError: true,
		},
		FlagTest{
			Generate: func() interface{} {
				var s = "random string"
				return &s
			},
			ExpectedError: true,
		},
	}

	for idx, test := range tests {
		t.Logf("test %d", idx)
		DefaultFlagSet = flag.NewFlagSet("test", flag.ContinueOnError)

		st := test.Generate()
		err := Flag(st)
		if err != nil {
			if test.ExpectedError {
				t.Logf("got expected error: %s", err)
				continue
			}
			t.Errorf("got unexpected error: %s", err)
			continue
		} else if test.ExpectedError {
			t.Errorf("expected an error")
			continue
		}

		var numFlags int
		DefaultFlagSet.VisitAll(func(f *flag.Flag) {
			t.Logf("got flag [name: %q, default: %q, usage: %q]", f.Name, f.DefValue, f.Usage)
			numFlags++
		})
		if len(test.ExpectedFlags) != numFlags {
			t.Errorf("expected %d flags, got %d", len(test.ExpectedFlags), numFlags)
			continue
		}

		var args []string
		for _, expectedFlag := range test.ExpectedFlags {
			arg := fmt.Sprintf("--%s=%s", expectedFlag.Name, expectedFlag.Value)
			args = append(args, arg)
		}

		if err := Parse(args); err != nil {
			t.Errorf("Parse got unexpected error: %s", err)
		}

		for _, expectedFlag := range test.ExpectedFlags {
			fl := DefaultFlagSet.Lookup(expectedFlag.Name)

			if fl == nil {
				t.Errorf("got unexpected nil Flag")
				continue
			}

			if fl.Name != expectedFlag.Name {
				t.Errorf("expected flag name to be %q, got %q", expectedFlag.Name, fl.Name)
			}

			if fl.Usage != expectedFlag.Usage {
				t.Errorf("expected flag usage to be %q, got %q", expectedFlag.Usage, fl.Usage)
			}

			if fl.DefValue != expectedFlag.DefValue {
				t.Errorf("expected flag default value to be %q, got %q", expectedFlag.DefValue, fl.DefValue)
			}

			if fl.Value.String() != expectedFlag.Value {
				t.Errorf("expected flag value to be %s, got %s", expectedFlag.Value, fl.Value)
			}
		}
	}
}
