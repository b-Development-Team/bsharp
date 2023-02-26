use super::*;

impl IR {
    pub fn build_typ(&mut self, id: usize) {
        let ast = self.types[id].ast.clone();
        if let Some(ast) = &ast {
            // Add scope
            self.scopes
                .push(Scope::new(ScopeKind::Type, self.types[id].pos));
            self.stack.push(self.scopes.len() - 1);

            // Build
            let res = self.build_node(ast);
            match res.typ().expect(res.pos, &Type::from(TypeData::TYPE)) {
                Err(err) => {
                    self.save_error(err);
                    return;
                }
                Ok(_) => {}
            };
            self.types[id].ast = None;
            self.types[id].typ = match res.data {
                IRNodeData::Type(v) => v,
                _ => unreachable!(),
            };
            self.types[id].scope = self.stack.pop().unwrap();
        }
    }

    pub fn build_fn(&mut self, id: usize) {
        // Build params
        if let Some(ast) = &self.funcs[id].params_ast {
            let pars = match &ast.data {
                ASTNodeData::Block(v) => v,
                _ => unreachable!(),
            }
            .clone();

            let mut generics = Vec::new();
            let mut params = Vec::new();
            let mut ret = None;
            let mut ret_def = Some(ast.pos);

            // Add type scope
            self.scopes.push(Scope::new(ScopeKind::Type, ast.pos));
            self.stack.push(self.scopes.len() - 1);

            for par in pars.iter() {
                let n = self.build_node(par);
                match n.data {
                    IRNodeData::Generic(v) => {
                        if params.len() > 0 {
                            self.save_error(IRError::InvalidArgument {
                                expected: TypeData::PARAM,
                                got: n.clone(),
                            });
                        }
                        generics.push(v);
                    }
                    IRNodeData::Param(ind) => {
                        params.push(ind);
                    }
                    IRNodeData::Returns(ref t) => {
                        if ret.is_some() {
                            self.save_error(IRError::InvalidArgument {
                                expected: TypeData::PARAM,
                                got: n.clone(),
                            });
                        }
                        ret = Some(t.clone());
                        ret_def = Some(n.pos);
                    }
                    IRNodeData::Invalid | IRNodeData::Void => {}
                    _ => self.save_error(IRError::InvalidArgument {
                        expected: TypeData::PARAM,
                        got: n,
                    }),
                }
            }
            self.stack.pop().unwrap();

            self.funcs[id].params_ast = None;
            self.funcs[id].params = params;
            self.funcs[id].generic_params = generics;
            if let Some(_) = ret {
                self.funcs[id].ret_typ = ret.unwrap();
            } else {
                self.funcs[id].ret_typ = Type::from(TypeData::VOID);
            }
            self.funcs[id].ret_typ_definition = ret_def.unwrap();
        }

        if let Some(ast) = &self.funcs[id].body_ast {
            // Add params
            let mut sc = Scope::new(ScopeKind::Function(id), ast.pos);
            for par in self.funcs[id].params.iter() {
                self.variables[*par].scope = self.scopes.len();
                sc.vars.insert(self.variables[*par].name.clone(), *par);
            }
            for gen in self.funcs[id].generic_params.iter() {
                sc.types.insert(self.types[*gen].name.clone(), *gen);
            }

            // Add scope
            self.scopes.push(sc);
            self.stack.push(self.scopes.len() - 1);

            // Build
            let res = self.build_node(&ast.clone());
            self.funcs[id].body_ast = None;
            self.funcs[id].body = res;
            self.funcs[id].scope = self.stack.pop().unwrap();
        }
    }

    pub fn build_node(&mut self, v: &ASTNode) -> IRNode {
        match match &v.data {
            ASTNodeData::Comment(_) => Ok(IRNode::new(IRNodeData::Void, v.pos, v.pos)),
            ASTNodeData::Stmt {
                name,
                name_pos,
                args,
            } => match name.as_str() {
                "ARRAY" => self.build_array(*name_pos, v.pos, args),
                "STRUCT" => self.build_struct(*name_pos, v.pos, args),
                "TUPLE" => self.build_tuple(*name_pos, v.pos, args),
                "ENUM" => self.build_enum(*name_pos, v.pos, args),
                "INTERFACE" => self.build_interface(*name_pos, v.pos, args),
                "FIELD" => self.build_field(*name_pos, v.pos, args),
                "GENERIC" => self.build_generic(*name_pos, v.pos, args),
                "CHAR" | "INT" | "FLOAT" | "BOOL" => self.build_prim(name, *name_pos, v.pos, args),
                "PARAM" => self.build_param(*name_pos, v.pos, args),
                "RETURNS" => self.build_returns(*name_pos, v.pos, args),
                "RETURN" => self.build_return(*name_pos, v.pos, args),
                "+" => self.build_mathop(*name_pos, v.pos, args, MathOperator::ADD),
                "-" => self.build_mathop(*name_pos, v.pos, args, MathOperator::SUBTRACT),
                "/" => self.build_mathop(*name_pos, v.pos, args, MathOperator::DIVIDE),
                "*" => self.build_mathop(*name_pos, v.pos, args, MathOperator::MULTIPLY),
                "%" => self.build_mathop(*name_pos, v.pos, args, MathOperator::MODULO),
                "DEFINE" => self.build_define(*name_pos, v.pos, args),
                _ => Err(IRError::UnknownStmt {
                    pos: *name_pos,
                    name: name.clone(),
                }),
            },
            ASTNodeData::Type(name) => self.build_typeval(v.pos, name.clone()),
            ASTNodeData::Variable(name) => self.build_var(v.pos, name.clone()),
            ASTNodeData::Block(pars) => self.build_block(v.pos, pars),
            ASTNodeData::Integer(val) => Ok(IRNode::new(IRNodeData::Int(*val), v.pos, v.pos)),
            _ => Err(IRError::UnexpectedNode(v.clone())),
        } {
            Ok(res) => res,
            Err(err) => {
                self.save_error(err);
                IRNode::invalid()
            }
        }
    }
}
