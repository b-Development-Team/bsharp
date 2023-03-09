use super::*;

impl IR {
    pub fn build_box(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        let val = self.build_node(&args[0]);
        Ok(IRNode::new(IRNodeData::NewBox(Box::new(val)), range, pos))
    }

    pub fn build_peek(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = self.typecheck(pos, &args, &vec![TypeData::BOX, TypeData::TYPE])?;
        Ok(IRNode::new(
            IRNodeData::Peek {
                bx: Box::new(args[0].clone()),
                typ: match &args[1].data {
                    IRNodeData::Type(t) => t.clone(),
                    IRNodeData::Invalid => Type::from(TypeData::INVALID),
                    _ => unreachable!(),
                },
            },
            range,
            pos,
        ))
    }

    pub fn build_unbox(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = self.typecheck(pos, &args, &vec![TypeData::BOX, TypeData::TYPE])?;
        Ok(IRNode::new(
            IRNodeData::Unbox {
                bx: Box::new(args[0].clone()),
                typ: match &args[1].data {
                    IRNodeData::Type(t) => t.clone(),
                    IRNodeData::Invalid => Type::from(TypeData::INVALID),
                    _ => unreachable!(),
                },
            },
            range,
            pos,
        ))
    }
}
