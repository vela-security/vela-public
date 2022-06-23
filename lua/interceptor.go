package lua

type interceptorFunc func(*LState, uint32, *callFrame) bool

type IndexEx interface {
	Index(*LState, string) LValue
}

type NewIndexEx interface {
	NewIndex(*LState, string, LValue)
}

//MetaTableEx 通过 string 获取 内容 a:key()
type MetaTableEx interface {
	MetaTable(*LState, string) LValue
}

//MetaEx  通过 LValue 获取 a[1]
type MetaEx interface {
	Meta(*LState, LValue) LValue
}

//NewMetaEx 通过LValue 设置 a[1] = 123
type NewMetaEx interface {
	NewMeta(*LState, LValue, LValue)
}

func interceptorCall(L *LState, inst uint32, baseframe *callFrame) bool {
	/*
		reg := L.reg
		RA := L.currentFrame.LocalBase + (int(inst>>18) & 0xff)
		lv := reg.GetCodeVM(RA)
		if lv.Typ() != LTGFunction {
			return false
		}

		B := int(inst & 0x1ff)    //GETB
		C := int(inst>>9) & 0x1ff //GETC
		nargs := B - 1
		top := reg.Top()
		if B == 0 {
			nargs = top - (RA + 1)
		}
		nret := C - 1

		lv.(GFunction).xcall(L , reg , RA , nargs , nret)
	*/
	return false
}

func interceptorSetTableHelper(L *LState, v interface{}, key LValue, val LValue) {
	switch vx := v.(type) {
	case NewMetaEx:
		vx.NewMeta(L, key, val)
	default:
		L.RaiseError("interceptor attempt to Index a non-table object(%v) with key '%s'", L.Type().String(), key.String())
	}
}

func interceptorSetTable(L *LState, inst uint32, baseframe *callFrame) bool {
	reg := L.reg
	cf := L.currentFrame
	lbase := cf.LocalBase
	A := int(inst>>18) & 0xff //GETA
	RA := lbase + A
	lv := reg.Get(RA).Peek()
	B := int(inst & 0x1ff)    //GETB
	C := int(inst>>9) & 0x1ff //GETC

	key := L.rkValue(B)
	val := L.rkValue(C)

	switch lv.Type() {
	case LTProcData:
		lv.(*ProcData).Data.NewMeta(L, key, val)
	case LTObject:
		interceptorSetTableHelper(L, lv, key, val)
	case LTAnyData:
		interceptorSetTableHelper(L, lv.(*AnyData).Data, key, val)
	case LTSlice:
		lv.(Slice).NewMeta(L, key, val)
	case LTMap:
		lv.(*Map).NewMeta(L, key, val)

	case LTKv, LTSkv:
		L.RaiseError("interceptor attempt to Index a non-table object(%v) with key '%s'", L.Type().String(), key.String())
	default:
		return false
	}

	return true
}

func interceptorSetTableEks(L *LState, inst uint32, baseframe *callFrame) bool {
	reg := L.reg
	RA := L.currentFrame.LocalBase + (int(inst>>18) & 0xff)
	lv := reg.Get(RA).Peek()
	B := int(inst & 0x1ff)    //GETB
	C := int(inst>>9) & 0x1ff //GETC
	key := L.rkString(B)

	switch lv.Type() {
	case LTProcData:
		lv.(*ProcData).Data.NewIndex(L, key, L.rkValue(C))
		return true
	case LTAnyData:
		lv.(*AnyData).NewIndex(L, key, L.rkValue(C))
		return true

	case LTMap:
		lv.(*Map).NewIndex(L, key, L.rkValue(C))
		return true

	case LTSlice:
		//goto
		return true

	case LTObject:
		switch vx := lv.(type) {
		case NewIndexEx:
			vx.NewIndex(L, key, L.rkValue(C))
		default:
			L.RaiseError("interceptor attempt to Index a object not found  with key '%s'", key)
		}

		return true

	case LTKv, LTSkv:
		L.RaiseError("interceptor attempt to Index a non-table object(%v) with key '%s'", L.Type().String(), key)
		return true
	}
	return false
}

func interceptorGetTableHelp(L *LState, v interface{}, key LValue) LValue {
	switch vx := v.(type) {
	case MetaEx:
		return vx.Meta(L, key)

	default:
		L.RaiseError("attempt to meta %v with %v", L.Type().String(), key)
	}
	return LNil
}

