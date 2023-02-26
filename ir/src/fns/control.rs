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
}
