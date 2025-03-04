package textual_test

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"

	"cosmossdk.io/math"
	"cosmossdk.io/x/tx/textual"
)

func TestIntJsonTestcases(t *testing.T) {
	type integerTest []string
	var testcases []integerTest
	raw, err := os.ReadFile("./internal/testdata/integers.json")
	require.NoError(t, err)
	err = json.Unmarshal(raw, &testcases)
	require.NoError(t, err)

	textual := textual.NewTextual(nil)

	for _, tc := range testcases {
		t.Run(tc[0], func(t *testing.T) {
			// Parse test case strings as protobuf uint64
			i, err := strconv.ParseUint(tc[0], 10, 64)
			if err == nil {
				r, err := textual.GetFieldValueRenderer(fieldDescriptorFromName("UINT64"))
				require.NoError(t, err)

				checkNumberTest(t, r, protoreflect.ValueOf(i), tc[1])
			}

			// Parse test case strings as protobuf uint32
			i, err = strconv.ParseUint(tc[0], 10, 32)
			if err == nil {
				r, err := textual.GetFieldValueRenderer(fieldDescriptorFromName("UINT32"))
				require.NoError(t, err)

				checkNumberTest(t, r, protoreflect.ValueOf(i), tc[1])
			}

			// Parse test case strings as sdk.Ints
			_, ok := math.NewIntFromString(tc[0])
			if ok {
				r, err := textual.GetFieldValueRenderer(fieldDescriptorFromName("SDKINT"))
				require.NoError(t, err)

				checkNumberTest(t, r, protoreflect.ValueOf(tc[0]), tc[1])
			}
		})
	}
}

// checkNumberTest checks that the output of a number value renderer
// matches the expected string. Only use it to test numbers.
func checkNumberTest(t *testing.T, r textual.ValueRenderer, pv protoreflect.Value, expected string) {
	screens, err := r.Format(context.Background(), pv)
	require.NoError(t, err)
	require.Len(t, screens, 1)
	require.Equal(t, 0, screens[0].Indent)
	require.Equal(t, false, screens[0].Expert)

	require.Equal(t, expected, screens[0].Content)

	// Round trip.
	value, err := r.Parse(context.Background(), screens)
	require.NoError(t, err)

	v, err := math.LegacyNewDecFromStr(pv.String())
	require.NoError(t, err)

	v1, err := math.LegacyNewDecFromStr(value.String())
	require.NoError(t, err)

	require.True(t, v.Equal(v1))
}
