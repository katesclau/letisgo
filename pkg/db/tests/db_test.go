package tests

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"mnesis.com/pkg/db"
)

type TestStruct struct {
	ID    string
	Value string
}

type TestProperties struct {
	Value string `json:"value" dynamodbav:"value,omitempty"`
}

type TestRecord struct {
	db.Record
	TestProperties
}

func (tr *TestRecord) GetStruct() TestStruct {
	return TestStruct{
		ID:    tr.Sk,
		Value: tr.Value,
	}
}

func (t *TestStruct) Record() any {
	// TODO make this type assertion available to all models
	typ := strings.Split(reflect.TypeOf(t).String(), ".")[1]
	r := TestRecord{
		db.Record{
			Pk:      typ,
			Sk:      t.ID,
			Type:    typ,
			Version: 0,
		},
		TestProperties{
			Value: t.Value,
		},
	}
	return r
}

func Test_DB_Provider(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	s := newTestStage(t)

	t.Run("Should be able insert item in DB", func(t *testing.T) {
		ctx := context.Background()

		o, err := s.ddb.Insert(ctx, &TestStruct{
			ID:    uuid.NewString(),
			Value: "Some string",
		})
		require.Nil(t, err)

		require.NotNil(t, o)
	})

	t.Run("Should be able get item in DB", func(t *testing.T) {
		ctx := context.Background()

		test := TestStruct{
			ID:    uuid.NewString(),
			Value: "Some string",
		}

		o, err := s.ddb.Insert(ctx, &test)
		require.Nil(t, err)
		require.NotNil(t, o)

		record, err := s.ddb.Get(ctx, reflect.TypeOf(test).Name(), test.ID)
		require.Nil(t, err)
		require.Equal(t, test, record.GetStruct())
	})

	t.Run("Should be able to delete record, and return its values", func(t *testing.T) {
		ctx := context.Background()

		test := TestStruct{
			ID:    uuid.NewString(),
			Value: "Some string to delete",
		}

		// Insert the record to be deleted
		o, err := s.ddb.Insert(ctx, &test)
		require.Nil(t, err)
		require.NotNil(t, o)

		// Delete the record
		deletedRecord, err := s.ddb.Delete(ctx, reflect.TypeOf(test).Name(), test.ID)
		require.Nil(t, err)
		require.NotNil(t, deletedRecord)

		// Verify the deleted record matches the original
		require.Equal(t, test, deletedRecord.GetStruct())

		// Attempt to retrieve the deleted record
		_, err = s.ddb.Get(ctx, reflect.TypeOf(test).Name(), test.ID)
		require.NotNil(t, err)

	})

	t.Run("Should be able to batch get items in DB", func(t *testing.T) {
		ctx := context.Background()

		// Create test records
		test1 := TestStruct{
			ID:    uuid.NewString(),
			Value: "Batch item 1",
		}
		test2 := TestStruct{
			ID:    uuid.NewString(),
			Value: "Batch item 2",
		}

		// Insert test records
		_, err := s.ddb.Insert(ctx, &test1)
		require.Nil(t, err)
		_, err = s.ddb.Insert(ctx, &test2)
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
		require.Equal(t, test1, records[0].GetStruct())
		require.Equal(t, test2, records[1].GetStruct())
	})
}
