package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	logrus.SetLevel(logrus.TraceLevel)
	ts := newTestStage(t)

	t.Run("Should be able to store data, and retrieve it", func(t *testing.T) {

		item := Data{
			PK:        "PartitionKey",
			SK:        fmt.Sprintf("SortKey-%d", time.Now().UnixNano()),
			Name:      "Test",
			Age:       25,
			CreatedAt: time.Now(),
		}

		err := ts.store.Create(ctx, item)
		require.Nil(t, err, err)

		retrieved, err := ts.store.Get(ctx, []string{item.PK, item.SK})
		require.Nil(t, err, err)
		require.Equal(t, item.PK, retrieved.PK)
		require.Equal(t, item.SK, retrieved.SK)
		require.Equal(t, item.Name, retrieved.Name)
		require.Equal(t, item.Age, retrieved.Age)
		require.Equal(t, item.CreatedAt.UnixNano(), retrieved.CreatedAt.UnixNano())
	})

	t.Run("Should be able to store data, retrieve it, update it", func(t *testing.T) {
		item := Data{
			PK:        "PartitionKey",
			SK:        fmt.Sprintf("SortKey-%d", time.Now().UnixNano()),
			Name:      "Test",
			Age:       25,
			CreatedAt: time.Now(),
		}

		err := ts.store.Create(ctx, item)
		require.Nil(t, err, err)

		retrieved, err := ts.store.Get(ctx, []string{item.PK, item.SK})
		require.Nil(t, err, err)
		require.Equal(t, item.PK, retrieved.PK)
		require.Equal(t, item.SK, retrieved.SK)
		require.Equal(t, item.Name, retrieved.Name)
		require.Equal(t, item.Age, retrieved.Age)
		require.Equal(t, item.CreatedAt.UnixNano(), retrieved.CreatedAt.UnixNano())

		retrieved.Name = "UpdatedTest"
		retrieved.Age = 30

		updatedItem, err := ts.store.Update(ctx, []string{retrieved.PK, retrieved.SK}, *retrieved)
		require.Nil(t, err, err)

		require.Equal(t, retrieved.PK, updatedItem.PK)
		require.Equal(t, retrieved.SK, updatedItem.SK)
		require.Equal(t, retrieved.Name, updatedItem.Name)
		require.Equal(t, retrieved.Age, updatedItem.Age)
		require.Equal(t, retrieved.CreatedAt.UnixNano(), updatedItem.CreatedAt.UnixNano())
	})

	t.Run("Should be able to create multiple, and search", func(t *testing.T) {
		items := []Data{
			{
				PK:        "PartitionKey",
				SK:        fmt.Sprintf("SortKey-%d", time.Now().UnixNano()),
				Name:      "Test",
				Age:       25,
				CreatedAt: time.Now(),
			},
			{
				PK:        "PartitionKey",
				SK:        fmt.Sprintf("SortKey-%d", time.Now().UnixNano()),
				Name:      "Test",
				Age:       25,
				CreatedAt: time.Now(),
			},
		}

		for _, item := range items {
			err := ts.store.Create(ctx, item)
			require.Nil(t, err, err)
		}

		// searchedItems, err := ts.store.Search(ctx, )
		// require.Nil(t, err, err)
		// require.Len(t, searchedItems, len(items))

		// for i, item := range items {
		// 	require.Equal(t, item.PK, searchedItems[i].PK)
		// 	require.Equal(t, item.SK, searchedItems[i].SK)
		// 	require.Equal(t, item.Name, searchedItems[i].Name)
		// 	require.Equal(t, item.Age, searchedItems[i].Age)
		// 	require.Equal(t, item.CreatedAt.UnixNano(), searchedItems[i].CreatedAt.UnixNano())
		// }
	})

	t.Run("Should be able to delete", func(t *testing.T) {
		item := Data{
			PK:        "PartitionKey",
			SK:        fmt.Sprintf("SortKey-%d", time.Now().UnixNano()),
			Name:      "Test",
			Age:       25,
			CreatedAt: time.Now(),
		}

		err := ts.store.Create(ctx, item)
		require.Nil(t, err, err)

		deleted, err := ts.store.Delete(ctx, []string{item.PK, item.SK})
		require.Nil(t, err, err)

		require.Equal(t, item.PK, deleted.PK)
		require.Equal(t, item.SK, deleted.SK)
		require.Equal(t, item.Name, deleted.Name)
		require.Equal(t, item.Age, deleted.Age)
		require.Equal(t, item.CreatedAt.UnixNano(), deleted.CreatedAt.UnixNano())

		got, err := ts.store.Get(ctx, []string{item.PK, item.SK})
		require.Nil(t, err, err)
		require.Nil(t, got)
	})
}
