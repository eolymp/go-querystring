package querystring

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var defaultOpts = UnmarshalOptions{ProtoOptions: protojson.UnmarshalOptions{DiscardUnknown: true}}
var strictOpts = UnmarshalOptions{ProtoOptions: protojson.UnmarshalOptions{}}

func Unmarshal(q url.Values, msg proto.Message) error {
	return defaultOpts.Unmarshal(q, msg)
}

func UnmarshalStrict(q url.Values, msg proto.Message) error {
	return strictOpts.Unmarshal(q, msg)
}

type UnmarshalOptions struct {
	ProtoOptions protojson.UnmarshalOptions
}

func (u UnmarshalOptions) Unmarshal(query url.Values, msg proto.Message) error {
	// backwards compatibility: entire message is specified in q parameter
	if q := query["q"]; len(q) > 0 {
		if q[0] == "" {
			return nil
		}

		return u.ProtoOptions.Unmarshal([]byte(q[0]), msg)
	}

	ref := msg.ProtoReflect()
	fields := ref.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)

		if err := u.setField(ref, field, query[string(field.Name())]); err != nil {
			return fmt.Errorf("invalid value for %s: %w", field.Name(), err)
		}
	}

	return nil
}

// setField sets a single proto field from query values
func (u UnmarshalOptions) setField(msg protoreflect.Message, field protoreflect.FieldDescriptor, values []string) error {
	if len(values) == 0 {
		return nil
	}

	switch {
	case field.IsList():
		list := msg.Mutable(field).List()

		for _, value := range values {
			val, err := u.toProtoValue(msg, field, value)
			if err != nil {
				return err
			}

			list.Append(val)
		}

	default:
		value, err := u.toProtoValue(msg, field, values[len(values)-1])
		if err != nil {
			return err
		}

		msg.Set(field, value)
	}

	return nil
}

// toProtoValue converts a string value to the appropriate protoreflect.Value
func (u UnmarshalOptions) toProtoValue(msg protoreflect.Message, field protoreflect.FieldDescriptor, value string) (protoreflect.Value, error) {
	switch field.Kind() {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(value), nil

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		val, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as int32: %w", value, err)
		}
		return protoreflect.ValueOfInt32(int32(val)), nil

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as int64: %w", value, err)
		}
		return protoreflect.ValueOfInt64(val), nil

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		val, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as uint32: %w", value, err)
		}
		return protoreflect.ValueOfUint32(uint32(val)), nil

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as uint64: %w", value, err)
		}

		return protoreflect.ValueOfUint64(val), nil

	case protoreflect.FloatKind:
		val, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as float32: %w", value, err)
		}

		return protoreflect.ValueOfFloat32(float32(val)), nil

	case protoreflect.DoubleKind:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as float64: %w", value, err)
		}

		return protoreflect.ValueOfFloat64(val), nil

	case protoreflect.BoolKind:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as bool: %w", value, err)
		}

		return protoreflect.ValueOfBool(val), nil

	case protoreflect.EnumKind:
		// Try to parse as enum number first, then as enum name
		if val, err := strconv.ParseInt(value, 10, 32); err == nil {
			return protoreflect.ValueOfEnum(protoreflect.EnumNumber(val)), nil
		} else {
			opts := field.Enum().Values()
			for i := 0; i < opts.Len(); i++ {
				opt := opts.Get(i)
				if string(opt.Name()) == value {
					return protoreflect.ValueOfEnum(opt.Number()), nil
				}
			}

			return protoreflect.Value{}, fmt.Errorf("unknown enum value %q for field %s", value, field.Name())
		}

	case protoreflect.MessageKind:
		switch field.Message().FullName() {
		case "google.protobuf.Timestamp":
			ts, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return protoreflect.Value{}, fmt.Errorf("cannot parse %q as timestamp: %w", value, err)
			}

			return protoreflect.ValueOfMessage(timestamppb.New(ts).ProtoReflect()), nil

		default:
			var sub protoreflect.Message
			switch {
			case field.IsList():
				sub = msg.NewField(field).List().NewElement().Message()
			default:
				sub = msg.NewField(field).Message()
			}

			if err := u.ProtoOptions.Unmarshal([]byte(value), sub.Interface()); err != nil {
				return protoreflect.Value{}, fmt.Errorf("cannot unmarshal message field %s: %w", field.Name(), err)
			}

			return protoreflect.ValueOfMessage(sub), nil
		}

	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes([]byte(value)), nil

	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported proto field kind: %s", field.Kind())
	}
}
