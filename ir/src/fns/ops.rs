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
            (TypeData::INT, TypeData::INT)
            | (TypeData::FLOAT, TypeData::FLOAT)
            | (TypeData::CHAR, TypeData::CHAR) => {}
            (TypeData::INTERFACE { .. }, TypeData::INTERFACE { .. }) => {
                let ok = if left
                    .typ()
                    .data
                    .matches(&vec![TypeData::INT, TypeData::FLOAT, TypeData::CHAR], self)
                {
                    match (left.typ().data, right.typ().data) {
                        (TypeData::DEF(a), TypeData::DEF(b)) => a == b,
                        _ => false,
                    }
                } else {
                    false
                };
                if !ok {
                    return Err(IRError::InvalidArgument {
                        expected: TypeData::FLOAT,
                        got: left,
                    });
                }
            }
            (TypeData::INVALID, _) | (_, TypeData::INVALID) => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::FLOAT,
                    got: left,
                });
            }
        };

        Ok(IRNode::new(
            IRNodeData::Math(Box::new(left), op, Box::new(right)),
            range,
            pos,
        ))
    }

    pub fn build_compop(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
        op: ComparisonOperator,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any, ASTNodeDataType::Any])?;
        let left = self.build_node(&args[0]);
        let right = self.build_node(&args[1]);

        // Check types
        match (
            left.typ().data.concrete(self),
            right.typ().data.concrete(self),
        ) {
            (TypeData::INT, TypeData::INT)
            | (TypeData::FLOAT, TypeData::FLOAT)
            | (TypeData::CHAR, TypeData::CHAR)
            | (TypeData::BOOL, TypeData::BOOL) => {}
            (TypeData::INTERFACE { .. }, TypeData::INTERFACE { .. }) => {
                let ok = if left.typ().data.matches(
                    &vec![
                        TypeData::INT,
                        TypeData::FLOAT,
                        TypeData::CHAR,
                        TypeData::BOOL,
                    ],
                    self,
                ) {
                    match (left.typ().data, right.typ().data) {
                        (TypeData::DEF(a), TypeData::DEF(b)) => a == b,
                        _ => false,
                    }
                } else {
                    false
                };
                if !ok {
                    return Err(IRError::InvalidArgument {
                        expected: TypeData::FLOAT,
                        got: left,
                    });
                }
            }
            (TypeData::INVALID, _) | (_, TypeData::INVALID) => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::FLOAT,
                    got: left,
                });
            }
        };

        Ok(IRNode::new(
            IRNodeData::Comparison(Box::new(left), op, Box::new(right)),
            range,
            pos,
        ))
    }
}
