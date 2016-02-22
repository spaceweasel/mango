package mango

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

func newTree() *tree {
	t := tree{}
	t.paramValidator = newParameterValidators()
	return &t
}

type tree struct {
	root             *treenode
	paramValidator   RouteParamValidators
	GlobalCORSConfig *CORSConfig
}

// AddRouteParamValidator adds a new validator to the collection.
// AddRouteParamValidator panics if a validator with the same Type()
// exists.
func (t *tree) AddRouteParamValidator(v ParamValidator) {
	t.paramValidator.AddValidator(v)
}

// AddRouteParamValidators adds a slice of new validators to the collection.
// AddRouteParamValidators panics if a validator with the same Type()
// exists.
func (t *tree) AddRouteParamValidators(validators []ParamValidator) {
	t.paramValidator.AddValidators(validators)
}

// SetRouteParamValidators sets the internal paramValidator collection.
func (t *tree) SetRouteParamValidators(v RouteParamValidators) {
	t.paramValidator = v
}

// Root returns the tree root node, assigning a new empty node if
// one has not been set already.
func (t *tree) Root() *treenode {
	if t.root == nil {
		t.root = &treenode{}
	}
	return t.root
}

// SetGlobalCORS sets the CORS configuration that will be used for
// a resource if it has no CORS configuration of its own. If the
// resource has no CORSConfig and tree.GlobalCORSConfig is nil
// then CORS request are treated like any other.
func (t *tree) SetGlobalCORS(config CORSConfig) {
	t.GlobalCORSConfig = &config
}

// SetCORS sets the CORS configuration that will be used for
// the resource matching the pattern.
// These settings override any global settings.
func (t *tree) SetCORS(pattern string, config CORSConfig) {
	node, _ := t.Root().addNode(pattern)
	node.CORSConfig = &config
}

// AddCORS sets the CORS configuration that will be used for
// the resource matching the pattern, by merging the supplied
// config with any globalCORSConfig.
// SetGlobalCORS MUST be called before this method!
func (t *tree) AddCORS(pattern string, config CORSConfig) {
	node, _ := t.Root().addNode(pattern)
	if t.GlobalCORSConfig == nil {
		node.CORSConfig = &config
		return
	}
	c := t.GlobalCORSConfig.clone()
	c.merge(config)
	node.CORSConfig = c
}

// AddHandlerFunc adds a new handlerFunc for the supplied pattern and method.
// If a handlerFunc already exists for the pattern-method combination,
// AddHandlerFunc panics.
func (t *tree) AddHandlerFunc(pattern, method string, handlerFunc ContextHandlerFunc) {
	node, pNames := t.Root().addNode(pattern)

	if node.handlers == nil {
		node.handlers = make(map[string]ContextHandlerFunc)
		node.paramNames = pNames
	} else if _, e := node.handlers[method]; e {
		panic(fmt.Sprintf("duplicate route handler method: \"%s %s\"", method, pattern))
	}
	node.handlers[method] = handlerFunc
}

// Resource is a container holding the Handler functions for
// the various HTTP methods, a RouteParams map of values obtained
// from the request path and a CORS configuration.
// The CORS config may be nil.
type Resource struct {
	Handlers    map[string]ContextHandlerFunc
	RouteParams map[string]string
	CORSConfig  *CORSConfig
}

// GetResource traverses the tree looking for a leaf nodes which match the supplied path.
// If found, GetResource returns the resource held at the leaf node.
// If the leaf node journey involves parameter nodes, then associated values
// will be extracted from the path and added to the resource RouteParams map.
func (t *tree) GetResource(path string) (*Resource, bool) {
	n, pValues, ok := t.search(t.Root().children, path)
	if !ok {
		return nil, false
	}
	res := Resource{}
	res.RouteParams = make(map[string]string)
	if n.paramNames != nil {
		for i, n := range n.paramNames.items {
			res.RouteParams[n] = pValues.items[i]
		}
	}
	res.Handlers = n.handlers
	res.CORSConfig = n.CORSConfig
	if res.CORSConfig == nil {
		res.CORSConfig = t.GlobalCORSConfig
	}
	return &res, true
}

func (t *tree) search(nodes []*treenode, path string) (*treenode, *stringList, bool) {
	i := strings.IndexByte(path, byte('/'))
	for _, node := range nodes {
		if node.isParam {
			if i >= 0 {
				value := path[:i]
				if !t.paramValidator.IsValid(value, node.paramConstraint) {
					continue
				}
				path = path[i:]
				t, paramValues, s := t.search(node.children, path)
				paramValues = addItem(paramValues, value)
				return t, paramValues, s
			}
			if !t.paramValidator.IsValid(path, node.paramConstraint) {
				continue
			}
			paramValues := newStringList(path)
			return node, paramValues, true
		}

		n := node
		for j := 0; j < len(node.label); j++ {
			if len(path) <= j || path[j] != node.label[j] {
				n = nil
				break
			}
		}
		if n != nil {
			path = path[len(n.label):]
			if len(path) == 0 {
				return n, nil, true
			}
			return t.search(n.children, path)
		}
	}
	return nil, nil, false
}

func (t *tree) Print() {
	for _, n := range t.Root().children {
		n.Print(0)
	}
}

func (t *tree) GetStats() treeStats {
	stats := treeStats{totalNodes: t.Root().Count()}
	return stats
}

