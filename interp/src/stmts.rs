use super::*;

impl Interp {
    pub fn exec(&mut self, node: &IRNode) -> Result<Value, InterpError> {
        match &node.data {
            IRNodeData::Block { body, .. } => self.exec_block(&body),
            IRNodeData::FnCall { func, args, .. } => {
                let args = self.exec_args(args)?;
                self.run_fn(*func, args)
            }
            IRNodeData::Boolean(left, op, right) => self.exec_boolop(*op, left, right),
            IRNodeData::Comparison(left, op, right) => self.exec_comp(*op, left, right),
            IRNodeData::Int(v) => Ok(Value::Int(*v)),
            IRNodeData::Float(v) => Ok(Value::Float(*v)),
            IRNodeData::Char(v) => Ok(Value::Char(*v)),
            _ => Err(InterpError::UnknownNode(node.clone())),
        }
    }
}
