package data

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Path struct {
	segments []string
}

func (path Path) ToString() string {
	return strings.Join(path.segments,"/")
}

func (path *Path) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&path.segments)
}

func NewPath(segs ...string) Path {
	return Path{
		segments: segs,
	}
}

func (path Path) Extend(segment string) Path {
	newSegments := make([]string, len(path.segments)+1)
	if path.segments!=nil {
		copy(newSegments, path.segments)
	}
	newSegments[len(newSegments)-1] = segment
	return Path{
		segments: newSegments,
	}
}

func (this Path) Match(that Path) bool {
	if len(this.segments) != len(that.segments) {
		return false;
	}
	for i := 0; i < len(this.segments); i++ {
		r, err := regexp.Compile("^" + this.segments[i] + "$")
		if err != nil {
			panic(fmt.Errorf("Path segment is not a regexp %s in %s", this.segments[i], this.segments))
		}
		if !r.Match([]byte(that.segments[i])) {
			return false
		}
	}
	return true
}

func (this Path) MatchSegments(segments ...string) bool {
	return this.Match(NewPath(segments...))
}

func (path Path) Length() int {
	return len(path.segments)
}
func (path Path) Segment(i int) string {
	if i >= 0 {
		return path.segments[i]
	} else if len(path.segments)+i >= 0 {
		return path.segments[len(path.segments)+i]
	} else {
		return ""
	}
}
func (path Path) IsEmpty() bool {
	return len(path.segments) == 0
}

func (path Path) Equal(other Path) bool {
	return reflect.DeepEqual(path.segments, other.segments)
}
func (path Path) Parent() Path {
	if len(path.segments) > 1 {
		return Path{
			segments: path.segments[0 : len(path.segments)-1],
		}
	} else {
		return NewPath()
	}
}
func (path Path) Last() string {
	return path.segments[len(path.segments)-1]

}
