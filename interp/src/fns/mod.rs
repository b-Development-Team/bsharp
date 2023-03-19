use super::*;

mod composite;
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

    pub fn exec_print(&mut self, val: &IRNode) -> Result<Value, InterpError> {
        let chars: Vec<_> = match self.exec(val)? {
            Value::Array(vals) => vals
                .borrow()
                .iter()
                .map(|v| match v {
                    Value::Char(c) => *c,
                    _ => unreachable!(),
                })
                .collect(),
            _ => unreachable!(),
        };
        self.out.write_all(&chars)?;
        Ok(Value::Void)
    }
}
