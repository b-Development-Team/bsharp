use std::cell::RefCell;

use super::*;

impl Interp {
    pub fn exec_arrlit(&mut self, args: &Vec<IRNode>) -> Result<Value, InterpError> {
        let vals = self.exec_args(args)?;
        Ok(Value::Array(RefCell::new(vals)))
    }
}
