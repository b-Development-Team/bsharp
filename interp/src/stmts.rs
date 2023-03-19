use super::*;

impl Interp {
    pub fn exec(&mut self, node: &IRNode) -> Result<Value, InterpError> {
        match &node.data {
            IRNodeData::Block { body, .. } => self.exec_block(&body),
            IRNodeData::FnCall { func, args, .. } => {
                let args = self.exec_args(args)?;
                self.run_fn(*func, args)
            }
            _ => Err(InterpError::UnknownNode(node.clone())),
        }
    }
}
