use super::*;

impl Interp {
    pub fn exec_var(&mut self, val: usize) -> Result<Value, InterpError> {
        Ok(self.stack.last().unwrap().vars[&val].clone())
    }

    pub fn exec_define(&mut self, var: usize, val: &IRNode) -> Result<Value, InterpError> {
        let val = self.exec(val)?;
        self.stack.last_mut().unwrap().vars.insert(var, val.clone());
        Ok(val)
    }
}
