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
        op: MathOperator,
        left: &IRNode,
        right: &IRNode,
    ) -> Result<Value, InterpError> {
        let left = self.exec(left)?;
        let right = self.exec(right)?;
        match (left, right) {
            (Value::Int(left), Value::Int(right)) => match op {
                MathOperator::ADD => Ok(Value::Int(left + right)),
                MathOperator::SUBTRACT => Ok(Value::Int(left - right)),
                MathOperator::MULTIPLY => Ok(Value::Int(left * right)),
                MathOperator::DIVIDE => Ok(Value::Int(left / right)),
                MathOperator::MODULO => Ok(Value::Int(left % right)),
                MathOperator::XOR => Ok(Value::Int(left ^ right)),
                MathOperator::SHIFT => Ok(Value::Int(left << right)),
                MathOperator::BOR => Ok(Value::Int(left | right)),
            },
            (Value::Float(left), Value::Float(right)) => match op {
                MathOperator::ADD => Ok(Value::Float(left + right)),
                MathOperator::SUBTRACT => Ok(Value::Float(left - right)),
                MathOperator::MULTIPLY => Ok(Value::Float(left * right)),
                MathOperator::DIVIDE => Ok(Value::Float(left / right)),
                MathOperator::MODULO => Ok(Value::Float(left % right)),
                MathOperator::XOR => Ok(Value::Float((left as i64 ^ right as i64) as f64)),
                MathOperator::SHIFT => Ok(Value::Float(((left as i64) << (right as i64)) as f64)),
                MathOperator::BOR => Ok(Value::Float((left as i64 | right as i64) as f64)),
            },
            _ => unreachable!(),
        }
    }
}
