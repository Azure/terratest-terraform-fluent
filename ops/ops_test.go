package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Azure/terratest-terraform-fluent/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		actual := map[string]interface{}{
			"test_map_key": "test",
		}
		mock := mockOperativeType(actual)
		err := mock.Query("test_map_key").HasValue("test").AsError()
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		mock := mockOperativeType(any(nil))
		mock.Exist = false
		err := mock.Query("test_map_key.test_nested_key").HasValue("test").AsError()
		assert.ErrorContains(t, err, "not found when expected")
	})

	t.Run("NestedMap", func(t *testing.T) {
		t.Parallel()
		actual := map[string]interface{}{
			"test_map_key": map[string]interface{}{
				"test_nested_key": "test",
			},
		}
		mock := mockOperativeType(actual)
		err := mock.Query("test_map_key.test_nested_key").HasValue("test").AsError()
		assert.NoError(t, err)
	})

	t.Run("QueryReturnNil", func(t *testing.T) {
		t.Parallel()

		actual := map[string]interface{}{
			"test_map_key": "test",
		}
		mock := mockOperativeType(actual)
		o := mock.Query("not_exist")
		err := o.HasValue("nil").AsError()
		assert.ErrorContains(t, err, "not found when expected")
		assert.False(t, o.Exist)
	})

	t.Run("QueryLength", func(t *testing.T) {
		t.Parallel()

		actual := map[string]any{
			"test_list": []any{
				"test",
				"test2",
			},
		}
		mock := mockOperativeType(actual)
		err := mock.Query("test_list.#").HasValue(2).AsError()
		assert.NoError(t, err)
	})

	t.Run("EscapedJsonString", func(t *testing.T) {
		t.Parallel()
		actual := "{\"properties\":{\"testProperty\":{\"testArray\":[\"testArrayMember\"]},\"testProperty2\":{\"TestArray2\":[]}}}"
		mock := mockOperativeType(actual)
		err := mock.Query("properties.testProperty.testArray.0").HasValue("testArrayMember").AsError()
		assert.NoError(t, err)
	})

	t.Run("InvalidJson", func(t *testing.T) {
		t.Parallel()
		mock := mockOperativeType("invalid json")
		err := mock.Query(".").HasValue("nil").AsError()
		assert.ErrorContains(t, err, "not valid JSON")
	})
}

func TestContainsJsonValue(t *testing.T) {
	t.Parallel()

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		mock := mockOperativeType(any(nil))
		mock.Exist = false
		err := mock.ContainsJsonValue(JsonAssertionFunc(func(input json.RawMessage) (*bool, error) { return to.Ptr(true), nil })).AsError()
		assert.ErrorContains(t, err, "not found when expected")
	})

	t.Run("NotAString", func(t *testing.T) {
		t.Parallel()
		mock := mockOperativeType(func() {})
		err := mock.ContainsJsonValue(JsonAssertionFunc(func(input json.RawMessage) (*bool, error) { return to.Ptr(true), nil })).AsError()
		assert.ErrorContains(t, err, "value is not a string")
	})

	t.Run("SuccessArray", func(t *testing.T) {
		t.Parallel()
		f := func(input json.RawMessage) (*bool, error) {
			i := make([]any, 0, 1)
			if err := json.Unmarshal(input, &i); err != nil {
				return nil, fmt.Errorf("JSON input is not an array")
			}
			if len(i) == 0 {
				return nil, fmt.Errorf("JSON input is empty")
			}
			if i[0].(map[string]any)["test"] != "test" {
				return nil, fmt.Errorf("JSON input key name is not equal to test")
			}

			return to.Ptr(true), nil
		}

		// Create a mock ThatTypeWithKey object
		mock := mockOperativeType("[{\"test\":\"test\"}]")
		err := mock.ContainsJsonValue(JsonAssertionFunc(f))
		require.NoError(t, err.AsError())
	})

	t.Run("AssertionFailed", func(t *testing.T) {
		t.Parallel()

		f := JsonAssertionFunc(
			func(input json.RawMessage) (*bool, error) {
				return to.Ptr(false), nil
			},
		)

		mock := mockOperativeType("{\"test\": \"test\"}")
		err := mock.ContainsJsonValue(f).AsError()
		assert.ErrorContains(t, err, "test_resource.test_key: assertion failed for")
	})

	t.Run("SuccessObject", func(t *testing.T) {
		t.Parallel()

		f := JsonAssertionFunc(
			func(input json.RawMessage) (*bool, error) {
				i := make(map[string]any)
				if err := json.Unmarshal(input, &i); err != nil {
					return nil, fmt.Errorf("JSON input is not an map")
				}
				if len(i) == 0 {
					return nil, fmt.Errorf("JSON input is empty")
				}
				if i["test"] != "test" {
					return to.Ptr(false), nil
				}
				return to.Ptr(true), nil
			},
		)

		mock := mockOperativeType("{\"test\": \"test\"}")
		err := mock.ContainsJsonValue(f).AsError()
		assert.NoError(t, err)
	})

	t.Run("EmptyKey", func(t *testing.T) {
		t.Parallel()

		f := JsonAssertionFunc(
			func(input json.RawMessage) (*bool, error) {
				return to.Ptr(true), nil
			},
		)

		mock := mockOperativeType("")
		err := mock.ContainsJsonValue(f).AsError()
		assert.ErrorContains(t, err, "test_resource.test_key: is empty")
	})

	t.Run("NilKey", func(t *testing.T) {
		t.Parallel()

		f := JsonAssertionFunc(
			func(input json.RawMessage) (*bool, error) {
				return to.Ptr(true), nil
			},
		)

		// Create a mock ThatTypeWithKey object
		mock := mockOperativeType(any(nil))
		err := mock.ContainsJsonValue(f).AsError()
		assert.ErrorContains(t, err, "test_resource.test_key: is empty")
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		f := JsonAssertionFunc(
			func(input json.RawMessage) (*bool, error) {
				return to.Ptr(true), errors.New("test error")
			},
		)

		// Create a mock ThatTypeWithKey object
		mock := mockOperativeType("{\"test\": \"test\"}")
		err := mock.ContainsJsonValue(f).AsError()
		assert.ErrorContains(t, err, "test error")
	})
}

