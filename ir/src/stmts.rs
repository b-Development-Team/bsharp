use super::*;

impl super::IR {
    pub fn build_node(&mut self, v: &ASTNode) -> Result<IRNode, IRError> {
        match v.data {
            ASTNodeData::Comment(_) => Ok(IRNode::new(IRNodeData::Void, v.pos)),
            _ => Err(IRError::UnexpectedNode(v.clone())),
        }
    }
}
