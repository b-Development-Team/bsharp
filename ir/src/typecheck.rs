use super::*;

pub fn typecheck_ast(
    pos: Pos,
    params: &Vec<ASTNode>,
    typs: &Vec<ASTNodeDataType>,
) -> Result<Vec<ASTNode>, IRError> {
    if params.len() != typs.len() {
        return Err(IRError::InvalidArgumentCount {
            pos,
            expected: typs.len(),
            got: params.len(),
        });
    }
    let mut res = Vec::new();
    for (i, typ) in typs.iter().enumerate() {
        if *typ == ASTNodeDataType::Any {
            res.push(params[i].clone());
            continue;
        }
        if params[i].data.typ() == ASTNodeDataType::Comment {
            continue;
        }
        if params[i].data.typ() != *typ {
            return Err(IRError::InvalidASTArgument {
                expected: *typ,
                got: params[i].clone(),
            });
        }
        res.push(params[i].clone());
    }
    Ok(res)
}

impl IR {
    pub fn typecheck(
        &mut self,
        pos: Pos,
        params: &Vec<ASTNode>,
        typs: &Vec<TypeData>,
    ) -> Result<Vec<IRNode>, IRError> {
        if params.len() != typs.len() {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: typs.len(),
                got: params.len(),
            });
        }

        let mut res = Vec::new();
        for (i, node) in params.iter().enumerate() {
            if node.data.typ() == ASTNodeDataType::Comment {
                continue;
            }

            let v = self.build_node(node);
            if v.typ(self).data != typs[i] && v.typ(self).data != TypeData::INVALID {
                self.save_error(IRError::InvalidArgument {
                    got: v.clone(),
                    expected: typs[i].clone(),
                });
            }

            res.push(v);
        }

        if res.len() != typs.len() {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: typs.len(),
                got: res.len(),
            });
        }

        Ok(res)
    }

    pub fn typecheck_variadic(
        &mut self,
        pos: Pos,
        params: &Vec<ASTNode>,
        typs: &Vec<TypeData>,
        end: TypeData,
        maxlen: usize, // 0 for unlimited
    ) -> Result<Vec<IRNode>, IRError> {
        if params.len() < typs.len() || (maxlen > 0 && params.len() > (typs.len() + maxlen)) {
            return Err(IRError::InvalidArgumentCount {
                pos,
                expected: typs.len(),
                got: params.len(),
            });
        }

        let mut res = Vec::new();
        for (i, node) in params.iter().enumerate() {
            let v = self.build_node(node);
            if v.typ(self).data == TypeData::VOID {
                continue;
            }
            if i >= typs.len() {
                if v.typ(self).data != end && v.typ(self).data != TypeData::INVALID {
                    return Err(IRError::InvalidArgument {
                        got: v,
                        expected: end.clone(),
                    });
                }
            } else {
                if v.typ(self).data != typs[i] && v.typ(self).data != TypeData::INVALID {
                    return Err(IRError::InvalidArgument {
                        got: v,
                        expected: typs[i].clone(),
                    });
                }
            }

            res.push(v);
        }
        Ok(res)
    }
}
