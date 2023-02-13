use super::*;

impl super::IR {
    pub fn build_typ(&mut self, id: usize) -> Result<(), IRError> {
        let ast = self.types[id].ast.clone();
        if let Some(ast) = &ast {
            // Build
            let res = self.build_node(ast)?;
            res.typ().expect(res.pos, &Type::from(TypeData::TYPE))?;
            self.types[id].ast = None;
            self.types[id].typ = res.typ();
        }
        Ok(())
    }

    pub fn build_node(&mut self, v: &ASTNode) -> Result<IRNode, IRError> {
        match v.data {
            ASTNodeData::Comment(_) => Ok(IRNode::new(IRNodeData::Void, v.pos, v.pos)),
            _ => Err(IRError::UnexpectedNode(v.clone())),
        }
    }
}
