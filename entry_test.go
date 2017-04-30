package tui

import (
	"image"
	"testing"
)

var entrySizeTests = []struct {
	test        string
	setup       func() *Entry
	minSizeHint image.Point
	sizeHint    image.Point
	size        image.Point
}{
	{
		test: "Entry default size",
		setup: func() *Entry {
			return NewEntry()
		},
		minSizeHint: image.Point{1, 1},
		sizeHint:    image.Point{10, 1},
		size:        image.Point{100, 100},
	},
}

func TestEntry_Size(t *testing.T) {
	for _, tt := range entrySizeTests {
		e := tt.setup()
		e.Resize(image.Point{100, 100})

		if got := e.Size(); got != tt.size {
			t.Errorf("e.Size() = %s; want = %s", got, tt.size)
		}
		if got := e.SizeHint(); got != tt.sizeHint {
			t.Errorf("e.SizeHint() = %s; want = %s", got, tt.sizeHint)
		}
		if got := e.MinSizeHint(); got != tt.minSizeHint {
			t.Errorf("e.MinSizeHint() = %s; want = %s", got, tt.minSizeHint)
		}
	}
}

var drawEntryTests = []struct {
	test  string
	size  image.Point
	setup func() *Entry
	want  string
}{
	{
		test: "Empty entry",
		size: image.Point{15, 5},
		setup: func() *Entry {
			return NewEntry()
		},
		want: `
               
...............
...............
...............
...............
`,
	},
	{
		test: "Entry with text",
		size: image.Point{15, 5},
		setup: func() *Entry {
			e := NewEntry()
			e.SetText("test")
			return e
		},
		want: `
test           
...............
...............
...............
...............
`,
	},
	{
		test: "Scrolling entry",
		size: image.Point{15, 5},
		setup: func() *Entry {
			e := NewEntry()
			e.SetText("Lorem ipsum dolor sit amet")
			return e
		},
		want: `
 dolor sit amet
...............
...............
...............
...............
`,
	},
	{
		test: "Scrolling entry when focused",
		size: image.Point{15, 5},
		setup: func() *Entry {
			e := NewEntry()
			e.SetText("Lorem ipsum dolor sit amet")
			e.SetFocused(true)
			return e
		},
		want: `
dolor sit amet 
...............
...............
...............
...............
`,
	},
}

func TestEntry_Draw(t *testing.T) {
	for _, tt := range drawEntryTests {
		tt := tt
		t.Run(tt.test, func(t *testing.T) {
			var surface *testSurface
			if tt.size.X == 0 && tt.size.Y == 0 {
				surface = newTestSurface(10, 5)
			} else {
				surface = newTestSurface(tt.size.X, tt.size.Y)
			}
			painter := NewPainter(surface, NewTheme())

			b := tt.setup()

			b.Resize(surface.size)
			b.Draw(painter)

			if surface.String() != tt.want {
				t.Errorf("got = \n%s\n\nwant = \n%s", surface.String(), tt.want)
			}
		})
	}
}

func TestEntry_OnChanged(t *testing.T) {
	e := NewEntry()

	var invoked bool
	e.OnChanged(func(e *Entry) {
		invoked = true
		if e.Text() != "t" {
			t.Errorf("e.Text() = %s; want = %s", e.Text(), "t")
		}
	})

	ev := Event{
		Type: EventKey,
		Ch:   't',
	}

	t.Run("When entry is not focused", func(t *testing.T) {
		e.OnEvent(ev)
		if invoked {
			t.Errorf("entry should not be submitted")
		}
	})

	invoked = false
	e.SetFocused(true)

	t.Run("When entry is focused", func(t *testing.T) {
		e.OnEvent(ev)
		if !invoked {
			t.Errorf("entry should be submitted")
		}
	})
}

func TestEntry_OnSubmit(t *testing.T) {
	e := NewEntry()

	var invoked bool
	e.OnSubmit(func(e *Entry) {
		invoked = true
	})

	ev := Event{
		Type: EventKey,
		Key:  KeyEnter,
	}

	t.Run("When entry is not focused", func(t *testing.T) {
		e.OnEvent(ev)
		if invoked {
			t.Errorf("entry should not be submitted")
		}
	})

	invoked = false
	e.SetFocused(true)

	t.Run("When entry is focused", func(t *testing.T) {
		e.OnEvent(ev)
		if !invoked {
			t.Errorf("button should be submitted")
		}
	})
}

