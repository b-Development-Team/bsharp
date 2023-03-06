use super::*;

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
            TypeData::ARRAY { .. } => {
                let pars = self.typecheck_variadic(pos, &args, &vec![], TypeData::INT, 1)?;
                if pars.len() == 1 {
                    IRNodeData::NewArray(typ, Some(Box::new(pars[0].clone())))
                } else {
                    IRNodeData::NewArray(typ, None)
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
}
