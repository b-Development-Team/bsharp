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

    pub fn build_case(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        // Get scope
        let kind = match &self.scopes[*self.stack.last().unwrap()].kind {
            ScopeKind::TypeMatch(vals) => ScopeKind::TypeMatch(vals.clone()),
            ScopeKind::Match(val) => ScopeKind::Match(val.clone()),
            _ => return Err(IRError::CaseOutsideMatch(range)),
        };

        // Build args
        let typs = if let ScopeKind::TypeMatch(_) = kind {
            vec![
                ASTNodeDataType::Variable,
                ASTNodeDataType::Any,
                ASTNodeDataType::Block,
            ]
        } else {
            vec![ASTNodeDataType::Any, ASTNodeDataType::Block]
        };
        let args = typecheck_ast(pos, &args, &typs)?;

        // Create scope
        let mut varind = None;
        let mut arg = None;
        let blkind = if let ScopeKind::TypeMatch(ref allowed) = kind {
            let name = match &args[0].data {
                ASTNodeData::Variable(name) => name.clone(),
                _ => unreachable!(),
            };
            let t = self.build_node(&args[1]);
            let typ = match &t.data {
                IRNodeData::Invalid => Type::from(TypeData::INVALID),
                IRNodeData::Type(t) => t.clone(),
                _ => {
                    return Err(IRError::InvalidArgument {
                        expected: TypeData::TYPE,
                        got: t,
                    })
                }
            };
            if !allowed.iter().any(|x| x.data == typ.data) {
                return Err(IRError::InvalidArgument {
                    expected: allowed[0].data.clone(),
                    got: t,
                });
            }
            self.variables.push(Variable {
                name: name.clone(),
                typ: typ.clone(),
                scope: self.scopes.len(),
                definition: args[0].pos,
            });
            varind = Some((self.variables.len() - 1, typ));
            self.scopes.push(Scope::new(ScopeKind::Case, pos));
            self.scopes
                .last_mut()
                .unwrap()
                .vars
                .insert(name, self.variables.len() - 1);
            self.stack.push(self.scopes.len() - 1);

            2
        } else {
            let par = self.build_node(&args[0]);
            let ok = match &par.data {
                IRNodeData::NewArrayLiteral(t, _) => t.data == TypeData::DEF(0),
                IRNodeData::Invalid => true,
                IRNodeData::Char(_) | IRNodeData::Int(_) | IRNodeData::Float(_) => true,
                _ => false,
            };
            if !ok {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::INT,
                    got: par,
                });
            }
            arg = Some(par);
            self.scopes.push(Scope::new(ScopeKind::Case, pos));
            self.stack.push(self.scopes.len() - 1);

            1
        };
        let blk = self.build_node(&args[blkind]);

        // Build node
        let node = if let ScopeKind::TypeMatch(_) = kind {
            IRNodeData::TypeCase {
                var: varind.as_ref().unwrap().0.clone(),
                typ: varind.as_ref().unwrap().1.clone(),
                body: Box::new(blk),
            }
        } else {
            IRNodeData::Case {
                val: Box::new(arg.unwrap()),
                body: Box::new(blk),
            }
        };

        self.stack.pop().unwrap();

        Ok(IRNode::new(node, range, pos))
    }
}
