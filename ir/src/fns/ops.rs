use super::*;

impl IR {
    pub fn build_mathop(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
        op: MathOperator,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any, ASTNodeDataType::Any])?;
        let left = self.build_node(&args[0]);
        let right = self.build_node(&args[1]);

        // Check types
        match (
            left.typ().data.concrete(self),
            right.typ().data.concrete(self),
        ) {
            (TypeData::INT, TypeData::INT) => Type::from(TypeData::INT),
            (TypeData::FLOAT, TypeData::FLOAT) => Type::from(TypeData::FLOAT),
            (TypeData::INVALID, _) | (_, TypeData::INVALID) => Type::from(TypeData::INVALID),
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::FLOAT,
                    got: left,
                })
            }
        };

        Ok(IRNode::new(
            IRNodeData::Math(Box::new(left), op, Box::new(right)),
            range,
            pos,
        ))
    }
}
