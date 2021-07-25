package myutils

type Rectangle struct {
	X uint32
	Y uint32
	W uint32
	H uint32
}

func (s *Rectangle) SetRect(x, y, w, h uint32) {
	s.X = x
	s.Y = y
	s.W = w
	s.H = h
}

func (s *Rectangle) SetRect2(x, y, x2, y2 uint32) {
	s.X = MinUint32(x, x2)
	s.Y = MinUint32(y, y2)
	s.W = uint32(AbsInt(int32(x - x2)))
	s.H = uint32(AbsInt(int32(y - y2)))
}

func (s *Rectangle) IsInRect(x, y uint32) bool {
	w, h := s.W, s.H
	if w == 0 {
	}
	if h == 0 {
	}
	return x >= s.X && x <= s.X+w && y >= s.Y && y <= s.Y+h
}

type RectangleList struct {
	list []*Rectangle
}

func (s *RectangleList) IsInRect(x, y uint32) bool {
	for _, v := range s.list {
		if v.IsInRect(x, y) {
			return true
		}
	}
	return false
}

func (s *RectangleList) AddRect(x, y, w, h uint32) {
	if s.list == nil {
		s.list = make([]*Rectangle, 0)
	}
	rect := &Rectangle{}
	rect.SetRect(x, y, w, h)
	s.list = append(s.list, rect)
}

func (s *RectangleList) AddRect2(x, y, x2, y2 uint32) {
	if s.list == nil {
		s.list = make([]*Rectangle, 0)
	}
	rect := &Rectangle{}
	rect.SetRect2(x, y, x2, y2)
	s.list = append(s.list, rect)
}

func (s *RectangleList) Clear() {
	s.list = make([]*Rectangle, 0)
}

func (s *RectangleList) Size() int {
	return len(s.list)
}

func GetAroundPoint(x, y uint32, index uint32) (dx, dy uint32) {
	if index == 0 {
		return x, y
	}
	var i uint32
	for i = 1; i < 100; i += 2 {
		if i*i > index {
			break
		}
	}
	n2 := index - (i-2)*(i-2) //
	if n2 < i {
		dx = x - (i-1)/2
		dy = y + n2 - (i-1)/2
		return dx, dy
	}
	if n2 >= i+(i-2)*2 {
		dx = x + (i-1)/2
		dy = y + (n2 - i - (i-2)*2) - (i-1)/2
		return dx, dy
	}
	if n2%2 == 0 {
		dx = x + (n2-i)/2 - (i-3)/2
		dy = y - (i-1)/2
	} else {
		dx = x + (n2-i)/2 - (i-3)/2
		dy = y + (i-1)/2
	}
	return dx, dy
}