func interceptorGetTable(L *LState, inst uint32, baseframe *callFrame) bool {
	reg := L.reg
	RA := L.currentFrame.LocalBase + (int(inst>>18) & 0xff)
	B := L.currentFrame.LocalBase + int(inst&0x1ff) //GETB
	C := int(inst>>9) & 0x1ff                       //GETC
	lv := reg.Get(B).Peek()

	switch lv.Type() {
	case LTProcData:
		reg.Set(RA, lv.(*ProcData).Data.Meta(L, L.rkValue(C)))
		return true

	case LTObject:
		reg.Set(RA, interceptorGetTableHelp(L, lv, L.rkValue(C)))
		return true

	case LTAnyData:
		reg.Set(RA, interceptorGetTableHelp(L, lv.(*AnyData).Data, L.rkValue(C)))

	case LTSlice:
		reg.Set(RA, lv.(Slice).Meta(L, L.rkValue(C)))
		return true

	case LTMap:
		reg.Set(RA, lv.(*Map).Meta(L, L.rkValue(C)))
		return true

	case LTKv, LTSkv:
		L.RaiseError("attempt to meta %v with %v", L.Type().String(), L.rkValue(C))
		return true

	}
	return false
}

func interceptorGetTableEks(L *LState, inst uint32, baseframe *callFrame) bool {
	reg := L.reg
	RA := L.currentFrame.LocalBase + (int(inst>>18) & 0xff)
	B := L.currentFrame.LocalBase + int(inst&0x1ff) //GETB
	C := int(inst>>9) & 0x1ff                       //GETC
	lv := reg.Get(B).Peek()

	switch lv.Type() {
	case LTProcData:
		reg.Set(RA, lv.(*ProcData).Data.Index(L, L.rkString(C)))
		return true
	case LTAnyData:
		reg.Set(RA, lv.(*AnyData).Index(L, L.rkString(C)))
		return true
	case LTKv:
		reg.Set(RA, lv.(*userKV).Get(L.rkString(C)))
		return true
	case LTSkv:
		reg.Set(RA, lv.(*safeUserKV).Get(L.rkString(C)))
		return true

	case LTMap:
		reg.Set(RA, lv.(*Map).Index(L, L.rkString(C)))
		return true

	case LTSlice:
		reg.Set(RA, lv.(Slice).Index(L, L.rkString(C)))
		//L.RaiseError("interceptor slice Index not found")
		return true

	case LTObject:
		if obj, ok := lv.(IndexEx); ok {
			reg.Set(RA, obj.Index(L, L.rkString(C)))
		} else {
			L.RaiseError("interceptor attempt to Index a object(%v) with %v", L.Type().String(), L.rkValue(C))
		}
		return true
	}
	return false
}

func interceptorSelf(L *LState, inst uint32, baseframe *callFrame) bool {
	reg := L.reg
	cf := L.currentFrame
	lbase := cf.LocalBase
	A := int(inst>>18) & 0xff //GETA
	RA := lbase + A
	B := int(inst & 0x1ff)    //GETB
	C := int(inst>>9) & 0x1ff //GETC
	obj := reg.Get(lbase + B)
	switch obj.Type() {
	case LTProcData:
		reg.Set(RA, obj.(*ProcData).Data.Index(L, L.rkString(C)))
		reg.Set(RA+1, obj)
		return true
	case LTAnyData:
		reg.Set(RA, obj.(*AnyData).Index(L, L.rkString(C)))
		reg.Set(RA+1, obj)
		return true
	case LTKv:
		reg.Set(RA, obj.(*userKV).Get(L.rkString(C)))
		reg.Set(RA+1, obj)
		return true
	case LTSkv:
		reg.Set(RA, obj.(*safeUserKV).Get(L.rkString(C)))
		reg.Set(RA+1, obj)
		return true
	case LTSlice:
		obj.(Slice).MetaTable(L, L.rkString(C))
		return true
	case LTMap:
		obj.(*Map).MetaTable(L, L.rkString(C))
		return true
	case LTObject:
		if lv, ok := obj.(MetaTableEx); ok {
			reg.Set(RA, lv.MetaTable(L, L.rkString(C)))
			reg.Set(RA+1, obj)
		} else {
			L.RaiseError("attempt to Index a object(%v) with %v", L.Type().String(), L.rkValue(C))
		}
		return true
	}

	return false
}

func interceptorGetG(L *LState, inst uint32, baseframe *callFrame) bool {
	return false
}

func interceptorPass(L *LState, inst uint32, baseframe *callFrame) bool {
	return false
}

func interceptorTable(op int) interceptorFunc {
	switch op {
	case OP_SELF:
		return interceptorSelf
	case OP_CALL:
		return interceptorCall
	case OP_SETTABLE:
		return interceptorSetTable
	case OP_SETTABLEKS:
		return interceptorSetTableEks
	case OP_GETTABLE:
		return interceptorGetTable
	case OP_GETTABLEKS:
		return interceptorGetTableEks
	case OP_GETGLOBAL:
		return interceptorGetG
	}

	return interceptorPass
}
