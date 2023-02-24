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

        // Add variable
        self.variables.push(Variable {
            name: name.clone(),
            typ: typ.clone(),
            scope: 0,
            definition: pos,
        });

        Ok(IRNode::new(
            IRNodeData::Param(self.variables.len() - 1),
            range,
            pos,
        ))
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

    pub fn build_block(&mut self, pos: Pos, args: &Vec<ASTNode>) -> Result<IRNode, IRError> {
        let mut body = Vec::new();
        self.scopes.push(Scope::new(ScopeKind::Block, pos));
        self.stack.push(self.scopes.len() - 1);
        for arg in args.iter() {
            let n = self.build_node(arg);
            match n.data {
                IRNodeData::Invalid | IRNodeData::Void => {}
                _ => body.push(n),
            }
        }
        Ok(IRNode::new(
            IRNodeData::Block {
                scope: self.stack.pop().unwrap(),
                body,
            },
            pos,
            pos,
        ))
    }

    pub fn build_return(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        // Get ret type
        let mut id = usize::MAX;
        for sc in self.stack.iter().rev() {
            if let ScopeKind::Function(ind) = &self.scopes[*sc].kind {
                id = *ind;
                break;
            }
        }
        if id == usize::MAX {
            return Err(IRError::ReturnStatementOutsideFunction(pos));
        }
        let ret_typ = self.funcs[id].ret_typ.clone();

        // Typecheck
        let expected_typ = if ret_typ.data == TypeData::VOID {
            vec![]
        } else {
            vec![ret_typ.data]
        };
        let par = self.typecheck(pos, args, &expected_typ)?;

        if expected_typ.len() == 0 {
            return Ok(IRNode::new(IRNodeData::Return(None), range, pos));
        }

        Ok(IRNode::new(
            IRNodeData::Return(Some(Box::new(par[0].clone()))),
            range,
            pos,
        ))
    }
}
