use super::*;

impl IR {
    pub fn build_cast_int(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        let val = self.build_node(&args[0]);
        match val.typ(self).data {
            TypeData::CHAR | TypeData::FLOAT | TypeData::INVALID => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::FLOAT,
                    got: val,
                });
            }
        }
        Ok(IRNode::new(
            IRNodeData::Cast(Box::new(val), Type::from(TypeData::INT)),
            range,
            pos,
        ))
    }

    pub fn build_cast_float(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        let val = self.build_node(&args[0]);
        match val.typ(self).data {
            TypeData::CHAR | TypeData::INT | TypeData::INVALID => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::INT,
                    got: val,
                });
            }
        }
        Ok(IRNode::new(
            IRNodeData::Cast(Box::new(val), Type::from(TypeData::FLOAT)),
            range,
            pos,
        ))
    }

    pub fn build_cast_char(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        let val = self.build_node(&args[0]);
        match val.typ(self).data {
            TypeData::FLOAT | TypeData::INT | TypeData::INVALID => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::INT,
                    got: val,
                });
            }
        }
        Ok(IRNode::new(
            IRNodeData::Cast(Box::new(val), Type::from(TypeData::CHAR)),
            range,
            pos,
        ))
    }
}
