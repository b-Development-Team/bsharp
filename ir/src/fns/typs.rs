use super::*;

impl IR {
    pub fn build_array(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = self.typecheck(pos, &args, &vec![TypeData::TYPE])?;

        let dat = match &args[0].data {
            IRNodeData::Type(t) => {
                IRNodeData::Type(Type::from(TypeData::ARRAY(Box::new(t.clone()))))
            }
            IRNodeData::Invalid => IRNodeData::Type(Type::from(TypeData::ARRAY(Box::new(
                Type::from(TypeData::INVALID),
            )))),
            _ => unreachable!(),
        };
        Ok(IRNode::new(dat, range, pos))
    }

    pub fn build_struct(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let pars = self.typecheck_variadic(pos, args, &Vec::new(), TypeData::FIELD, 0)?;

        let mut fields = Vec::new();
        for field in pars {
            match field.data {
                IRNodeData::Field { name, typ } => fields.push(Field { name, typ }),
                _ => continue,
            }
        }

        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::STRUCT(fields))),
            range,
            pos,
        ))
    }

    pub fn build_tuple(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let pars = self.typecheck_variadic(pos, args, &Vec::new(), TypeData::TYPE, 0)?;
        let mut vals = Vec::with_capacity(pars.len());
        for par in pars {
            match par.data {
                IRNodeData::Type(t) => vals.push(t),
                IRNodeData::Invalid => vals.push(Type::from(TypeData::INVALID)),
                _ => unreachable!(),
            }
        }
        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::TUPLE(vals))),
            range,
            pos,
        ))
    }

    pub fn build_enum(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let pars = self.typecheck_variadic(pos, args, &Vec::new(), TypeData::TYPE, 0)?;
        let mut vals = Vec::with_capacity(pars.len());
        for par in pars {
            match par.data {
                IRNodeData::Type(t) => vals.push(t),
                IRNodeData::Invalid => vals.push(Type::from(TypeData::INVALID)),
                _ => unreachable!(),
            }
        }
        Ok(IRNode::new(
            IRNodeData::Type(Type::from(TypeData::ENUM(vals))),
            range,
            pos,
        ))
    }

    pub fn build_prim(
        &mut self,
        name: &String,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        match self.typecheck(pos, args, &Vec::new()) {
            Ok(_) => (),
            Err(e) => self.save_error(e),
        };
        let t = match name.as_str() {
            "CHAR" => TypeData::CHAR,
            "INT" => TypeData::INT,
            "FLOAT" => TypeData::FLOAT,
            "BOOL" => TypeData::BOOL,
            "BX" => TypeData::BOX,
            _ => unreachable!(),
        };
        Ok(IRNode::new(IRNodeData::Type(Type::from(t)), range, pos))
    }

    pub fn build_field(
        &mut self,
        pos: Pos,
        range: Pos,
        args: &Vec<ASTNode>,
    ) -> Result<IRNode, IRError> {
        let args = typecheck_ast(
            pos,
            args,
            &vec![ASTNodeDataType::Field, ASTNodeDataType::Any],
        )?;
        let typ = self.build_node(&args[1]);
        typ.typ(self).expect(typ.pos, &Type::from(TypeData::TYPE))?;

        let name = match &args[0].data {
            ASTNodeData::Field(value) => value.clone(),
            _ => unreachable!(),
        };
        let t = match typ.data {
            IRNodeData::Type(t) => t,
            IRNodeData::Invalid => Type::from(TypeData::INVALID),
            _ => unreachable!(),
        };

        Ok(IRNode::new(IRNodeData::Field { name, typ: t }, range, pos))
    }

    pub fn build_typeval(&mut self, pos: Pos, val: String) -> Result<IRNode, IRError> {
        if let Some(ind) = self.typemap.get(&val) {
            Ok(IRNode::new(
                IRNodeData::Type(Type::from(TypeData::DEF(*ind))),
                pos,
                pos,
            ))
        } else {
            Err(IRError::UnknownType { pos, name: val })
        }
    }
}
