package oracle

import (
	"context"
	"sync"

	oracle "github.com/godoes/gorm-oracle"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Oracle struct {
	DB  *gorm.DB
	DSN string
}

func New(dsn string) *Oracle {
	return &Oracle{&gorm.DB{}, dsn}
}

func (o *Oracle) Open() error {
	var err error
	o.DB, err = gorm.Open(oracle.Open(o.DSN), &gorm.Config{})
	return err
}

func OpenArr(oracles []*Oracle) ([]*Oracle, error) {
	var (
		eg  errgroup.Group
		mu  sync.Mutex
		res []*Oracle
	)
	for _, v := range oracles {
		oracle := v
		eg.Go(func() error {
			err := oracle.Open()
			if err != nil {
				return err
			}
			mu.Lock()
			res = append(res, oracle)
			mu.Unlock()
			return nil
		})
	}
	return res, eg.Wait()
}

func QueryArr[T any](c context.Context, oracles []*Oracle, q string) ([]T, error) {
	var (
		eg  errgroup.Group
		mu  sync.Mutex
		res []T
	)
	for _, v := range oracles {
		oracle := v
		eg.Go(func() error {
			var t []T
			err := oracle.DB.WithContext(c).Raw(q).Scan(&t).Error
			if err != nil {
				return err
			}
			mu.Lock()
			res = append(res, t...)
			mu.Unlock()
			return nil
		})

	}
	return res, eg.Wait()
}
