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
}
