package Panizza

import (
	"log"
	"strings"
	"sync"
	"container/list"
	"reflect"
)

//全局存放注册的Aspect
var Aspecters = sync.Map{}

type Aspect interface {
	Config() string                                                    //配置切面 return "xxxxx|xxxxx|xxxxx"  xxxxx表示需要配置的目标函数名称
	Before(ctx *HandleContext, HandleName string)                      //目标方法执行之前执行
	After(ctx *HandleContext, HandleName string)                       //目标方法执行之后执行
	AfterPanic(err interface{}, ctx *HandleContext, HandleName string) //目标方法出现异常后执行
}

//传入Aspect指针，注册Aspect
func RegistAspecter(aspect interface{}) {
	t := reflect.TypeOf(aspect)
	v := reflect.ValueOf(aspect)
	for i := 0; i < v.NumMethod(); i++ {
		name := t.Method(i).Name
		if name == "OnCreate" {
			v.Method(i).Call([]reflect.Value{})
			continue
		}
	}
	handleNamesString := (aspect.(Aspect)).Config()
	if handleNamesString == "" {
		return
	}
	log.Println("Aspect at HandleName:", handleNamesString)

	handleNamesSlice := strings.Split(handleNamesString, "|")

	m := map[string]*list.List{}

	for key, val := range controllerInstence.methodTree {
		li := val
		for _, name := range handleNamesSlice {
			for e := li.Front(); e != nil; e = e.Next(){
				n := e.Value.(*node)
				if n.HandleName == name{
					n.HasAspect = true
					e.Value = n
				}
			}
		}
		m[key] = li
	}


	for k, v := range m {
		controllerInstence.methodTree[k] = v
	}


	for _, v := range handleNamesSlice {
		Aspecters.Store(v, aspect)
	}
}
