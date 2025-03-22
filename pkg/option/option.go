package option

type Option[T any] interface {
	apply(t *T) (Settable[T], error)
}

func Apply[T any](t *T, must bool, opts ...Option[T]) {
	for _, opt := range opts {
		if _, err := opt.apply(t); err != nil && must {
			panic(err)
		}
	}
}

func ApplyError[T any](t *T, opts ...Option[T]) error {
	for _, opt := range opts {
		if _, err := opt.apply(t); err != nil {
			return err
		}
	}
	return nil
}

func ApplyRevert[T any](t *T, must bool, opts ...Option[T]) (revertOpts []Settable[T]) {
	for _, opt := range opts {
		if revert, err := opt.apply(t); err != nil && must {
			panic(err)
		} else if revert != nil {
			revertOpts = append(revertOpts, revert)
		}
	}
	return revertOpts
}

func ApplyRevertError[T any](t *T, opts ...Option[T]) (revertOpts []Settable[T], err error) {
	for _, opt := range opts {
		if revert, err := opt.apply(t); err != nil {
			return nil, err
		} else if revert != nil {
			revertOpts = append(revertOpts, revert)
		}
	}
	return revertOpts, nil
}

type Settable[T any] func(*T)

var _ Option[uint8] = Settable[uint8](nil)

//lint:ignore U1000 used by Apply* functions operating on the Option[T] interface
func (o Settable[T]) apply(t *T) (Settable[T], error) {
	o(t)
	return nil, nil
}

type SetErrorable[T any] func(*T) error

var _ Option[uint8] = SetErrorable[uint8](nil)

//lint:ignore U1000 used by Apply* functions operating on the Option[T] interface
func (o SetErrorable[T]) apply(t *T) (Settable[T], error) {
	return nil, o(t)
}

type Revertable[T any] func(*T) func(*T)

var _ Option[uint8] = Revertable[uint8](nil)

//lint:ignore U1000 used by Apply* functions operating on the Option[T] interface
func (o Revertable[T]) apply(t *T) (Settable[T], error) {
	return o(t), nil
}

type RevertErrorable[T any] func(*T) (Settable[T], error)

var _ Option[uint8] = RevertErrorable[uint8](nil)

//lint:ignore U1000 used by Apply* functions operating on the Option[T] interface
func (o RevertErrorable[T]) apply(t *T) (Settable[T], error) {
	return o(t)
}
