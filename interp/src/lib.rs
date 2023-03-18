use ir::*;

mod values;
use values::*;

mod errors;
use errors::*;

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

    pub fn run_fn(
        &mut self,
        func: &String,
        args: Vec<Value>,
    ) -> Result<Option<Value>, InterpError> {
        let funind = self.ir.funcs.iter().position(|x| x.name == *func).unwrap();

        // Put in params
        let mut frame = StackFrame::default();
        for (i, arg) in args.iter().enumerate() {
            frame
                .vars
                .insert(self.ir.funcs[funind].params[i], arg.clone());
        }
        self.stack.push(frame);

        // Run
        self.exec(&self.ir.funcs[funind].body.clone())?;

        // Get return
        let ret = self.stack.last().unwrap().ret.clone();
        self.stack.pop();

        Ok(ret)
    }
}
