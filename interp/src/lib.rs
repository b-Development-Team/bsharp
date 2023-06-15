use ir::*;
use std::io::Write;

mod values;
use values::*;

mod errors;
use errors::*;

mod fns;
mod stmts;

pub struct Interp {
    pub stack: Vec<StackFrame>,
    pub ir: IR,
    pub out: Box<dyn Write>,
}

impl Interp {
    pub fn new(ir: IR, out: Box<dyn Write>) -> Self {
        Self {
            stack: Vec::new(),
            ir,
            out,
        }
    }

    pub fn run_fn(
        &mut self,
        func: usize,
        args: Vec<Value>,
        callpos: Pos,
    ) -> Result<Value, InterpError> {
        // Put in params
        let mut frame = StackFrame::default();
        frame.pos = callpos;
        frame.func = func;
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
