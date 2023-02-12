use super::*;

pub fn typecheck_ast(
    pos: Pos,
    params: &Vec<ASTNode>,
    typs: &Vec<ASTNodeDataType>,
) -> Result<(), IRError> {
    if params.len() != typs.len() {
        return Err(IRError::InvalidArgumentCount {
            pos,
            expected: typs.len(),
            got: params.len(),
        });
    }
    for (i, typ) in typs.iter().enumerate() {
        if params[i].data.typ() != *typ {
            return Err(IRError::InvalidASTArgument {
                expected: *typ,
                got: params[i].clone(),
            });
        }
    }
    Ok(())
}
