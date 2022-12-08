package gopool

type Option func(*Pool)

func WithBlock(block bool) Option {
	return func(p *Pool) {
		p.block = block
	}
}

func WithPreAlloc(preAlloc bool) Option {
	return func(p *Pool) {
		p.preAlloc = preAlloc
	}
}
