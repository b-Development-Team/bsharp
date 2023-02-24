use super::*;

impl IR {
    pub fn build_param(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(
            pos,
            args,
            &vec![ASTNodeDataType::Variable, ASTNodeDataType::Any],
        )?;
        let name = match &args[0].data {
            ASTNodeData::Variable(v) => v.clone(),
            _ => unreachable!(),
        };
        let typ = match self.build_node(&args[1]).data {
            IRNodeData::Type(t) => t,
            IRNodeData::Invalid => Type::from(TypeData::INVALID),
            _ => unreachable!(),
        };

        Ok(IRNode::new(IRNodeData::Param { name, typ }, range, pos))
    }

    pub fn build_returns(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let par = &self.typecheck(pos, args, &vec![TypeData::TYPE])?[0];
        let typ = match &par.data {
            IRNodeData::Type(t) => t.clone(),
            IRNodeData::Invalid => Type::from(TypeData::INVALID),
            _ => unreachable!(),
        };
        Ok(IRNode::new(IRNodeData::Returns(typ), range, pos))
    }
}
