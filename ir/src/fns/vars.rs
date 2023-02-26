use super::*;

impl IR {
    pub fn build_var(&mut self, pos: Pos, val: String) -> Result<IRNode, IRError> {
        for scope in self.stack.iter().rev() {
            let var = match self.scopes[*scope].vars.get(&val) {
                Some(t) => *t,
                None => {
                    continue;
                }
            };
            return Ok(IRNode::new(
                IRNodeData::Variable(var, self.variables[var].typ.clone()),
                pos,
                pos,
            ));
        }

        Err(IRError::UnknownVariable { pos, name: val })
    }

    pub fn build_define(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(
            pos,
            args,
            &vec![ASTNodeDataType::Variable, ASTNodeDataType::Any],
        )?;
        let name = match &args[0].data {
            ASTNodeData::Variable(v) => v.clone(),
            _ => unreachable!(),
        };
        let value = self.build_node(&args[1]);

        // Check if name is already on stack
        if let Some(ind) = self.scopes[*self.stack.last().unwrap()].vars.get(&name) {
            if self.variables[*ind].typ.data == value.typ().data {
                return Ok(IRNode::new(
                    IRNodeData::Define {
                        var: *ind,
                        val: Box::new(value),
                        edit: true,
                    },
                    range,
                    pos,
                ));
            } else {
                return Err(IRError::DuplicateVariable(pos, *ind));
            }
        }

        self.variables.push(Variable {
            name: name.clone(),
            typ: value.typ().clone(),
            scope: *self.stack.last().unwrap(),
            definition: pos,
        });
        self.scopes[*self.stack.last().unwrap()]
            .vars
            .insert(name.clone(), self.variables.len() - 1);

        Ok(IRNode::new(
            IRNodeData::Define {
                var: self.variables.len() - 1,
                val: Box::new(value),
                edit: false,
            },
            range,
            pos,
        ))
    }
}
