package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"go.uber.org/zap"
)

func TestWriterFailWithKafkaSyncer(t *testing.T) {
	config := sarama.NewConfig()
	p := mocks.NewAsyncProducer(t, config)

	buf := make([]byte, 0, 256)
	w := bytes.NewBuffer(buf)
	w.Write([]byte("hello"))

	logger := New(NewKafakSyncer(p, "test", NewFileSyncer(w)), 0)

	p.ExpectInputAndFail(errors.New("producer error"))
	p.ExpectInputAndFail(errors.New("producer error"))

	// all below will be written to the fallback sycner
	logger.Info("demo1", zap.String("status", "ok")) // write to the kafka syncer
	logger.Info("demo2", zap.String("status", "ok")) // write to the kafka syncer

	// make sure the goroutine which handles the error writes the log to the fallback syncer
	time.Sleep(2 * time.Second)

	s := string(w.Bytes())
	if !strings.Contains(s, "demo1") {
		t.Errorf("want true, actual false")
	}
	if !strings.Contains(s, "demo2") {
		t.Errorf("want true, actual false")
	}

	if err := p.Close(); err != nil {
		t.Error(err)
	}
}