func TestHasValue(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		actual := "test"
		mock := mockOperativeType(actual)
		err := mock.HasValue("test").AsError()
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		actual := "test"
		mock := mockOperativeType(actual)
		err := mock.HasValue("not_test").AsError()
		assert.ErrorContains(t, err, "expected value not_test not equal to actual test")
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		mock := mockOperativeType("test")
		mock.Exist = false
		err := mock.HasValue("not_test").AsError()
		assert.ErrorContains(t, err, "test_resource.test_key: not found when expected")
	})

	t.Run("InvalidArgs", func(t *testing.T) {
		t.Parallel()

		mock := mockOperativeType(func() {})
		err := mock.HasValue("test").AsError()
		assert.ErrorContains(t, err, "invalid operation")
	})
}

func TestContainsString(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		mock := mockOperativeType("contains string")
		err := mock.ContainsString("ains str").AsError()
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		mock := mockOperativeType("contains string")
		err := mock.ContainsString("not_found").AsError()
		assert.ErrorContains(t, err, "expected value not_found not contained within contains string")
	})

	t.Run("NotExists", func(t *testing.T) {
		mock := mockOperativeType(any(nil))
		mock.Exist = false
		err := mock.ContainsString("ains str").AsError()
		assert.ErrorContains(t, err, "not found when expected")
	})

	t.Run("NotAString", func(t *testing.T) {
		mock := mockOperativeType(any(nil))
		err := mock.ContainsString("").AsError()
		assert.ErrorContains(t, err, "Cannot convert value to string")
	})

}

func TestExists(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		mock := mockOperativeType(nil)
		err := mock.Exists().AsError()
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		mock := mockOperativeType(nil)
		mock.Exist = false
		err := mock.Exists().AsError()
		assert.ErrorContains(t, err, "not found when expected")
	})
}

func TestDoesNotExist(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		mock := mockOperativeType(nil)
		mock.Exist = false
		err := mock.DoesNotExist().AsError()
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		mock := mockOperativeType(nil)
		err := mock.DoesNotExist().AsError()
		assert.ErrorContains(t, err, "found when not expected")
	})
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		mock := mockOperativeType("test")
		val, err := mock.GetValue()
		assert.NoError(t, err)
		assert.Equal(t, "test", val)
	})

	t.Run("Failure", func(t *testing.T) {
		mock := mockOperativeType(any(nil))
		mock.Exist = false
		_, err := mock.GetValue()
		assert.ErrorContains(t, err, "not found when expected")
	})
}

func TestValidateEqualArgsFuncFail(t *testing.T) {
	t.Parallel()
	require.Nil(t, validateEqualArgs(nil, nil))
	f := func() {}
	assert.ErrorContains(t, validateEqualArgs(f, nil), "cannot take func type as argument")
	assert.ErrorContains(t, validateEqualArgs(nil, f), "cannot take func type as argument")
}

func TestIsFunction(t *testing.T) {
	t.Parallel()
	f := func() {}
	assert.True(t, isFunction(f))
	i := 1
	assert.False(t, isFunction(i))
	var s *string
	assert.False(t, isFunction(s))
}

func mockOperativeType(val any) Operative {
	return Operative{
		Reference: "test_resource.test_key",
		Actual:    val,
		Exist:     true,
	}
}
