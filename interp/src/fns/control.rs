use super::*;

impl Interp {
    pub fn exec_block(&mut self, body: &Vec<IRNode>) -> Result<Value, InterpError> {
        for node in body {
            self.exec(node)?;
            if self.stack.last().unwrap().ret.is_some() {
                break; // Return
            }
        }
        Ok(Value::Void)
    }

    pub fn exec_if(
        &mut self,
        cond: &IRNode,
        body: &IRNode,
        els: &IRNode,
        ret_typ: &Option<Type>,
    ) -> Result<Value, InterpError> {
        let cond = self.exec(cond)?;
        let res = if let Value::Bool(cond) = cond {
            if cond {
                self.exec(body)?
            } else {
                self.exec(els)?
            }
        } else {
            unreachable!();
        };
        if ret_typ.is_some() {
            return Ok(res);
        }
        Ok(Value::Void)
    }

    pub fn exec_typematch(
        &mut self,
        cond: &IRNode,
        cases: &Vec<IRNode>,
    ) -> Result<Value, InterpError> {
        let cond = self.exec(cond)?;
        let (t, val) = match cond {
            Value::Enum(t, val) => (t, val),
            _ => unreachable!(),
        };
        for case in cases {
            let (var, typ, body) = match &case.data {
                IRNodeData::TypeCase { var, typ, body } => (var, typ, body),
                _ => unreachable!(),
            };
            if *typ == t {
                self.stack
                    .last_mut()
                    .unwrap()
                    .vars
                    .insert(*var, *val.clone());
                self.exec(body)?;
                break;
            }
        }
        Ok(Value::Void)
    }

    pub fn exec_match(&mut self, cond: &IRNode, cases: &Vec<IRNode>) -> Result<Value, InterpError> {
        let cond = self.exec(cond)?;
        for case in cases {
            let (val, body) = match &case.data {
                IRNodeData::Case { val, body } => (val, body),
                _ => unreachable!(),
            };

            let val = self.exec(val)?;
            if val == cond {
                self.exec(body)?;
                break;
            }
        }
        Ok(Value::Void)
    }
}
