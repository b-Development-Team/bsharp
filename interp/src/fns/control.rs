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
}
