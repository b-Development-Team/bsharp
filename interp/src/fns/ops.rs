use super::*;

impl Interp {
    pub fn exec_boolop(
        &mut self,
        op: BooleanOperator,
        left: &IRNode,
        right: &Option<Box<IRNode>>,
    ) -> Result<Value, InterpError> {
        if op == BooleanOperator::NOT {
            if let Value::Bool(val) = self.exec(left)? {
                return Ok(Value::Bool(!val));
            }
        }

        let (left, right) = match (self.exec(left)?, self.exec(right.as_ref().unwrap())?) {
            (Value::Bool(left), Value::Bool(right)) => (left, right),
            _ => unreachable!(),
        };
        match op {
            BooleanOperator::AND => Ok(Value::Bool(left && right)),
            BooleanOperator::OR => Ok(Value::Bool(left || right)),
            _ => unreachable!(),
        }
    }

    pub fn exec_comp(
        &mut self,
        op: ComparisonOperator,
        left: &IRNode,
        right: &IRNode,
    ) -> Result<Value, InterpError> {
        let left = self.exec(left)?;
        let right = self.exec(right)?;
        match (left, right) {
            (Value::Char(left), Value::Char(right)) => match op {
                ComparisonOperator::GREATER => Ok(Value::Bool(left > right)),
                ComparisonOperator::LESS => Ok(Value::Bool(left < right)),
                ComparisonOperator::GREATEREQUAL => Ok(Value::Bool(left >= right)),
                ComparisonOperator::LESSEQUAL => Ok(Value::Bool(left <= right)),
                ComparisonOperator::EQUAL => Ok(Value::Bool(left == right)),
                ComparisonOperator::NOTEQUAL => Ok(Value::Bool(left != right)),
            },
            (Value::Int(left), Value::Int(right)) => match op {
                ComparisonOperator::GREATER => Ok(Value::Bool(left > right)),
                ComparisonOperator::LESS => Ok(Value::Bool(left < right)),
                ComparisonOperator::GREATEREQUAL => Ok(Value::Bool(left >= right)),
                ComparisonOperator::LESSEQUAL => Ok(Value::Bool(left <= right)),
                ComparisonOperator::EQUAL => Ok(Value::Bool(left == right)),
                ComparisonOperator::NOTEQUAL => Ok(Value::Bool(left != right)),
            },
            (Value::Float(left), Value::Float(right)) => match op {
                ComparisonOperator::GREATER => Ok(Value::Bool(left > right)),
                ComparisonOperator::LESS => Ok(Value::Bool(left < right)),
                ComparisonOperator::GREATEREQUAL => Ok(Value::Bool(left >= right)),
                ComparisonOperator::LESSEQUAL => Ok(Value::Bool(left <= right)),
                ComparisonOperator::EQUAL => Ok(Value::Bool(left == right)),
                ComparisonOperator::NOTEQUAL => Ok(Value::Bool(left != right)),
            },
            (Value::Bool(left), Value::Bool(right)) => match op {
                ComparisonOperator::EQUAL => Ok(Value::Bool(left == right)),
                ComparisonOperator::NOTEQUAL => Ok(Value::Bool(left != right)),
                _ => unreachable!(),
            },
            _ => unreachable!(),
        }
    }

    pub fn exec_math(
        &mut self,
        pos: Pos,
        op: MathOperator,
        left: &IRNode,
        right: &IRNode,
    ) -> Result<Value, InterpError> {
        let left = self.exec(left)?;
        let right = self.exec(right)?;
        match (left, right) {
            (Value::Int(left), Value::Int(right)) => match op {
                MathOperator::ADD => Ok(Value::Int(left.wrapping_add(right))),
                MathOperator::SUBTRACT => Ok(Value::Int(left - right)),
                MathOperator::MULTIPLY => Ok(Value::Int(left.wrapping_mul(right))),
                MathOperator::DIVIDE => {
                    if right == 0 {
                        return Err(InterpError::DivideByZero(pos));
                    }
                    Ok(Value::Int(left / right))
                }
                MathOperator::MODULO => Ok(Value::Int(left % right)),
                MathOperator::XOR => Ok(Value::Int(left ^ right)),
                MathOperator::SHIFT => Ok(Value::Int(left.wrapping_shl(right as u32))),
                MathOperator::BOR => Ok(Value::Int(left | right)),
            },
            (Value::Float(left), Value::Float(right)) => match op {
                MathOperator::ADD => Ok(Value::Float(left + right)),
                MathOperator::SUBTRACT => Ok(Value::Float(left - right)),
                MathOperator::MULTIPLY => Ok(Value::Float(left * right)),
                MathOperator::DIVIDE => {
                    if right == 0.0 {
                        return Err(InterpError::DivideByZero(pos));
                    }
                    Ok(Value::Float(left / right))
                }
                MathOperator::MODULO => Ok(Value::Float(left % right)),
                MathOperator::XOR => Ok(Value::Float((left as i64 ^ right as i64) as f64)),
                MathOperator::SHIFT => Ok(Value::Float(((left as i64) << (right as i64)) as f64)),
                MathOperator::BOR => Ok(Value::Float((left as i64 | right as i64) as f64)),
            },
            (Value::Char(left), Value::Char(right)) => match op {
                MathOperator::ADD => Ok(Value::Char(left + right)),
                MathOperator::SUBTRACT => Ok(Value::Char(left - right)),
                MathOperator::MULTIPLY => Ok(Value::Char(left * right)),
                MathOperator::DIVIDE => {
                    if right == 0 {
                        return Err(InterpError::DivideByZero(pos));
                    }
                    Ok(Value::Char(left / right))
                }
                MathOperator::MODULO => Ok(Value::Char(left % right)),
                MathOperator::XOR => Ok(Value::Char(left ^ right)),
                MathOperator::SHIFT => Ok(Value::Char(left << right)),
                MathOperator::BOR => Ok(Value::Char(left | right)),
            },
            _ => unreachable!(),
        }
    }

    pub fn exec_cast(&mut self, arg: &IRNode, target: &Type) -> Result<Value, InterpError> {
        let val = self.exec(arg)?;
        match target.data.concrete(&self.ir) {
            TypeData::CHAR => match val {
                Value::Int(val) => Ok(Value::Char(val as u8)),
                Value::Float(val) => Ok(Value::Char(val as u8)),
                _ => unreachable!(),
            },
            TypeData::INT => match val {
                Value::Char(val) => Ok(Value::Int(val as i64)),
                Value::Float(val) => Ok(Value::Int(val as i64)),
                _ => unreachable!(),
            },
            TypeData::FLOAT => match val {
                Value::Char(val) => Ok(Value::Float(val as f64)),
                Value::Int(val) => Ok(Value::Float(val as f64)),
                _ => unreachable!(),
            },
            _ => unreachable!(),
        }
    }
}
