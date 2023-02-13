use super::*;

impl IR {
    pub fn defpass(&mut self) -> Result<(), IRError> {
        // Go through files
        for i in 0..self.fset.files.len() {
            let f = self.fset.files[i].ast.ast.clone();
            self.defpass_file(f)?;
        }

        Ok(())
    }

    fn defpass_file(&mut self, src: Vec<ASTNode>) -> Result<(), FSetError> {
        for node in src {
            match &node.data {
                ASTNodeData::Comment(_) => {}
                ASTNodeData::Stmt {
                    name,
                    name_pos,
                    args,
                } => match name.as_str() {
                    "TYPE" => {
                        match typecheck_ast(
                            *name_pos,
                            args,
                            &vec![ASTNodeDataType::Type, ASTNodeDataType::Stmt],
                        ) {
                            Ok(_) => {}
                            Err(err) => {
                                self.save_error(err);
                                continue;
                            }
                        };

                        let name = match &args[0].data {
                            ASTNodeData::Type(name) => name.clone(),
                            _ => unreachable!(),
                        };

                        let scopeind = *self.stack.last().unwrap();
                        let scope = &mut self.scopes[scopeind];
                        scope.types.insert(name.clone(), self.types.len());

                        self.types.push(TypeDef {
                            scope: scopeind,
                            name,
                            pos: *name_pos,
                            ast: Some(args[1].clone()),
                            typ: Type::void(),
                        })
                    }
                    "FUNC" => {
                        match typecheck_ast(
                            *name_pos,
                            args,
                            &vec![
                                ASTNodeDataType::Function,
                                ASTNodeDataType::Block,
                                ASTNodeDataType::Block,
                            ],
                        ) {
                            Ok(_) => {}
                            Err(err) => {
                                self.save_error(err);
                                continue;
                            }
                        }

                        let name = match &args[0].data {
                            ASTNodeData::Function(name) => name.clone(),
                            _ => unreachable!(),
                        };

                        self.funcs.push(Function {
                            definition: *name_pos,
                            name,

                            params: Vec::new(),
                            ret_typ: Type::void(),
                            ret_typ_definition: Pos::default(),
                            body: IRNode::void(),

                            params_ast: Some(args[1].clone()),
                            body_ast: Some(args[2].clone()),
                        })
                    }
                    "IMPORT" => {
                        match typecheck_ast(*name_pos, args, &vec![ASTNodeDataType::String]) {
                            Ok(_) => {}
                            Err(err) => {
                                self.save_error(err);
                                continue;
                            }
                        };
                        let name = match &args[0].data {
                            ASTNodeData::String(name) => name.clone(),
                            _ => unreachable!(),
                        };
                        self.fset.import(&name)?;
                    }
                    _ => self.save_error(IRError::UnexpectedNode(node)),
                },
                _ => self.save_error(IRError::UnexpectedNode(node)),
            }
        }

        Ok(())
    }
}
