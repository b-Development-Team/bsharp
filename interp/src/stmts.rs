use super::*;

impl Interp {
    pub fn exec(&mut self, node: &IRNode) -> Result<Value, InterpError> {
        match node.data {
            _ => Err(InterpError::UnknownNode(node.clone())),
        }
    }
}
