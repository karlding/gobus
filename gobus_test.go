package gobus

import (
	"sync/atomic"
	"testing"
)

func TestNew(t *testing.T) {
	goBus := New()

	if goBus == nil {
		t.Errorf("could not create GoBus")
		t.Fail()
	}
}

func TestSubscribe(t *testing.T) {
	goBus := New()

	// should successfully create a new Subscriber on the Message Bus
	if goBus.Subscribe("test", func() {}) != nil {
		t.Fail()
	}

	if goBus.Subscribe("test", func() {}) != nil {
		t.Fail()
	}
}

func TestSubscribeNonFunction(t *testing.T) {
	goBus := New()

	// should fail since "String" is not a callback
	if goBus.Subscribe("test", "String") == nil {
		t.Fail()
	}
}

func TestSubscribeMultiple(t *testing.T) {
	goBus := New()

	callbackOne := false
	callbackTwo := false

	goBus.Subscribe("test", func() {
		callbackOne = true
	})

	goBus.Subscribe("test", func() {
		callbackTwo = true
	})

	goBus.Publish("test")

	goBus.Wait()

	if callbackOne == false {
		t.Fail()
	}
	if callbackTwo == false {
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	goBus := New()

	goBus.Subscribe("test", func(a int, b int) {
		if a != b {
			t.Fail()
		}
	})
	goBus.Publish("test", 10, 10)

	goBus.Wait()
}

func TestPublishIgnore(t *testing.T) {
	goBus := New()

	// the callback should not run since it is not subscribed
	goBus.Subscribe("test", func(a string) {
		t.Fail()
	})

	goBus.Publish("solar", "car")
}

func TestPublishNoSubscriber(t *testing.T) {
	goBus := New()

	goBus.Publish("test", 10, 10)
}

func TestPublishWaitOrder(t *testing.T) {
	goBus := New()
	ops := uint64(0)
	sum := uint64(0)

	goBus.Subscribe("sums", func(a uint64) {
		atomic.AddUint64(&ops, a)
	})

	for i := uint64(1); i <= 100000; i++ {
		goBus.Publish("sums", i)

		goBus.Wait()
		sum = sum + i
		if ops != sum {
			t.Errorf("Expected %d but got %d\n", sum, ops)
			t.Fail()
		}
	}
}

func TestPublishWait(t *testing.T) {
	goBus := New()
	ops := uint64(0)
	calls := uint64(100000)

	goBus.Subscribe("sums", func(a uint64) {
		atomic.AddUint64(&ops, 1)
	})

	for i := uint64(1); i <= calls; i++ {
		goBus.Publish("sums", i)
	}

	goBus.Wait()

	if ops != calls {
		t.Errorf("Expected %d but got %d\n", calls, ops)
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	goBus := New()

	handler := func(a int) {
		if a != 1 {
			t.Errorf("Expected %d but got %d\n", 1, a)
			t.Fail()
		}
	}

	goBus.Subscribe("test", handler)
	goBus.Publish("test", 1)

	goBus.Wait()
	goBus.Unsubscribe("test", handler)
}

func TestUnsubscribeFail(t *testing.T) {
	goBus := New()
	handler := func() {}

	goBus.Subscribe("topic", handler)
	if goBus.Unsubscribe("topic", handler) != nil {
		t.Fail()
	}
	if goBus.Unsubscribe("topic", handler) == nil {
		t.Fail()
	}
}

func TestUnsubscribeFailTopic(t *testing.T) {
	goBus := New()
	handler := func() {}

	goBus.Subscribe("topic", handler)

	if goBus.Unsubscribe("missing", "yolo") == nil {
		t.Fail()
	}
}
