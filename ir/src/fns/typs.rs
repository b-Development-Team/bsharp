use std::collections::HashMap;

use super::*;

impl IR {
    fn push_scope(&mut self, pos: Pos) {
        self.scopes.push(Scope {
            kind: ScopeKind::Type,
            vars: HashMap::new(),
            types: HashMap::new(),
            pos,
        });
        self.stack.push(self.scopes.len() - 1);
    }

    pub fn build_array(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(pos, args, &vec![ASTNodeDataType::Any])?;
        self.push_scope(range);
        let v = self.build_node(&args[0]);
        let sc = self.stack.pop().unwrap();

        let dat = match v.data {
            IRNodeData::Type(t) => IRNodeData::Type(Type::from(TypeData::ARRAY {
                param: None,
                body: Box::new(t),
                scope: sc,
            })),
            IRNodeData::Generic(v) => IRNodeData::Type(Type::from(TypeData::ARRAY {
                param: Some(v),
                body: Box::new(Type::from(TypeData::DEF(v))),
                scope: sc,
            })),
            IRNodeData::Invalid => IRNodeData::Type(Type::from(TypeData::ARRAY {
                param: None,
                body: Box::new(Type::from(TypeData::INVALID)),
                scope: sc,
            })),
            _ => {
                self.save_error(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: v,
                });
                IRNodeData::Type(Type::from(TypeData::ARRAY {
                    param: None,
                    body: Box::new(Type::from(TypeData::INVALID)),
                    scope: sc,
                }))
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
        self.push_scope(range);
        let pars = self.typecheck_variadic(pos, args, &Vec::new(), TypeData::FIELD, 0)?;
        let scope = self.stack.pop().unwrap();

        let mut fields = Vec::new();
        let mut params = Vec::new();
        for field in pars {
            match field.data {
                IRNodeData::Field { name, typ } => fields.push(Field { name, typ }),
                IRNodeData::Generic(v) => params.push(v),
                _ => continue,
            }
        }

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::STRUCT {
                params,
                fields,
                scope,
            })),
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

        self.push_scope(range);
        for par in args.iter() {
            let v = self.build_node(par);
            match v.data {
                IRNodeData::Generic(v) => params.push(v),
                IRNodeData::Type(t) => fields.push(t),
                IRNodeData::Invalid | IRNodeData::Void => {}
                _ => self.save_error(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: v,
                }),
            }
        }
        let scope = self.stack.pop().unwrap();

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::TUPLE {
                params,
                body: fields,
                scope,
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

        self.push_scope(range);
        for par in args.iter() {
            let v = self.build_node(par);
            match v.data {
                IRNodeData::Generic(v) => params.push(v),
                IRNodeData::Type(t) => fields.push(t),
                IRNodeData::Invalid | IRNodeData::Void => {}
                _ => self.save_error(IRError::InvalidArgument {
                    expected: TypeData::TYPE,
                    got: v,
                }),
            }
        }
        let scope = self.stack.pop().unwrap();

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::ENUM {
                params,
                body: fields,
                scope,
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
            if self.scopes[*scope].kind != ScopeKind::Type {
                continue;
            }
            let typ = match self.scopes[*scope].types.get(&val) {
                Some(t) => *t,
                None => break,
            };
            self.build_typ(typ);
            return Ok(IRNode::new(
                IRNodeData::Type(Type::from(TypeData::DEF(typ))),
                pos,
                pos,
            ));
        }
        // Check global scope
        let typ = match self.scopes[0].types.get(&val) {
            Some(t) => *t,
            None => return Err(IRError::UnknownType { pos, name: val }),
        };
        self.build_typ(typ);
        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::DEF(typ))),
            pos,
            pos,
        ))
    }

    pub fn build_interface(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let mut fields = Vec::new();
        let mut params = Vec::new();

        self.push_scope(range);

        for par in args.iter() {
            let v = self.build_node(par);
            match v.data {
                IRNodeData::Generic(v) => params.push(v),
                IRNodeData::Type(ref t) => match &t.data.concrete(self) {
                    TypeData::INT
                    | TypeData::FLOAT
                    | TypeData::CHAR
                    | TypeData::BOX
                    | TypeData::BOOL
                    | TypeData::ARRAY { .. }
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

        let scope = self.stack.pop().unwrap();

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::INTERFACE {
                params,
                body: fields,
                scope,
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
        match typ.data.concrete(self) {
            TypeData::INTERFACE { .. } | TypeData::INVALID => {}
            _ => self.save_error(IRError::InvalidArgument {
                expected: TypeData::INTERFACE {
                    params: Vec::new(),
                    body: Vec::new(),
                    scope: 0,
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

        Ok(IRNode::new(
            IRNodeData::Generic(self.types.len() - 1),
            range,
            pos,
        ))
    }
}
