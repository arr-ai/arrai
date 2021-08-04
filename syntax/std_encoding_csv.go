package syntax

import (
	"bytes"
	"context"
	"io"

	"encoding/csv"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
)

func newDecodeConfig() decodeConfig {
	return decodeConfig{comma: ','}
}

func newEncodeConfig() encodeConfig {
	return encodeConfig{comma: ','}
}

type decodeConfig struct {
	fieldsPerRecord  int
	comma            rune
	comment          rune
	trimLeadingSpace bool
	lazyQuotes       bool
}

type encodeConfig struct {
	comma rune
	crlf  bool
}

func stdEncodingCSV() rel.Attr {
	return rel.NewTupleAttr(
		"csv",
		rel.NewNativeFunctionAttr(decodeAttr, func(_ context.Context, value rel.Value) (rel.Value, error) {
			return csvDecodeFnBody("csv.decode", value, newDecodeConfig())
		}),

		rel.NewNativeFunctionAttr(decoderAttr, func(_ context.Context, configValue rel.Value) (rel.Value, error) {
			fn := "csv.decoder"
			config := newDecodeConfig()

			configTuple, ok := configValue.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(configValue))
			}

			comma, err := getConfigInt(configTuple, fn, "comma", int(config.comma))
			if err != nil {
				return nil, err
			}
			config.comma = rune(comma)

			comment, err := getConfigInt(configTuple, fn, "comment", int(config.comment))
			if err != nil {
				return nil, err
			}
			config.comment = rune(comment)

			config.trimLeadingSpace, _ = getConfigBool(configTuple, "trimLeadingSpace")

			fieldsPerRecord, err := getConfigInt(configTuple, fn, "fieldsPerRecord", config.fieldsPerRecord)
			if err != nil {
				return nil, err
			}
			config.fieldsPerRecord = fieldsPerRecord

			config.lazyQuotes, _ = getConfigBool(configTuple, "lazyQuotes")

			return rel.NewNativeFunction("decode", func(_ context.Context, value rel.Value) (rel.Value, error) {
				return csvDecodeFnBody("csv.decoder payload", value, config)
			}), nil
		}),

		rel.NewNativeFunctionAttr("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
			return csvEncodeFnBody("csv.encode", value, newEncodeConfig())
		}),

		rel.NewNativeFunctionAttr("encoder", func(_ context.Context, configValue rel.Value) (rel.Value, error) {
			fn := "csv.encoder"
			config := newEncodeConfig()

			configTuple, ok := configValue.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(configValue))
			}

			comma, err := getConfigInt(configTuple, fn, "comma", int(config.comma))
			if err != nil {
				return nil, err
			}
			config.comma = rune(comma)

			config.crlf, _ = getConfigBool(configTuple, "crlf")

			return rel.NewNativeFunction("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
				return csvEncodeFnBody("csv.encoder payload", value, config)
			}), nil
		}),
	)
}

func csvDecodeFnBody(fn string, value rel.Value, config decodeConfig) (rel.Value, error) {
	var bs []byte
	switch t := value.(type) {
	case rel.String:
		bs = []byte(t.String())
	case rel.Bytes:
		bs = t.Bytes()
	default:
		return nil, errors.Errorf("first arg to %s must be string or bytes, not %s", fn, rel.ValueTypeAsString(value))
	}

	reader := csv.NewReader(bytes.NewBuffer(bs))
	reader.Comma = config.comma
	reader.Comment = config.comment
	reader.TrimLeadingSpace = config.trimLeadingSpace
	reader.FieldsPerRecord = config.fieldsPerRecord
	reader.LazyQuotes = config.lazyQuotes
	reader.ReuseRecord = true

	var rows []rel.Value
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make([]rel.Value, len(record))
		for j, item := range record {
			row[j] = rel.NewString([]rune(item))
		}
		rows = append(rows, rel.NewArray(row...))
	}
	return rel.NewArray(rows...), nil
}

func csvEncodeFnBody(fn string, value rel.Value, config encodeConfig) (rel.Value, error) {
	arr, ok := rel.AsArray(value)
	if !ok {
		return nil, errors.Errorf("first arg to %s must be array, not %s", fn, rel.ValueTypeAsString(value))
	}

	records := make([][]string, arr.Count())
	for i, row := range arr.Values() {
		rowArray, ok := rel.AsArray(row)
		if !ok {
			return nil, errors.Errorf("record %v must be array, not %v", i, rel.ValueTypeAsString(row))
		}
		record := make([]string, rowArray.Count())
		for j, value := range rowArray.Values() {
			s, ok := rel.AsString(value)
			if !ok {
				return nil, errors.Errorf("value %v of record %v must be string, not %v", j, i, rel.ValueTypeAsString(value))
			}
			record[j] = s.String()
		}
		records[i] = record
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	writer.Comma = config.comma
	writer.UseCRLF = config.crlf
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}

	return rel.NewBytes(buffer.Bytes()), nil
}
