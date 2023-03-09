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
                        let args = match typecheck_ast(
                            *name_pos,
                            args,
                            &vec![ASTNodeDataType::Type, ASTNodeDataType::Stmt],
                        ) {
                            Ok(v) => v,
                            Err(err) => {
                                self.save_error(err);
                                continue;
                            }
                        };

                        let name = match &args[0].data {
                            ASTNodeData::Type(name) => name.clone(),
                            _ => unreachable!(),
                        };

                        if let Some(i) = self.typemap.get(&name) {
                            self.save_error(IRError::DuplicateType(*name_pos, *i));
                            continue;
                        }
                        self.typemap.insert(name.clone(), self.types.len());
                        self.types.push(TypeDef {
                            name,
                            pos: *name_pos,
                            ast: Some(args[1].clone()),
                            typ: Type::void(),
                        });
                    }
                    "FUNC" => {
                        let args = match typecheck_ast(
                            *name_pos,
                            args,
                            &vec![
                                ASTNodeDataType::Function,
                                ASTNodeDataType::Block,
                                ASTNodeDataType::Block,
                            ],
                        ) {
                            Ok(v) => v,
                            Err(err) => {
                                self.save_error(err);
                                continue;
                            }
                        };

                        let name = match &args[0].data {
                            ASTNodeData::Function(name) => name.clone(),
                            _ => unreachable!(),
                        };

                        // Check if exists
                        if let Some(i) = self.funcs.iter().position(|v| v.name == name) {
                            self.save_error(IRError::DuplicateFunction(*name_pos, i));
                            continue;
                        }

                        self.funcs.push(Function {
                            definition: *name_pos,
                            name,

                            params: Vec::new(),
                            ret_typ: Type::void(),
                            ret_typ_definition: Pos::default(),
                            body: IRNode::void(),

                            params_ast: Some(args[1].clone()),
                            body_ast: Some(args[2].clone()),
                            scope: 0,
                        })
                    }
                    "IMPORT" => {
                        let args =
                            match typecheck_ast(*name_pos, args, &vec![ASTNodeDataType::String]) {
                                Ok(v) => v,
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
