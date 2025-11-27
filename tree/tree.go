package tree

import (
	"github.com/mohamadrezamomeni/telecraft/handler"
	"github.com/mohamadrezamomeni/telecraft/pkg/telecrafterror"
)

type Tree struct {
	Handler  handler.HandlerFunc
	path     string
	children map[string]*Tree
}

func New(
	path string,
	handler handler.HandlerFunc,
) *Tree {
	return &Tree{
		Handler:  handler,
		path:     path,
		children: make(map[string]*Tree),
	}
}

func (t *Tree) Set(paths []string, handler handler.HandlerFunc) {
	scope := "tree.set"

	cur := t
	for _, path := range paths {
		if child, ok := cur.children[path]; ok {
			cur = child
		} else {
			cur.children[path] = New(path, nil)
			cur = cur.children[path]
		}
	}
	if cur.Handler != nil {
		panic(
			telecrafterror.
				Scope(scope).
				Errorf("duplicate registration has happened"),
		)
	}
	cur.Handler = handler
}

func (t *Tree) MatchPath(paths []string) (*Tree, map[string]string) {
	return t.matchPathRecursive(paths, 0)
}

func (t *Tree) matchPathRecursive(paths []string, i int) (*Tree, map[string]string) {
	if i >= len(paths) && t.Handler != nil {
		return t, make(map[string]string)
	} else if i >= len(paths) {
		return nil, nil
	}

	curPath := paths[i]

	var res *Tree
	var params map[string]string
	for subPath, child := range t.children {
		if subPath == curPath || (len(subPath) > 0 && subPath[0] == ':') {
			res, params = child.matchPathRecursive(paths, i+1)
		}
		if params != nil && len(subPath) > 0 && subPath[0] == ':' {
			params[subPath[1:]] = curPath
		}
		if res != nil {
			return res, params
		}
	}

	return nil, nil
}
