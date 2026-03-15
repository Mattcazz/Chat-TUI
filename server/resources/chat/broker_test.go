package chat

import (
	"testing"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

func TestBroker_SubscribeAndPublish(t *testing.T) {
	broker := NewBroker()
	ch := broker.Subscribe(1)

	msg := pkg.MsgResponse{UserName: "alice", Content: "hello"}
	broker.Publish(1, msg)

	select {
	case received := <-ch:
		if received.Content != msg.Content || received.UserName != msg.UserName {
			t.Errorf("got %+v, want %+v", received, msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for message")
	}
}

func TestBroker_NoMessageOnWrongConvID(t *testing.T) {
	broker := NewBroker()
	ch := broker.Subscribe(1)

	broker.Publish(2, pkg.MsgResponse{Content: "wrong conv"})

	select {
	case msg := <-ch:
		t.Errorf("expected no message, got %+v", msg)
	case <-time.After(100 * time.Millisecond):
		// correct — nothing received
	}
}

func TestBroker_MultipleSubscribers(t *testing.T) {
	broker := NewBroker()
	ch1 := broker.Subscribe(1)
	ch2 := broker.Subscribe(1)

	msg := pkg.MsgResponse{Content: "broadcast"}
	broker.Publish(1, msg)

	for i, ch := range []chan pkg.MsgResponse{ch1, ch2} {
		select {
		case received := <-ch:
			if received.Content != msg.Content {
				t.Errorf("subscriber %d: got %q, want %q", i, received.Content, msg.Content)
			}
		case <-time.After(time.Second):
			t.Errorf("subscriber %d: timed out", i)
		}
	}
}

func TestBroker_Unsubscribe(t *testing.T) {
	broker := NewBroker()
	ch := broker.Subscribe(1)
	broker.Unsubscribe(1, ch)

	broker.Publish(1, pkg.MsgResponse{Content: "after unsub"})

	select {
	case msg := <-ch:
		t.Errorf("expected no message after unsubscribe, got %+v", msg)
	case <-time.After(100 * time.Millisecond):
		// correct
	}
}

func TestBroker_UnsubscribeOneOfMany(t *testing.T) {
	broker := NewBroker()
	ch1 := broker.Subscribe(1)
	ch2 := broker.Subscribe(1)

	broker.Unsubscribe(1, ch1)

	msg := pkg.MsgResponse{Content: "only ch2"}
	broker.Publish(1, msg)

	select {
	case <-ch1:
		t.Error("ch1 should not receive after unsubscribe")
	case <-time.After(100 * time.Millisecond):
	}

	select {
	case received := <-ch2:
		if received.Content != msg.Content {
			t.Errorf("ch2: got %q, want %q", received.Content, msg.Content)
		}
	case <-time.After(time.Second):
		t.Error("ch2 timed out")
	}
}
