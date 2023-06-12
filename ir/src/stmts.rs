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
            match res.typ(self).expect(res.pos, &Type::from(TypeData::TYPE)) {
                Err(err) => {
                    self.save_error(err);
                    return;
                }
                Ok(_) => {}
            };
            self.types[id].ast = None;
            self.types[id].typ = match res.data {
                IRNodeData::Type(v) => v,
                IRNodeData::Invalid => Type::from(TypeData::INVALID),
                _ => unreachable!(),
            };
        }
    }

    pub fn build_fn(&mut self, id: usize, sig: bool) {
        // Build params
        if let Some(ast) = &self.funcs[id].params_ast {
            let pars = match &ast.data {
                ASTNodeData::Block(v) => v,
                _ => unreachable!(),
            }
            .clone();

            let mut params = Vec::new();
            let mut ret = None;
            let mut ret_def = Some(ast.pos);

            // Add type scope
            self.scopes.push(Scope::new(ScopeKind::Type, ast.pos));
            self.stack.push(self.scopes.len() - 1);

            for par in pars.iter() {
                let n = self.build_node(par);
                match n.data {
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
            if let Some(_) = ret {
                self.funcs[id].ret_typ = ret.unwrap();
            } else {
                self.funcs[id].ret_typ = Type::from(TypeData::VOID);
            }
            self.funcs[id].ret_typ_definition = ret_def.unwrap();
        }

        if !sig {
            if let Some(ast) = &self.funcs[id].body_ast {
                // Add params
                let mut sc = Scope::new(ScopeKind::Function(id), ast.pos);
                for par in self.funcs[id].params.iter() {
                    self.variables[*par].scope = self.scopes.len();
                    sc.vars.insert(self.variables[*par].name.clone(), *par);
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
                "FIELD" => self.build_field(*name_pos, v.pos, args),
                "CHAR" | "INT" | "FLOAT" | "BOOL" | "BX" => {
                    self.build_prim(name, *name_pos, v.pos, args)
                }
                "PARAM" => self.build_param(*name_pos, v.pos, args),
                "RETURNS" => self.build_returns(*name_pos, v.pos, args),
                "RETURN" => self.build_return(*name_pos, v.pos, args),
                "ARR" => self.build_arr(*name_pos, v.pos, args),
                "+" => self.build_mathop(*name_pos, v.pos, args, MathOperator::ADD),
                "-" => self.build_mathop(*name_pos, v.pos, args, MathOperator::SUBTRACT),
                "/" => self.build_mathop(*name_pos, v.pos, args, MathOperator::DIVIDE),
                "*" => self.build_mathop(*name_pos, v.pos, args, MathOperator::MULTIPLY),
                "%" => self.build_mathop(*name_pos, v.pos, args, MathOperator::MODULO),
                "XOR" => self.build_mathop(*name_pos, v.pos, args, MathOperator::XOR),
                "SHIFT" => self.build_mathop(*name_pos, v.pos, args, MathOperator::SHIFT),
                "BOR" => self.build_mathop(*name_pos, v.pos, args, MathOperator::BOR),
                ">" => self.build_compop(*name_pos, v.pos, args, ComparisonOperator::GREATER),
                ">=" => self.build_compop(*name_pos, v.pos, args, ComparisonOperator::GREATEREQUAL),
                "<" => self.build_compop(*name_pos, v.pos, args, ComparisonOperator::LESS),
                "<=" => self.build_compop(*name_pos, v.pos, args, ComparisonOperator::LESSEQUAL),
                "=" => self.build_compop(*name_pos, v.pos, args, ComparisonOperator::EQUAL),
                "NEQ" => self.build_compop(*name_pos, v.pos, args, ComparisonOperator::NOTEQUAL),
                "NOT" => self.build_boolop(*name_pos, v.pos, args, BooleanOperator::NOT),
                "&" => self.build_boolop(*name_pos, v.pos, args, BooleanOperator::AND),
                "|" => self.build_boolop(*name_pos, v.pos, args, BooleanOperator::OR),
                "DEFINE" => self.build_define(*name_pos, v.pos, args),
                "WHILE" => self.build_while(*name_pos, v.pos, args),
                "LEN" => self.build_len(*name_pos, v.pos, args),
                "APPEND" => self.build_append(*name_pos, v.pos, args),
                "GET" => self.build_get(*name_pos, v.pos, args),
                "NEW" => self.build_new(*name_pos, v.pos, args),
                "PRINT" => self.build_print(*name_pos, v.pos, args),
                "IF" => self.build_if(*name_pos, v.pos, args),
                "BOX" => self.build_box(*name_pos, v.pos, args),
                "PEEK" => self.build_peek(*name_pos, v.pos, args),
                "UNBOX" => self.build_unbox(*name_pos, v.pos, args),
                ":" => self.build_structop(*name_pos, v.pos, args),
                "SET" => self.build_set(*name_pos, v.pos, args),
                "MATCH" => self.build_match(*name_pos, v.pos, args),
                "CASE" => self.build_case(*name_pos, v.pos, args),
                "TOI" => self.build_cast_int(*name_pos, v.pos, args),
                "TOF" => self.build_cast_float(*name_pos, v.pos, args),
                "TOC" => self.build_cast_char(*name_pos, v.pos, args),
                _ => Err(IRError::UnknownStmt {
                    pos: *name_pos,
                    name: name.clone(),
                }),
            },
            ASTNodeData::Type(name) => self.build_typeval(v.pos, name.clone()),
            ASTNodeData::Variable(name) => self.build_var(v.pos, name.clone()),
            ASTNodeData::Block(pars) => self.build_block(v.pos, pars),
            ASTNodeData::Integer(val) => Ok(IRNode::new(IRNodeData::Int(*val), v.pos, v.pos)),
            ASTNodeData::Float(val) => Ok(IRNode::new(IRNodeData::Float(*val), v.pos, v.pos)),
            ASTNodeData::Char(val) => Ok(IRNode::new(IRNodeData::Char(*val as u8), v.pos, v.pos)),
            ASTNodeData::String(val) => Ok(self.build_string(v.pos, val)),
            _ => Err(IRError::UnexpectedNode(v.clone())),
        } {
            Ok(res) => res,
            Err(err) => {
                self.save_error(err);
                IRNode::invalid(v.pos)
            }
        }
    }
}
