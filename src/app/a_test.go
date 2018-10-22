package main

import (
	"errors"
	"testing"
	"time"

	xlog "github.com/qiniu/x/xlog.v7"
)

func do1() (err error) {
	var (
		xl = xlog.New("main")
		// err error
	)
	ch := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 2)
		err = errors.New("hello")
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		xl.Infof("%v", err)
		return err
	case <-time.After(time.Second * 2):
		err = errors.New("request timeout")
		return err
	}

}

func TestReadRace(t *testing.T) {
	do1()
	time.Sleep(time.Second * 3)
}
