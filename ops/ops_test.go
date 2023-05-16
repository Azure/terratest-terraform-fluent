package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		err := mock.Query("not_exist").HasValue(nil).AsError()
		assert.NoError(t, err)
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

		actual := "test"
		mock := mockOperativeType(actual)
		mock.Exist = false
		err := mock.HasValue("not_test").AsError()
		assert.ErrorContains(t, err, "test_resource.test_key: not found when expected")
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

func mockOperativeType(val any) Operative {
	return Operative{
		Reference: "test_resource.test_key",
		Actual:    val,
		Exist:     true,
	}
}
