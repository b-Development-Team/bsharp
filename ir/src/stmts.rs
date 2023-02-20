use super::*;

impl IR {
    pub fn build_typ(&mut self, id: usize) {
        let ast = self.types[id].ast.clone();
        if let Some(ast) = &ast {
            // Add scope
            self.scopes
                .push(Scope::new(ScopeKind::Type, self.types[id].pos));

            // Check if within another type
            let old = if self.stack.len() > 2 {
                let v = self.stack[1];
                self.stack[1] = self.scopes.len() - 1;
                Some(v)
            } else {
                self.stack.push(self.scopes.len() - 1);
                None
            };

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
            self.types[id].typ.name = Some(self.types[id].name.clone());
            let ind = if let Some(v) = old {
                let val = self.stack[1];
                self.stack[1] = v;
                val
            } else {
                self.stack.pop().unwrap()
            };
            self.types[id].scope = ind;
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
                _ => Err(IRError::UnknownStmt {
                    pos: *name_pos,
                    name: name.clone(),
                }),
            },
            ASTNodeData::Type(name) => self.build_typeval(v.pos, name.clone()),
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
