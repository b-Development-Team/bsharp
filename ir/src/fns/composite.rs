use super::*;

impl IR {
    fn check_array_par(&mut self, par: &IRNode) -> Result<Type, IRError> {
        match par.typ().data.concrete(self) {
            TypeData::ARRAY { body, .. } => Ok(*body),
            TypeData::INVALID => Ok(Type::from(TypeData::INVALID)),
            TypeData::INTERFACE { body, .. } => {
                for t in body {
                    match t.data.concrete(self) {
                        TypeData::ARRAY { body, .. } => return Ok(*body),
                        _ => {}
                    }
                }

                Err(IRError::InvalidArgument {
                    expected: TypeData::ARRAY {
                        param: None,
                        body: Box::new(Type::void()),
                        scope: 0,
                    },
                    got: par.clone(),
                })
            }
            _ => Err(IRError::InvalidArgument {
                expected: TypeData::ARRAY {
                    param: None,
                    body: Box::new(Type::void()),
                    scope: 0,
                },
                got: par.clone(),
            }),
        }
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

    pub fn build_append(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let val = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any, ASTNodeDataType::Any])?;
        let arr = self.build_node(&val[0]);
        let body_typ = self.check_array_par(&arr)?;
        let val = self.build_node(&val[1]);
        if val.typ().data.concrete(self) != body_typ.data.concrete(self)
            && val.typ().data.concrete(self) != TypeData::INVALID
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
}
