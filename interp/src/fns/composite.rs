use std::{cell::RefCell, collections::HashMap};

use super::*;

impl Interp {
    pub fn exec_arrlit(&mut self, args: &Vec<IRNode>) -> Result<Value, InterpError> {
        let vals = self.exec_args(args)?;
        Ok(Value::Array(RefCell::new(vals)))
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
        Ok(Value::Struct(RefCell::new(fields)))
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
        Ok(Value::Tuple(RefCell::new(vals)))
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
}
