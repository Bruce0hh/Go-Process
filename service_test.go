package gorpc

import (
	"reflect"
	"testing"
)

// 测试服务注册
type Foo int
type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func Test_newService(t *testing.T) {
	var foo Foo
	type args struct {
		rec interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{"test1", args{&foo}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newService(tt.args.rec)
			if !reflect.DeepEqual(len(got.method), tt.want) && reflect.DeepEqual(got.method["Sum"] == nil, tt.want1) {
				t.Errorf("newService() = %v, want %v", len(got.method), tt.want)
			}
		})
	}
}

func Test_service_call(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	mType := s.method["Sum"]
	argv := mType.newArgv()
	reply := mType.newReply()
	argv.Set(reflect.ValueOf(Args{
		Num1: 1,
		Num2: 3,
	}))

	type args struct {
		m     *methodType
		argv  reflect.Value
		reply reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{
			m:     mType,
			argv:  argv,
			reply: reply,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.call(tt.args.m, tt.args.argv, tt.args.reply)
			if (err == nil && *reply.Interface().(*int) == 4 && mType.NumCalls() == 1) != tt.wantErr {
				t.Errorf("call() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
