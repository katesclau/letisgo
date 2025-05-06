package tests

import (
	"context"
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

func (t *TestStruct) Record() any {
	return TestRecord{
		db.Record{
			Pk:      "TestStruct",
			Sk:      t.ID,
			Type:    "TestStruct",
			Version: 0,
		},
		TestProperties{
			Value: t.Value,
		},
	}
}

func Test_DB_Provider(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	t.Run("Should be able insert item in DB", func(t *testing.T) {
		ctx := context.Background()
		s := newTestStage(t)

		o, err := s.ddb.Insert(ctx, &TestStruct{
			ID:    uuid.NewString(),
			Value: "Some string",
		})
		logrus.Info(err)
		require.Nil(t, err)

		require.NotNil(t, o)
	})
}
