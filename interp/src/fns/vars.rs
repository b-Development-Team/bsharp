use super::*;

impl Interp {
    pub fn exec_var(&mut self, val: usize) -> Result<Value, InterpError> {
        Ok(self.stack.last().unwrap().vars[&val].clone())
    }
}
