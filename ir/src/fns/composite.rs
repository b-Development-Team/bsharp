use super::*;

impl IR {
    fn check_array_par(&mut self, par: &IRNode) -> Result<(), IRError> {
        match par.typ().data.concrete(self) {
            TypeData::ARRAY { .. } | TypeData::INVALID => {}
            TypeData::INTERFACE { body, .. } => {
                let mut found = false;
                for t in body {
                    match t.data.concrete(self) {
                        TypeData::ARRAY { .. } => found = true,
                        _ => {}
                    }
                }
                if !found {
                    return Err(IRError::InvalidArgument {
                        expected: TypeData::ARRAY {
                            param: None,
                            body: Box::new(Type::void()),
                            scope: 0,
                        },
                        got: par.clone(),
                    });
                }
            }
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::ARRAY {
                        param: None,
                        body: Box::new(Type::void()),
                        scope: 0,
                    },
                    got: par.clone(),
                })
            }
        }
        Ok(())
    }

    pub fn build_len(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let val = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        let par = self.build_node(&val[0]);
        self.check_array_par(&par)?;
        Ok(IRNode::new(IRNodeData::Len(Box::new(par)), range, pos))
    }
}
