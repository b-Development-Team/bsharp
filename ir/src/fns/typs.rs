use super::*;

impl IR {
    pub fn build_array(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let dat = if args.len() != 1 {
            self.save_error(IRError::InvalidArgumentCount {
                pos,
                expected: 1,
                got: args.len(),
            });
            IRNodeData::Type(Type::from(TypeData::ARRAY(
                None,
                Box::new(Type::from(TypeData::INVALID)),
            )))
        } else {
            // TODO: Support comments in ARRAY
            let v = self.build_node(&args[0]);
            match v.data {
                IRNodeData::Type(t) => {
                    IRNodeData::Type(Type::from(TypeData::ARRAY(None, Box::new(t))))
                }
                IRNodeData::Generic { typ, name } => IRNodeData::Type(Type::from(TypeData::ARRAY(
                    Some(Box::new(Generic {
                        typ: typ.clone(),
                        name: name.clone(),
                    })),
                    Box::new(Type::new(typ.data, Some(name))),
                ))),
                IRNodeData::Invalid => IRNodeData::Type(Type::from(TypeData::ARRAY(
                    None,
                    Box::new(Type::from(TypeData::INVALID)),
                ))),
                _ => {
                    self.save_error(IRError::InvalidArgument {
                        expected: TypeData::TYPE,
                        got: v,
                    });
                    IRNodeData::Type(Type::from(TypeData::ARRAY(
                        None,
                        Box::new(Type::from(TypeData::INVALID)),
                    )))
                }
            }
        };

        Ok(IRNode::new(dat, range, pos))
    }

    pub fn build_struct(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let pars = self.typecheck_variadic(pos, args, &Vec::new(), TypeData::FIELD, 0)?;
        let mut fields = Vec::new();
        let mut params = Vec::new();
        for field in pars {
            match field.data {
                IRNodeData::Field { name, typ } => fields.push(Field { name, typ }),
                IRNodeData::Generic { name, typ } => params.push(Generic { name, typ }),
                _ => continue,
            }
        }

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::STRUCT { params, fields })),
            range,
            pos,
        ))
    }

    pub fn build_tuple(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut fields = Vec::new();
        let mut params = Vec::new();

        for par in args.iter() {
            let v = self.build_node(par);
            match v.data {
                IRNodeData::Generic { name, typ } => params.push(Generic { name, typ }),
                IRNodeData::Type(t) => fields.push(t),
                IRNodeData::Invalid | IRNodeData::Void => {}
                _ => self.save_error(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: v,
                }),
            }
        }

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::TUPLE {
                params,
                body: fields,
            })),
            range,
            pos,
        ))
    }

    pub fn build_enum(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut fields = Vec::new();
        let mut params = Vec::new();

        for par in args.iter() {
            let v = self.build_node(par);
            match v.data {
                IRNodeData::Generic { name, typ } => params.push(Generic { name, typ }),
                IRNodeData::Type(t) => fields.push(t),
                IRNodeData::Invalid | IRNodeData::Void => {}
                _ => self.save_error(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: v,
                }),
            }
        }

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::ENUM {
                params,
                body: fields,
            })),
            range,
            pos,
        ))
    }

    pub fn build_prim(
        &mut self,
        name: &String,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        match self.typecheck(pos, args, &Vec::new()) {
            Ok(_) => (),
            Err(e) => self.save_error(e),
        };
        let t = match name.as_str() {
            "CHAR" => TypeData::CHAR,
            "INT" => TypeData::INT,
            "FLOAT" => TypeData::FLOAT,
            "BOOL" => TypeData::BOOL,
            _ => unreachable!(),
        };
        Ok(IRNode::new(IRNodeData::Type(Type::from(t)), range, pos))
    }

    pub fn build_field(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(
            pos,
            args,
            &vec![ASTNodeDataType::Field, ASTNodeDataType::Any],
        )?;
        let typ = self.build_node(&args[1]);
        typ.typ().expect(typ.pos, &Type::from(TypeData::TYPE))?;

        let name = match &args[0].data {
            ASTNodeData::Field(value) => value.clone(),
            _ => unreachable!(),
        };
        let t = match typ.data {
            IRNodeData::Type(t) => t,
            IRNodeData::Invalid => Type::from(TypeData::INVALID),
            _ => unreachable!(),
        };

        Ok(IRNode::new(IRNodeData::Field { name, typ: t }, range, pos))
    }

    pub fn build_typeval(&mut self, pos: Pos, val: String) -> Result<IRNode, IRError> {
        for scope in self.stack.iter().rev() {
            let typ = match self.scopes[*scope].types.get(&val) {
                Some(t) => *t,
                None => continue,
            };
            self.build_typ(typ);
            return Ok(IRNode::new(
                IRNodeData::Type(self.types[typ].typ.clone()),
                pos,
                pos,
            ));
        }
        Err(IRError::UnknownType { pos, name: val })
    }

    pub fn build_interface(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut fields = Vec::new();
        let mut params = Vec::new();

        for par in args.iter() {
            let v = self.build_node(par);
            match v.data {
                IRNodeData::Generic { name, typ } => params.push(Generic { name, typ }),
                IRNodeData::Type(ref t) => match &t.data {
                    TypeData::INT
                    | TypeData::FLOAT
                    | TypeData::CHAR
                    | TypeData::BOX
                    | TypeData::BOOL
                    | TypeData::ARRAY(_, _)
                    | TypeData::STRUCT { .. }
                    | TypeData::TUPLE { .. }
                    | TypeData::ENUM { .. } => fields.push(t.clone()),
                    TypeData::INVALID => {}
                    _ => self.save_error(IRError::InvalidArgument {
                        expected: TypeData::TYPE,
                        got: v,
                    }),
                },
                IRNodeData::Invalid | IRNodeData::Void => {}
                _ => self.save_error(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: v,
                }),
            }
        }

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::INTERFACE {
                params,
                body: fields,
            })),
            range,
            pos,
        ))
    }

    pub fn build_generic(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(
            pos,
            args,
            &vec![ASTNodeDataType::Type, ASTNodeDataType::Any],
        )?;

        // Get type
        let typval = self.build_node(&args[1]);
        typval
            .typ()
            .expect(typval.pos, &Type::from(TypeData::TYPE))?;
        let typ = match &typval.data {
            IRNodeData::Type(v) => v.clone(),
            _ => unreachable!(),
        };
        match typ.data {
            TypeData::INTERFACE { .. } | TypeData::INVALID => {}
            _ => self.save_error(IRError::InvalidArgument {
                expected: TypeData::INTERFACE {
                    params: Vec::new(),
                    body: Vec::new(),
                },
                got: typval,
            }),
        }

        // Get name
        let name = match &args[0].data {
            ASTNodeData::Type(value) => value.clone(),
            _ => unreachable!(),
        };

        // Add to stack
        let scope = *self.stack.last().unwrap();
        self.types.push(TypeDef {
            scope,
            name: name.clone(),
            pos,
            ast: None,
            typ: typ.clone(),
        });
        self.scopes[scope]
            .types
            .insert(name.clone(), self.types.len() - 1);

        Ok(IRNode::new(IRNodeData::Generic { name, typ }, range, pos))
    }
}
