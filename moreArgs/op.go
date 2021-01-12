/**
 * @Author: zhangsan
 * @Description:
 * @File:  kv
 * @Version: 1.0.0
 * @Date: 2021/1/11 下午5:21
 */

package moreArgs

type Op struct {
	limit int64
	order bool
	group bool
	key string
	value string
}

type Options func (op *Op)

func WithLimit(limit int64) Options{
	return func(op *Op) {op.limit= limit}
}

func WithOrder(or bool) Options{
	return func(op *Op) {op.order= or}
}
func WithGroup(or bool) Options{
	return func(op *Op) {op.group= or}
}

func OpPut(opts... Options) Op{
	ret := Op{key: "string key",value: "string value"}
	ret.applyOpts(opts)
	return ret

}
func (op *Op)applyOpts(opts []Options){
	for _,opt := range opts{
		opt(op)
	}
}