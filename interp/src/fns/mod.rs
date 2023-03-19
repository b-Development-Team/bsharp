use super::*;

mod control;
mod ops;
mod vars;

impl Interp {
    pub fn exec_args(&mut self, args: &Vec<IRNode>) -> Result<Vec<Value>, InterpError> {
        let mut res = Vec::with_capacity(args.len());
        for arg in args {
            res.push(self.exec(arg)?);
        }
        Ok(res)
    }
}
