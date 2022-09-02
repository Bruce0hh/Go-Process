package gorpc

import (
	"reflect"
	"testing"
)

func Test_isExportedOrBuiltinType(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExportedOrBuiltinType(tt.args.t); got != tt.want {
				t.Errorf("isExportedOrBuiltinType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_methodType_NumCalls(t *testing.T) {
	type fields struct {
		method    reflect.Method
		ArgType   reflect.Type
		ReplyType reflect.Type
		numCalls  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &methodType{
				method:    tt.fields.method,
				ArgType:   tt.fields.ArgType,
				ReplyType: tt.fields.ReplyType,
				numCalls:  tt.fields.numCalls,
			}
			if got := m.NumCalls(); got != tt.want {
				t.Errorf("NumCalls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_methodType_newArgv(t *testing.T) {
	type fields struct {
		method    reflect.Method
		ArgType   reflect.Type
		ReplyType reflect.Type
		numCalls  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   reflect.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &methodType{
				method:    tt.fields.method,
				ArgType:   tt.fields.ArgType,
				ReplyType: tt.fields.ReplyType,
				numCalls:  tt.fields.numCalls,
			}
			if got := m.newArgv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newArgv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_methodType_newReply(t *testing.T) {
	type fields struct {
		method    reflect.Method
		ArgType   reflect.Type
		ReplyType reflect.Type
		numCalls  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   reflect.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &methodType{
				method:    tt.fields.method,
				ArgType:   tt.fields.ArgType,
				ReplyType: tt.fields.ReplyType,
				numCalls:  tt.fields.numCalls,
			}
			if got := m.newReply(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// TODO: Add test cases.
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

	type fields struct {
		name   string
		typ    reflect.Type
		rec    reflect.Value
		method map[string]*methodType
	}
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
		// TODO: Add test cases.
		{"test1", args{
			m:     mType,
			argv:  argv,
			reply: reply,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService(&foo)
			err := s.call(tt.args.m, tt.args.argv, tt.args.reply)
			if (err == nil && *reply.Interface().(*int) == 4 && mType.NumCalls() == 1) != tt.wantErr {
				t.Errorf("call() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_registerMethods(t *testing.T) {
	type fields struct {
		name   string
		typ    reflect.Type
		rec    reflect.Value
		method map[string]*methodType
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//s := &service{
			//	name:   tt.fields.name,
			//	typ:    tt.fields.typ,
			//	rec:    tt.fields.rec,
			//	method: tt.fields.method,
			//}
		})
	}
}
