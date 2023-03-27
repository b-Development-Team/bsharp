use std::{cell::RefCell, collections::HashMap, rc::Rc};

use super::*;

impl Interp {
    pub fn exec_arrlit(&mut self, args: &Vec<IRNode>) -> Result<Value, InterpError> {
        let vals = self.exec_args(args)?;
        Ok(Value::Array(Rc::new(RefCell::new(vals))))
    }

    pub fn exec_newstruct(&mut self, args: &Vec<IRNode>) -> Result<Value, InterpError> {
        let mut fields = HashMap::new();
        for arg in args {
            match &arg.data {
                IRNodeData::StructOp { field, val } => {
                    let val = self.exec(val)?;
                    fields.insert(field.clone(), val);
                }
                _ => unreachable!(),
            }
        }
        Ok(Value::Struct(Rc::new(RefCell::new(fields))))
    }

    pub fn exec_newenum(&mut self, arg: &IRNode) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        Ok(Value::Enum(arg.typ(&self.ir), Box::new(val)))
    }

    pub fn exec_newbox(&mut self, arg: &IRNode) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        Ok(Value::Box(arg.typ(&self.ir), Box::new(val)))
    }

    pub fn exec_peek(&mut self, arg: &IRNode, t: &Type) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Box(typ, _) = val {
            Ok(Value::Bool(typ == *t))
        } else {
            unreachable!()
        }
    }

    pub fn exec_unbox(&mut self, arg: &IRNode, t: &Type) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Box(typ, val) = val {
            if typ == *t {
                Ok(*val)
            } else {
                Err(InterpError::InvalidBoxType {
                    pos: arg.pos,
                    expected: t.clone(),
                    got: typ,
                })
            }
        } else {
            unreachable!()
        }
    }

    pub fn exec_newtuple(&mut self, args: &Vec<IRNode>) -> Result<Value, InterpError> {
        let vals = self.exec_args(args)?;
        Ok(Value::Tuple(Rc::new(RefCell::new(vals))))
    }

    pub fn exec_getenum(&mut self, arg: &IRNode, typ: &Type) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Enum(t, v) = val {
            if t == *typ {
                Ok(*v)
            } else {
                Err(InterpError::InvalidEnumType {
                    pos: arg.pos,
                    expected: typ.clone(),
                    got: t,
                })
            }
        } else {
            unreachable!()
        }
    }

    pub fn exec_getstruct(&mut self, arg: &IRNode, field: &String) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Struct(v) = val {
            let v = v.borrow();
            if let Some(val) = v.get(field) {
                Ok(val.clone())
            } else {
                unreachable!()
            }
        } else {
            unreachable!()
        }
    }

    pub fn exec_setstruct(
        &mut self,
        arg: &IRNode,
        ops: &Vec<IRNode>,
    ) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Struct(v) = val {
            let mut v = v.borrow_mut();
            for op in ops {
                match &op.data {
                    IRNodeData::StructOp { field, val } => {
                        let val = self.exec(val)?;
                        v.insert(field.clone(), val);
                    }
                    _ => unreachable!(),
                }
            }

            Ok(Value::Void)
        } else {
            unreachable!()
        }
    }

    pub fn exec_gettuple(&mut self, arg: &IRNode, idx: &usize) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Tuple(v) = val {
            let v = v.borrow();
            if let Some(val) = v.get(*idx) {
                Ok(val.clone())
            } else {
                unreachable!()
            }
        } else {
            unreachable!()
        }
    }

    pub fn exec_len(&mut self, arg: &IRNode) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Array(v) = val {
            let v = v.borrow();
            Ok(Value::Int(v.len() as i64))
        } else {
            unreachable!()
        }
    }

    pub fn exec_getarr(&mut self, arg: &IRNode, idx: &IRNode) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        if let Value::Array(v) = val {
            let idx = self.exec(idx)?;
            if let Value::Int(idx) = idx {
                let v = v.borrow();
                if let Some(val) = v.get(idx as usize) {
                    Ok(val.clone())
                } else {
                    Err(InterpError::ArrayIndexOutOfBounds {
                        pos: arg.pos,
                        len: v.len(),
                        index: idx as usize,
                    })
                }
            } else {
                unreachable!()
            }
        } else {
            unreachable!()
        }
    }

    pub fn exec_newarr(&mut self, cap: &Option<Box<IRNode>>) -> Result<Value, InterpError> {
        let cap = if let Some(n) = cap {
            let n = self.exec(n)?;
            if let Value::Int(n) = n {
                n
            } else {
                unreachable!()
            }
        } else {
            0
        };

        Ok(Value::Array(Rc::new(RefCell::new(Vec::with_capacity(
            cap as usize,
        )))))
    }

    pub fn exec_append(&mut self, arg: &IRNode, val: &IRNode) -> Result<Value, InterpError> {
        let val = self.exec(val)?;
        let arg = self.exec(arg)?;
        if let Value::Array(v) = arg {
            let mut v = v.borrow_mut();
            v.push(val);
            Ok(Value::Void)
        } else {
            unreachable!()
        }
    }

    pub fn exec_setarr(
        &mut self,
        arg: &IRNode,
        ind: &IRNode,
        val: &IRNode,
    ) -> Result<Value, InterpError> {
        let v = self.exec(arg)?;
        let ind = self.exec(ind)?;
        let val = self.exec(val)?;
        match (v, ind, val) {
            (Value::Array(arr), Value::Int(ind), v) => {
                let mut arr = arr.borrow_mut();
                if let Some(val) = arr.get_mut(ind as usize) {
                    *val = v;
                    Ok(Value::Void)
                } else {
                    Err(InterpError::ArrayIndexOutOfBounds {
                        pos: arg.pos,
                        len: arr.len(),
                        index: ind as usize,
                    })
                }
            }
            _ => unreachable!(),
        }
    }
}
