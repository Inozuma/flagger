package flagger

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DefaultFlagSet is the default FlagSet to use with flagger.
// Set to flag.CommandLine by default.
var DefaultFlagSet = flag.CommandLine

// DefaultFlagTag is the default tag when calling Flag
var DefaultFlagTag = "flag"

func boolDefValue(str string) (bool, error) {
	if str == "" {
		return false, nil
	}

	return strconv.ParseBool(str)
}

func intDefValue(str string) (int, error) {
	if str == "" {
		return 0, nil
	}

	return strconv.Atoi(str)
}

func int64DefValue(str string) (int64, error) {
	if str == "" {
		return 0, nil
	}

	return strconv.ParseInt(str, 10, 64)
}

func uintDefValue(str string) (uint, error) {
	if str == "" {
		return 0, nil
	}

	ui, err := strconv.ParseUint(str, 10, 0)
	return uint(ui), err
}

func uint64DefValue(str string) (uint64, error) {
	if str == "" {
		return 0, nil
	}

	return strconv.ParseUint(str, 10, 64)
}

func float64DefValue(str string) (float64, error) {
	if str == "" {
		return 0, nil
	}

	return strconv.ParseFloat(str, 64)
}

func durationDefValue(str string) (time.Duration, error) {
	if str == "" {
		return time.Duration(0), nil
	}

	return time.ParseDuration(str)
}

// Flag will parse given structure and declare the necessary flags from it.
func Flag(s interface{}) error {
	if s == nil {
		return fmt.Errorf("flagger: Flag nil structure")
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("flagger: %q not a pointer", v.Kind())
	}

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("flagger: %q not a structure", v.Kind())
	}

	st := v.Type()

	for idx := 0; idx < st.NumField(); idx++ {
		field := st.Field(idx)
		fieldValue := v.Field(idx)

		// ignore if we can not set the field
		if !fieldValue.CanSet() {
			continue
		}

		tag := field.Tag.Get(DefaultFlagTag)

		// explicit ignore
		if tag == "-" {
			continue
		}

		// get flag options if any were set
		flagName, flagDefValue, flagUsage := field.Name, "", ""

		if tag != "" {
			opts := strings.Split(tag, ",")
			if len(opts) >= 1 {
				flagName = opts[0]
			}
			if len(opts) >= 2 {
				flagDefValue = opts[1]
			}
			if len(opts) >= 3 {
				flagUsage = opts[2]
			}
		}

		// declare flags for the current structure field
		value := fieldValue.Addr().Interface()

		// if it implements flag.Value, call FlagSet.Var
		if fvalue, ok := value.(flag.Value); ok {
			DefaultFlagSet.Var(fvalue, flagName, flagUsage)
			continue
		}

		// otherwise check for supported type by flag package
		switch fieldValue.Kind() {
		case reflect.Bool:
			defValue, err := boolDefValue(flagDefValue)
			if err != nil {
				return fmt.Errorf("flagger: could not get boolean default value for %q: %s", flagName, err)
			}
			DefaultFlagSet.BoolVar(value.(*bool), flagName, defValue, flagUsage)
		case reflect.Int:
			defValue, err := intDefValue(flagDefValue)
			if err != nil {
				return fmt.Errorf("flagger: could not get int default value for %q: %s", flagName, err)
			}
			DefaultFlagSet.IntVar(value.(*int), flagName, defValue, flagUsage)
		case reflect.Int64:
			// Special case for time.Duration
			if field.Type.String() == "time.Duration" {
				defValue, err := durationDefValue(flagDefValue)
				if err != nil {
					return fmt.Errorf("flagger: could not get duration default value for %q: %s", flagName, err)
				}
				DefaultFlagSet.DurationVar(value.(*time.Duration), flagName, defValue, flagUsage)
			} else {
				defValue, err := int64DefValue(flagDefValue)
				if err != nil {
					return fmt.Errorf("flagger: could not get int64 default value for %q: %s", flagName, err)
				}
				DefaultFlagSet.Int64Var(value.(*int64), flagName, defValue, flagUsage)
			}
		case reflect.Uint:
			defValue, err := uintDefValue(flagDefValue)
			if err != nil {
				return fmt.Errorf("flagger: could not get uint default value for %q: %s", flagName, err)
			}
			DefaultFlagSet.UintVar(value.(*uint), flagName, defValue, flagUsage)
		case reflect.Uint64:
			defValue, err := uint64DefValue(flagDefValue)
			if err != nil {
				return fmt.Errorf("flagger: could not get uint64 default value for %q: %s", flagName, err)
			}
			DefaultFlagSet.Uint64Var(value.(*uint64), flagName, defValue, flagUsage)
		case reflect.Float64:
			defValue, err := float64DefValue(flagDefValue)
			if err != nil {
				return fmt.Errorf("flagger: could not get float64 default value for %q: %s", flagName, err)
			}
			DefaultFlagSet.Float64Var(value.(*float64), flagName, defValue, flagUsage)
		case reflect.String:
			DefaultFlagSet.StringVar(value.(*string), flagName, flagDefValue, flagUsage)
		default:
			// if it's not a flag supported type, return an error
			return fmt.Errorf("flagger: Flag unsupported kind %q for field %q", fieldValue.Kind(), field.Name)
		}
	}

	return nil
}

// Parse will call for DefaultFlagSet.Parse(args).
// You should call this if DefaultFlagSet is not flag.CommandLine, otherwise simply use flag.Parse().
func Parse(args []string) error {
	return DefaultFlagSet.Parse(args)
}
