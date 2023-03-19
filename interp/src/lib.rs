use ir::*;

mod values;
use values::*;

mod errors;
use errors::*;

mod fns;
mod stmts;

pub struct Interp {
    pub stack: Vec<StackFrame>,
    pub ir: IR,
}

impl Interp {
    pub fn new(ir: IR) -> Self {
        Self {
            stack: Vec::new(),
            ir,
        }
    }

    pub fn run_fn(&mut self, func: usize, args: Vec<Value>) -> Result<Value, InterpError> {
        // Put in params
        let mut frame = StackFrame::default();
        for (i, arg) in args.iter().enumerate() {
            frame
                .vars
                .insert(self.ir.funcs[func].params[i], arg.clone());
        }
        self.stack.push(frame);

        // Run
        self.exec(&self.ir.funcs[func].body.clone())?;

        // Get return
        let ret = self.stack.last().unwrap().ret.clone();
        self.stack.pop();

        Ok(if let Some(ret) = ret {
            ret
        } else {
            Value::Void
        })
    }
}
