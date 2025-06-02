package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/katesclau/letisgo/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	ID    string `json:"id" dynamodbav:"sk"`
	Value string `json:"value" dynamodbav:"value,omitempty"`
	Typ   string `json:"" dynamodbav:"pk"`
}

func NewTestStruct(v string) testStruct {
	return testStruct{
		ID:    uuid.NewString(),
		Value: v,
		Typ:   "TestStruct",
	}
}

func (ts *testStruct) GetType() string {
	return ts.Typ
}

func Test_DB_Provider(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	s := newTestStage(t)

	t.Run("Should be able insert item in DB", func(t *testing.T) {
		ctx := context.Background()

		o, err := s.ddb.Insert(ctx, NewTestStruct("Some text"))
		require.Nil(t, err)

		require.NotNil(t, o)
	})

	t.Run("Should be able get item in DB", func(t *testing.T) {
		ctx := context.Background()

		test := NewTestStruct("Some string")

		o, err := s.ddb.Insert(ctx, test)
		require.Nil(t, err)
		require.NotNil(t, o)

		record, err := s.ddb.Get(ctx, test.GetType(), test.ID)
		require.Nil(t, err)
		require.Equal(t, test, record)
	})

	t.Run("Should be able to delete record, and return its values", func(t *testing.T) {
		ctx := context.Background()

		test := NewTestStruct("Some string to delete")

		// Insert the record to be deleted
		o, err := s.ddb.Insert(ctx, test)
		require.Nil(t, err)
		require.NotNil(t, o)

		// Delete the record
		deletedRecord, err := s.ddb.Delete(ctx, test.GetType(), test.ID)
		require.Nil(t, err)
		require.NotNil(t, deletedRecord)

		// Verify the deleted record matches the original
		require.Equal(t, test, deletedRecord)

		// Attempt to retrieve the deleted record
		_, err = s.ddb.Get(ctx, reflect.TypeOf(test).Name(), test.ID)
		require.NotNil(t, err)

	})

	t.Run("Should be able to batch get items in DB", func(t *testing.T) {
		ctx := context.Background()

		// Create test records
		test1 := NewTestStruct("Batch item 1")
		test2 := NewTestStruct("Batch item 2")

		// Insert test records
		_, err := s.ddb.Insert(ctx, test1)
		require.Nil(t, err)
		_, err = s.ddb.Insert(ctx, test2)
		require.Nil(t, err)

		// Perform batch get
		keys := []db.KeyValues{
			{
				Pk: "TestStruct",
				Sk: test1.ID,
			},
			{
				Pk: "TestStruct",
				Sk: test2.ID,
			},
		}
		records, err := s.ddb.BatchGet(ctx, keys)
		require.Nil(t, err)
		require.Len(t, records, 2)

		// Verify the retrieved records match the originals
		require.Equal(t, test1, records[0])
		require.Equal(t, test2, records[1])
	})

	t.Run("Should be able to update item in DB", func(t *testing.T) {
		ctx := context.Background()

		// Create and insert a test record
		original := NewTestStruct("Original value")
		o, err := s.ddb.Insert(ctx, original)
		require.Nil(t, err)
		require.NotNil(t, o)

		// Update the record
		updated := original
		updated.Value = "Updated value"
		err = s.ddb.Update(ctx, updated)
		require.Nil(t, err)

		// Retrieve the updated record and verify it matches
		retrievedRecord, err := s.ddb.Get(ctx, updated.GetType(), updated.ID)
		require.Nil(t, err)
		require.Equal(t, updated, retrievedRecord)
	})
}