type treeStats struct {
	totalNodes int
}

type treenode struct {
	children        []*treenode
	label           string
	handlers        map[string]ContextHandlerFunc
	paramNames      *stringList
	isParam         bool
	paramConstraint string
	CORSConfig      *CORSConfig
}

func (n *treenode) insert(child *treenode) {
	n.append(child)
	for i := len(n.children) - 1; i > 0; i-- {
		n.children[i] = n.children[i-1]
	}
	n.children[0] = child
}

func (n *treenode) append(child *treenode) {
	n.children = append(n.children, child)
}

func (n *treenode) addParamNode(pattern string) (*treenode, *stringList, bool) {
	i := strings.IndexByte(pattern, byte('{'))
	if i >= 0 {
		node, paramNames := n.addNode(pattern[:i])
		pattern = pattern[i+1:]
		i := strings.IndexByte(pattern, byte('}'))
		if i < 0 {
			panic("invalid route syntax: {" + pattern)
		}

		nc := strings.Split(pattern[:i], ":") // split {name:constraint}
		name := strings.TrimSpace(nc[0])
		constraint := ""
		if len(nc) > 1 {
			constraint = strings.TrimSpace(strings.Join(nc[1:], ""))
		}
		pn := node.addParamChild(constraint)
		paramNames = addItem(paramNames, name)
		pattern = pattern[i+1:]
		if len(pattern) == 0 {
			return pn, paramNames, true
		}

		node, paramNames = pn.addNode(pattern)
		paramNames = addItem(paramNames, name)
		return node, paramNames, true
	}
	return nil, nil, false
}

func (n *treenode) addNode(pattern string) (*treenode, *stringList) {
	// handle any parameters first...
	node, params, ok := n.addParamNode(pattern)
	if ok {
		return node, params
	}

	for _, child := range n.children {
		if child.isParam {
			continue
		}

		j, r := 0, comlen(child.label, pattern)
		for ; j < r; j++ {
			if child.label[j] != pattern[j] {
				break
			}
		}

		switch {
		case j == 0:
			continue
		case len(child.label) == j:
			// pattern "starts with label"
			if len(pattern) == j {
				return child, nil
			}
			// continue with remainder of pattern...
			return child.addNode(pattern[j:])
		case len(child.label) > j:
			// both start with common string - need to split the current child node.
			// create new grandchild node using "uncommon" part of Label
			gc := &treenode{label: child.label[j:]}
			// move associated data from child to new grandchild, and set the current
			// child's children to a new slice containing only the new grandchiild node
			gc.children, child.children = child.children, []*treenode{gc}
			gc.handlers, child.handlers = child.handlers, nil
			//n.ParamNames, node.ParamNames = node.ParamNames, nil
			// reset current node Label to "common" part...
			child.label = pattern[:j]
			pattern = pattern[j:]
			if len(pattern) == 0 {
				return child, nil
			}
			// continue inserting with "uncommon" part...
			return child.addNode(pattern)
		}
	}

	// no matches - just add a new node with full pattern
	node = &treenode{label: pattern}
	// insert as first child to be ahead of any parametised siblings
	n.insert(node)
	return node, nil
}

func (n *treenode) addParamChild(constraint string) *treenode {
	pc := n.paramChild(constraint)
	if pc == nil {
		pc = &treenode{paramConstraint: constraint, isParam: true}
		// append as last child to ensure after any non-parametised siblings
		n.append(pc)
		// ensure that any empty constraint is always the last sibling by
		// checking penultimate sibling and moving to end if an empty constraint
		c := n.children
		l := len(c)
		if l > 1 && c[l-2].isParam && c[l-2].paramConstraint == "" {
			c[l-2], c[l-1] = c[l-1], c[l-2]
		}
	}
	return pc
}

func (n *treenode) paramChild(constraint string) *treenode {
	for _, c := range n.children {
		if c.isParam && c.paramConstraint == constraint {
			return c
		}
	}
	return nil
}

func (n *treenode) Print(l int) {
	tab := strings.Repeat("\t", l)
	handlers := ""
	pns := ""
	if n.handlers != nil {
		for k, h := range n.handlers {
			name := extractFnName(h)
			handlers += fmt.Sprintf("[%s: %s]", k, name)
		}
		handlers = "Handlers " + handlers

		pns = " ParamNames ["
		if n.paramNames != nil {
			for _, pn := range n.paramNames.items {
				pns += pn + ","
			}
		}
		pns += "]"
	}

	if n.isParam {
		log.Printf("%s>Label: %q\t%s\t%s", tab, n.label, handlers, pns)
	} else {
		log.Printf("%s>Param: %q\t%s\t%s", tab, n.paramConstraint, handlers, pns)
	}
	l++
	for _, n := range n.children {
		n.Print(l)
	}
}

func (n *treenode) Count() int {
	count := 0
	for _, child := range n.children {
		count += child.Count()
		count++
	}
	return count
}

type stringList struct {
	items []string
}

func (n *stringList) append(s string) {
	n.items = append(n.items, s)
}

func newStringList(s string) *stringList {
	return &stringList{items: []string{s}}
}

func addItem(n *stringList, s string) *stringList {
	if n == nil {
		return newStringList(s)
	}
	n.append(s)
	return n
}

func comlen(a, b string) int {
	al, bl := len(a), len(b)
	if al > bl {
		return bl
	}
	return al
}

func extractFnName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return name
}
