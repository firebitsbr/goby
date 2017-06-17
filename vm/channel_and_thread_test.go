package vm

import "testing"

func TestObjectMutationInThread(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		c = Channel.new

		i = 0
		thread do
		  i++
		  c.deliver(i)
		end

		# Used to block main process until thread is finished
		c.receive
		i
		`, 1},
		{`
		c = Channel.new

		i = 0
		thread do
		  i++
		  c.deliver(i)
		end

		i++
		# Used to block main process until thread is finished
		c.receive
		i
		`, 2},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}

func TestObjectDeliveryBetweenThread(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		c = Channel.new

		thread do
		  s = "123"
		  c.deliver(s)
		end

		c.receive
		`, "123"},
		{`
		c = Channel.new

		thread do
		  h = "Hello"
		  w = "World"
		  c.deliver(h)
		  c.deliver(w)
		end

		h = c.receive
		w = c.receive

		h + " " + w
		`, "Hello World"},
		{`
		class Foo
		  def bar
		    100
		  end
		end

		c = Channel.new

		thread do
		  f = Foo.new
		  c.deliver(f)
		end

		c.receive.bar
		`, 100},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}
