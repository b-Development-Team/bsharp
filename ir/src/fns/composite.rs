use super::*;

fn expectedargcount(cnt: usize, pos: Pos, got: usize) -> IRError {
    IRError::InvalidArgumentCount {
        pos,
        expected: cnt,
        got,
    }
}

impl IR {
    pub fn build_len(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let val = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        let par = self.build_node(&val[0]);
        match par.typ(self).data.concrete(self) {
            TypeData::ARRAY(_) => {}
            TypeData::INVALID => {}
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::ARRAY(Box::new(Type::void())),
                    got: par,
                })
            }
        }
        Ok(IRNode::new(IRNodeData::Len(Box::new(par)), range, pos))
    }

    pub fn build_append(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let val = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any, ASTNodeDataType::Any])?;
        let arr = self.build_node(&val[0]);
        let body_typ = match arr.typ(self).data.concrete(self) {
            TypeData::ARRAY(body) => *body,
            TypeData::INVALID => Type::from(TypeData::INVALID),
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::ARRAY(Box::new(Type::void())),
                    got: arr,
                });
            }
        };
        let val = self.build_node(&val[1]);

        if val.typ(self).data.concrete(self) != body_typ.data.concrete(self)
            && val.typ(self).data.concrete(self) != TypeData::INVALID
        {
            return Err(IRError::InvalidArgument {
                expected: body_typ.data,
                got: val,
            });
        }
        Ok(IRNode::new(IRNodeData::Len(Box::new(val)), range, pos))
    }

    pub fn build_new(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut args = args.clone();
        if args.len() < 1 {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: 1,
                got: 0,
            });
        }
        let t = self.build_node(&args[0]);
        let typ = match &t.data {
            IRNodeData::Type(t) => t.clone(),
            _ => Type::from(TypeData::INVALID),
        };
        args.remove(0);

        let dat = match typ.data.concrete(self) {
            TypeData::INVALID => IRNodeData::Invalid,
            TypeData::ENUM(v) => {
                if args.len() != 1 {
                    return Err(expectedargcount(1, pos, args.len()));
                }
                let arg = self.build_node(&args[0]);
                let ok = v.iter().any(|v| v.data == arg.typ(self).data)
                    || arg.typ(self).data == TypeData::INVALID;
                if !ok {
                    return Err(IRError::InvalidArgument {
                        expected: v[0].data.clone(),
                        got: arg,
                    });
                }
                IRNodeData::NewEnum(Box::new(arg), typ)
            }
            TypeData::ARRAY { .. } => {
                let pars = self.typecheck_variadic(pos, &args, &vec![], TypeData::INT, 1)?;
                if pars.len() == 1 {
                    IRNodeData::NewArray(typ, Some(Box::new(pars[0].clone())))
                } else {
                    IRNodeData::NewArray(typ, None)
                }
            }
            TypeData::STRUCT(fields) => {
                if args.len() != fields.len() {
                    return Err(expectedargcount(fields.len(), pos, args.len()));
                }
                self.scopes
                    .push(Scope::new(ScopeKind::Struct(fields), range));
                self.stack.push(self.scopes.len() - 1);
                let args = self.typecheck_variadic(pos, &args, &vec![], TypeData::VOID, 0)?;
                self.stack.pop().unwrap();
                for arg in args.iter() {
                    match &arg.data {
                        IRNodeData::SetStruct { .. } => {}
                        IRNodeData::Invalid => {}
                        _ => {
                            return Err(IRError::InvalidArgument {
                                expected: TypeData::VOID,
                                got: arg.clone(),
                            })
                        }
                    }
                }

                IRNodeData::NewStruct(typ, args)
            }
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: t,
                })
            }
        };
        Ok(IRNode::new(dat, range, pos))
    }

    pub fn build_arr(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut pars = Vec::new();
        for arg in args {
            let v = self.build_node(arg);
            if v.typ(self).data != TypeData::VOID {
                pars.push(v);
            }
        }
        if pars.len() < 1 {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: 1,
                got: pars.len(),
            });
        }

        let t = pars[0].typ(self).data;
        for v in pars.iter().skip(1) {
            let typ = v.typ(self).data;
            if typ != t && typ != TypeData::INVALID {
                self.save_error(IRError::InvalidArgument {
                    expected: t.clone(),
                    got: v.clone(),
                });
            }
        }

        Ok(IRNode::new(
            IRNodeData::NewArrayLiteral(Type::from(TypeData::ARRAY(Box::new(Type::from(t)))), pars),
            range,
            pos,
        ))
    }

    pub fn build_get(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut args = args.clone();
        if args.len() < 1 {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: 1,
                got: 0,
            });
        }
        let t = self.build_node(&args[0]);
        let typ = match &t.data {
            IRNodeData::Type(t) => t.clone(),
            _ => Type::from(TypeData::INVALID),
        };
        args.remove(0);

        let dat = match typ.data.concrete(self) {
            TypeData::INVALID => IRNodeData::Invalid,
            TypeData::ARRAY { .. } => {
                let pars = self.typecheck(pos, &args, &vec![TypeData::INT])?;
                IRNodeData::GetArr {
                    arr: Box::new(t),
                    ind: Box::new(pars[0].clone()),
                }
            }
            // TODO: Structs, enums, etc.
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: t,
                })
            }
        };
        Ok(IRNode::new(dat, range, pos))
    }

    pub fn build_structop(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let ind = *self.stack.last().unwrap();
        let fields = match &self.scopes[ind].kind {
            ScopeKind::Struct(fields) => fields.clone(),
            _ => return Err(IRError::StructOpOutsideDef(pos)),
        };
        let args = typecheck_ast(
            pos,
            args,
            &vec![ASTNodeDataType::Field, ASTNodeDataType::Any],
        )?;
        let name = match &args[0].data {
            ASTNodeData::Field(name) => name,
            _ => unreachable!(),
        };
        let field = fields.iter().position(|x| x.name == *name);
        if !field.is_some() {
            return Err(IRError::UnknownField {
                pos,
                name: name.clone(),
            });
        }
        let arg = self.build_node(&args[1]);
        if arg.typ(self).data != TypeData::INVALID
            && arg.typ(self).data != fields[field.unwrap()].typ.data
        {
            return Err(IRError::InvalidArgument {
                expected: fields[field.unwrap()].typ.data.clone(),
                got: arg,
            });
        }
        Ok(IRNode::new(
            IRNodeData::StructOp {
                field: name.clone(),
                val: Box::new(arg),
            },
            range,
            pos,
        ))
    }

    pub fn build_set(
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
        let str = self.build_node(&args[0]);
        let typ = match str.typ(self).data.concrete(self) {
            TypeData::STRUCT(v) => v,
            TypeData::INVALID => Vec::new(),
            _ => {
                return Err(IRError::InvalidArgument {
                    expected: TypeData::STRUCT(Vec::new()),
                    got: str,
                })
            }
        };
        self.scopes.push(Scope::new(ScopeKind::Struct(typ), range));
        self.stack.push(self.scopes.len() - 1);
        let mut pars = Vec::new();
        for v in args.iter().skip(1) {
            let val = self.build_node(v);
            match val.data {
                IRNodeData::Invalid => continue,
                IRNodeData::StructOp { .. } => {}
                _ => {
                    return Err(IRError::InvalidArgument {
                        expected: TypeData::FIELD,
                        got: val,
                    })
                }
            }
            pars.push(val);
        }
        self.stack.pop().unwrap();
        Ok(IRNode::new(
            IRNodeData::SetStruct {
                strct: Box::new(str),
                vals: pars,
            },
            range,
            pos,
        ))
    }
}
