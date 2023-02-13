use super::*;

impl IR {
    pub fn build_typ(&mut self, id: usize) {
        let ast = self.types[id].ast.clone();
        if let Some(ast) = &ast {
            // Add scope
            self.scopes
                .push(Scope::new(ScopeKind::Type, self.types[id].pos));

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
            self.types[id].typ = res.typ();
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
                "CHAR" | "INT" | "FLOAT" | "BOOL" => self.build_prim(name, *name_pos, v.pos, args),
                _ => Err(IRError::UnknownFunction {
                    pos: *name_pos,
                    name: name.clone(),
                }),
            },
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
