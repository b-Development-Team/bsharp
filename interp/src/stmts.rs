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
            IRNodeData::If {
                cond,
                body,
                els,
                ret_typ,
            } => self.exec_if(cond, body, els, ret_typ),
            IRNodeData::Variable(ind, _) => self.exec_var(*ind),
            IRNodeData::NewArrayLiteral(_, vals) => self.exec_arrlit(vals),
            IRNodeData::Print(v) => self.exec_print(v),
            IRNodeData::Define { var, val, .. } => self.exec_define(*var, val),
            _ => Err(InterpError::UnknownNode(node.clone())),
        }
    }
}
