package socketStream

import (
	"sync"
)

type PackagePacker struct {
	b []byte
	l sync.Mutex
}

func NewPackagePacker() *PackagePacker {
	return &PackagePacker{b: make([]byte, 0)}
}

func (s *PackagePacker) Clear() {
	s.b = s.b[0:0]
}

func (s *PackagePacker) Append(b []byte) {
	s.l.Lock()
	defer s.l.Unlock()
	s.b = append(s.b, b...)
}

func (s *PackagePacker) GetPackage() *TPackage {
	s.l.Lock()
	defer s.l.Unlock()
	length := len(s.b)
	if length > 0 {
		pkg := NewTPackage()
		err := pkg.ParseBuf(s.b)
		if err != nil {
			return nil
		}
		pkgLen := int(pkg.BodySize + 10)
		if pkgLen == length {
			s.Clear()
		} else {
			s.b = s.b[pkgLen:]
		}
		return pkg
	}
	return nil
}
