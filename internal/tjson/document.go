// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tjson

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

// documentType represents BSON Document type.
type documentType types.Document

// tjsontype implements tjsontype interface.
func (doc *documentType) tjsontype() {}

// UnmarshalJSONWithSchema unmarshals the JSON data with given JSON Schema.
func (doc *documentType) UnmarshalJSONWithSchema(data []byte, schema *Schema) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	if schema.Type != Object {
		panic(fmt.Sprintf("unexpected type %q", schema.Type))
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)

	var rawMessages map[string]json.RawMessage
	if err := dec.Decode(&rawMessages); err != nil {
		return lazyerrors.Error(err)
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Error(err)
	}

	b, ok := rawMessages["$k"]
	if !ok {
		return lazyerrors.Errorf("tjson.documentType.UnmarshalJSONWithSchema: missing $k")
	}

	var keys []string
	if err := json.Unmarshal(b, &keys); err != nil {
		return lazyerrors.Error(err)
	}
	if len(keys)+1 != len(rawMessages) {
		return lazyerrors.Errorf(
			"tjson.documentType.UnmarshalJSONWithSchema: %d elements in $k, %d in total",
			len(keys), len(rawMessages),
		)
	}

	td := must.NotFail(types.NewDocument())
	for _, key := range keys {
		b, ok = rawMessages[key]
		if !ok {
			return lazyerrors.Errorf("tjson.documentType.UnmarshalJSONWithSchema: missing key %q", key)
		}

		// If the field is set as null and is not present in the schema, it's a valid case.
		// If the field is set as something but null and is not present in the schema, we should return an error.
		s := schema.Properties[key]
		if s == nil && !bytes.Equal(b, []byte("null")) {
			return lazyerrors.Errorf("tjson.documentType.UnmarshalJSONWithSchema: no schema for key %q", key)
		}
		v, err := Unmarshal(b, s)
		if err != nil {
			return lazyerrors.Error(err)
		}
		if err = td.Set(key, v); err != nil {
			return lazyerrors.Error(err)
		}
	}

	*doc = documentType(*td)
	return nil
}

// MarshalJSON implements tjsontype interface.
func (doc *documentType) MarshalJSON() ([]byte, error) {
	td := types.Document(*doc)

	var buf bytes.Buffer

	buf.WriteString(`{"$k":`)
	keys := td.Keys()
	if keys == nil {
		keys = []string{}
	}
	b, err := json.Marshal(keys)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	buf.Write(b)

	for _, key := range keys {
		buf.WriteByte(',')

		if b, err = json.Marshal(key); err != nil {
			return nil, lazyerrors.Error(err)
		}
		buf.Write(b)
		buf.WriteByte(':')

		value, err := td.Get(key)
		if err != nil {
			return nil, lazyerrors.Error(err)
		}
		b, err := Marshal(value)
		if err != nil {
			return nil, lazyerrors.Error(err)
		}

		buf.Write(b)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// check interfaces
var (
	_ tjsontype = (*documentType)(nil)
)
