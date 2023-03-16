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

    pub fn build_match(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        if args.len() < 2 {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: 2,
                got: args.len(),
            });
        }
        let val = self.build_node(&args[0]);

        let typecase = if let TypeData::ENUM(vals) = val.typ(self).data.concrete(self) {
            self.scopes
                .push(Scope::new(ScopeKind::TypeMatch(vals.clone()), pos));
            true
        } else {
            let v = val.typ(self);
            match v.data {
                TypeData::INVALID
                | TypeData::DEF(0)
                | TypeData::CHAR
                | TypeData::INT
                | TypeData::FLOAT => {}
                _ => {
                    return Err(IRError::InvalidArgument {
                        expected: TypeData::CHAR,
                        got: val,
                    })
                }
            };
            self.scopes
                .push(Scope::new(ScopeKind::Match(v.clone()), pos));
            false
        };
        self.stack.push(self.scopes.len() - 1);

        let mut cases = Vec::new();
        for arg in args.iter().skip(1) {
            let v = self.build_node(arg);
            if !match v.data {
                IRNodeData::Invalid => true,
                IRNodeData::TypeCase { .. } => typecase,
                IRNodeData::Case { .. } => !typecase,
                _ => false,
            } {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::PARAM,
                    got: v,
                });
            }
            cases.push(v);
        }
        self.stack.pop().unwrap();

        Ok(IRNode::new(
            if typecase {
                IRNodeData::TypeMatch {
                    val: Box::new(val),
                    body: cases,
                }
            } else {
                IRNodeData::Match {
                    val: Box::new(val),
                    body: cases,
                }
            },
            range,
            pos,
        ))
    }
}
