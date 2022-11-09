// Package pool include helpers for work witch sync.Pool
package pool

import "sync"

// Pool generic for create any poll by T
//
// For example:
//
//	 type ShortedCreateDTOPool struct {
//	     pool.Pool[ShortedCreateDTO]
//	 }
//
//	 func (p *ShortedCreateDTOPool) Put(v *ShortedCreateDTO) {
//			v.URL = ""
//			p.Pool.Put(v)
//	 }
type Pool[T any] sync.Pool

// Get return struct from pool
func (p *Pool[T]) Get() *T {
	v := (*sync.Pool)(p).Get()
	if v == nil {
		var zero T
		v = &zero
	}
	return v.(*T)
}

// Put Return to pool
func (p *Pool[T]) Put(t *T) {
	(*sync.Pool)(p).Put(t)
}
