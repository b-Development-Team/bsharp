use super::*;

mod composite;
mod control;
mod fns;
mod ops;
mod typs;
mod vars;

impl IR {
    pub fn build_print(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = self.typecheck(
            pos,
            &args,
            &vec![TypeData::ARRAY(Box::new(Type::from(TypeData::CHAR)))],
        )?;
        Ok(IRNode::new(
            IRNodeData::Print(Box::new(args[0].clone())),
            range,
            pos,
        ))
    }
}
