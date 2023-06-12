use super::*;

impl Interp {
    pub fn exec(&mut self, node: &IRNode) -> Result<Value, InterpError> {
        match &node.data {
            IRNodeData::Block { body, .. } => self.exec_block(&body),
            IRNodeData::FnCall { func, args, .. } => {
                let args = self.exec_args(args)?;
                self.run_fn(*func, args)
            }
            IRNodeData::Boolean(left, op, right) => self.exec_boolop(*op, left, right),
            IRNodeData::Comparison(left, op, right) => self.exec_comp(*op, left, right),
            IRNodeData::Int(v) => Ok(Value::Int(*v)),
            IRNodeData::Float(v) => Ok(Value::Float(*v)),
            IRNodeData::Char(v) => Ok(Value::Char(*v)),
            IRNodeData::If {
                cond,
                body,
                els,
                ret_typ,
            } => self.exec_if(cond, body, els, ret_typ),
            IRNodeData::Variable(ind, _) => self.exec_var(*ind),
            IRNodeData::NewArrayLiteral(_, vals) => self.exec_arrlit(vals),
            IRNodeData::Print(v) => self.exec_print(v),
            IRNodeData::Define { var, val, .. } => self.exec_define(*var, val),
            IRNodeData::NewEnum(val, ..) => self.exec_newenum(val),
            IRNodeData::NewStruct(_, vals) => self.exec_newstruct(vals),
            IRNodeData::GetEnum { enm, typ } => self.exec_getenum(enm, typ),
            IRNodeData::GetStruct { strct, field } => self.exec_getstruct(strct, field),
            IRNodeData::SetStruct { strct, vals } => self.exec_setstruct(strct, vals),
            IRNodeData::TypeMatch { val, body } => self.exec_typematch(val, body),
            IRNodeData::Match { val, body } => self.exec_match(val, body),
            IRNodeData::NewTuple(_, args) => self.exec_newtuple(args),
            IRNodeData::NewBox(val) => self.exec_newbox(val),
            IRNodeData::Peek { bx, typ } => self.exec_peek(bx, typ),
            IRNodeData::Unbox { bx, typ } => self.exec_unbox(bx, typ),
            IRNodeData::GetTuple { tup, ind } => self.exec_gettuple(tup, ind),
            IRNodeData::While { cond, body } => self.exec_while(cond, body),
            IRNodeData::Len(arr) => self.exec_len(arr),
            IRNodeData::Math(left, op, right) => self.exec_math(node.pos, *op, left, right),
            IRNodeData::GetArr { arr, ind } => self.exec_getarr(arr, ind),
            IRNodeData::NewArray(_, cap) => self.exec_newarr(cap),
            IRNodeData::Append { arr, val } => self.exec_append(arr, val),
            IRNodeData::Return(val) => self.exec_ret(val),
            IRNodeData::Cast(val, typ) => self.exec_cast(val, typ),
            IRNodeData::SetArr { arr, ind, val } => self.exec_setarr(arr, ind, val),
            _ => Err(InterpError::UnknownNode(node.clone())),
        }
    }
}
