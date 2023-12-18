// Package orderednodes is a way to create a hierarchical tree of nodes,
// where the nodes are ordered and keep their order, and without needing
// or using Go generics.
//
// Node order is of course important for markup in general and XML mixed
// content in particular, while unimportant when XML is used purely for
// data records.
//
// interface [Norder] is implemented not for type [Nord] but rather for
// `*Nord´. NThis is so that nodes are writable. 
//
package orderednodes
