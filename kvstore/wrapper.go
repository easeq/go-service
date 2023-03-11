package kvstore

import (
	"context"
)

type Wrapper struct {
	s     KVStore
	trace *Trace
}

type PutCallback func(context.Context, *Record, ...SetOpt) (*Record, error)
type GetCallback func(context.Context, string, ...GetOpt) ([]*Record, error)
type DeleteCallback func(context.Context, string, ...DeleteOpt) error
type TxnCallback func(context.Context, TxnHandler) error
type SubscribeCallback func(context.Context, string, SubscribeHandler) error

// NewWrapper returns a new KVStore wrapper
func NewWrapper(s KVStore) *Wrapper {
	return &Wrapper{s, NewTrace(s)}
}

func (w *Wrapper) Put(
	ctx context.Context,
	record *Record,
	put PutCallback,
	opts ...SetOpt,
) (*Record, error) {
	if w.trace == nil {
		return put(ctx, record, opts...)
	}

	return w.trace.Put(ctx, record, put, opts...)
}

func (w *Wrapper) Get(
	ctx context.Context,
	key string,
	get GetCallback,
	opts ...GetOpt,
) ([]*Record, error) {
	if w.trace == nil {
		return get(ctx, key, opts...)
	}

	return w.trace.Get(ctx, key, get, opts...)
}

func (w *Wrapper) Delete(
	ctx context.Context,
	key string,
	delete DeleteCallback,
	opts ...DeleteOpt,
) error {
	if w.trace == nil {
		return delete(ctx, key, opts...)
	}

	return w.trace.Delete(ctx, key, delete, opts...)
}

func (w *Wrapper) Txn(
	ctx context.Context,
	handler TxnHandler,
	txn TxnCallback,
) error {
	if w.trace == nil {
		return txn(ctx, handler)
	}

	return w.trace.Txn(ctx, handler, txn)
}

func (w *Wrapper) Subscribe(
	ctx context.Context,
	key string,
	handler SubscribeHandler,
	subscribe SubscribeCallback,
) error {
	if w.trace == nil {
		return subscribe(ctx, key, handler)
	}

	return w.trace.Subscribe(ctx, key, handler, subscribe)
}

func (w *Wrapper) HandlerHandle(
	ctx context.Context,
	key string,
	handler SubscribeHandler,
	args ...interface{},
) error {
	if w.trace == nil {
		return handler.Handle(ctx, key, args...)
	}

	return w.trace.HandlerHandle(ctx, key, handler, args...)
}
