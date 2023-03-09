use super::*;

impl IR {
    pub fn build_while(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let pars = self.typecheck(pos, args, &vec![TypeData::BOOL, TypeData::VOID])?;
        match pars[1].data {
            IRNodeData::Block { .. } => {}
            IRNodeData::Invalid => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::VOID,
                    got: pars[1].clone(),
                })
            }
        }

        Ok(IRNode::new(
            IRNodeData::While {
                cond: Box::new(pars[0].clone()),
                body: Box::new(pars[1].clone()),
            },
            range,
            pos,
        ))
    }

    pub fn build_if(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(
            pos,
            args,
            &vec![
                ASTNodeDataType::Any,
                ASTNodeDataType::Any,
                ASTNodeDataType::Any,
            ],
        )?;
        let cond = self.build_node(&args[0]);
        cond.typ(self).expect(pos, &Type::from(TypeData::BOOL))?;
        let body = self.build_node(&args[1]);
        let els = self.build_node(&args[2]);
        let ret_typ = if body.typ(self) == els.typ(self) && body.typ(self).data != TypeData::VOID {
            Some(body.typ(self))
        } else {
            None
        };

        Ok(IRNode::new(
            IRNodeData::If {
                body: Box::new(body),
                els: Box::new(els),
                cond: Box::new(cond),
                ret_typ,
            },
            range,
            pos,
        ))
    }
}
