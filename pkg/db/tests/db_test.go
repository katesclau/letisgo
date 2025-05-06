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

}