var layoutEntryTests = []struct {
	test  string
	setup func() *Box
	want  string
}{
	{
		test: "Preferred",
		setup: func() *Box {
			e := NewEntry()
			e.SetSizePolicy(Preferred, Preferred)

			b := NewHBox(e)
			b.SetBorder(true)
			b.SetSizePolicy(Expanding, Expanding)

			return b
		},
		want: `
┌──────────────────┐
│                  │
│..................│
│..................│
└──────────────────┘
`,
	},
	{
		test: "Preferred/Preferred",
		setup: func() *Box {
			e1 := NewEntry()
			e1.SetSizePolicy(Preferred, Preferred)
			e1.SetText("0123456789foo")

			e2 := NewEntry()
			e2.SetSizePolicy(Preferred, Preferred)
			e2.SetText("0123456789bar")

			b := NewHBox(e1, e2)
			b.SetBorder(true)
			b.SetSizePolicy(Expanding, Expanding)

			return b
		},
		want: `
┌──────────────────┐
│456789foo456789bar│
│..................│
│..................│
└──────────────────┘
`,
	},
	{
		test: "Preferred/Minimum",
		setup: func() *Box {
			e1 := NewEntry()
			e1.SetSizePolicy(Preferred, Preferred)
			e1.SetText("0123456789foo")

			e2 := NewEntry()
			e2.SetSizePolicy(Minimum, Preferred)
			e2.SetText("0123456789bar")

			b := NewHBox(e1, e2)
			b.SetBorder(true)
			b.SetSizePolicy(Expanding, Expanding)

			return b
		},
		want: `
┌──────────────────┐
│0123456789foo89bar│
│..................│
│..................│
└──────────────────┘
`,
	},
	{
		test: "Minimum/Preferred",
		setup: func() *Box {
			e1 := NewEntry()
			e1.SetSizePolicy(Minimum, Preferred)
			e1.SetText("0123456789foo")

			e2 := NewEntry()
			e2.SetSizePolicy(Preferred, Preferred)
			e2.SetText("0123456789bar")

			b := NewHBox(e1, e2)
			b.SetBorder(true)
			b.SetSizePolicy(Expanding, Expanding)

			return b
		},
		want: `
┌──────────────────┐
│89foo0123456789bar│
│..................│
│..................│
└──────────────────┘
`,
	},
	{
		test: "Preferred/Expanding",
		setup: func() *Box {
			e1 := NewEntry()
			e1.SetSizePolicy(Preferred, Preferred)
			e1.SetText("foo")

			e2 := NewEntry()
			e2.SetSizePolicy(Expanding, Preferred)
			e2.SetText("bar")

			b := NewHBox(e1, e2)
			b.SetBorder(true)
			b.SetSizePolicy(Expanding, Expanding)

			return b
		},
		want: `
┌──────────────────┐
│foobar            │
│..................│
│..................│
└──────────────────┘
`,
	},
	{
		test: "Expanding/Preferred",
		setup: func() *Box {
			e1 := NewEntry()
			e1.SetSizePolicy(Expanding, Preferred)
			e1.SetText("foo")

			e2 := NewEntry()
			e2.SetSizePolicy(Preferred, Preferred)
			e2.SetText("bar")

			b := NewHBox(e1, e2)
			b.SetBorder(true)
			b.SetSizePolicy(Expanding, Expanding)

			return b
		},
		want: `
┌──────────────────┐
│foo            bar│
│..................│
│..................│
└──────────────────┘
`,
	},
}

func TestEntry_Layout(t *testing.T) {
	t.Skip("enable when layout engine is working correctly")

	for _, tt := range layoutEntryTests {
		tt := tt
		t.Run(tt.test, func(t *testing.T) {
			surface := newTestSurface(20, 5)
			painter := NewPainter(surface, NewTheme())

			b := tt.setup()

			b.Resize(surface.size)
			b.Draw(painter)

			if surface.String() != tt.want {
				t.Errorf("got = \n%s\n\nwant = \n%s", surface.String(), tt.want)
			}
		})
	}
}
