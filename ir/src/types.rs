use super::*;

#[derive(Debug, Clone)]
pub struct Type {
    pub data: TypeData,
}

impl PartialEq for Type {
    fn eq(&self, other: &Self) -> bool {
        self.data == other.data
    }
}

impl From<TypeData> for Type {
    fn from(data: TypeData) -> Self {
        Self { data }
    }
}

impl TypeData {
    pub fn concrete(&self, ir: &IR) -> Self {
        match self {
            TypeData::DEF(v) => ir.types[*v].typ.data.concrete(ir),
            _ => self.clone(),
        }
    }
}

impl Type {
    pub fn void() -> Self {
        Self {
            data: TypeData::VOID,
        }
    }

    pub fn expect(&self, pos: Pos, expected: &Type) -> Result<(), IRError> {
        if self == expected || self.data == TypeData::INVALID {
            Ok(())
        } else {
            Err(IRError::InvalidType {
                pos,
                expected: expected.clone(),
                got: self.clone(),
            })
        }
    }
}

#[derive(Debug, PartialEq, Clone)]
pub enum TypeData {
    INT,
    FLOAT,
    CHAR,
    BOOL,
    BOX,
    ARRAY(Box<Type>),
    STRUCT(Vec<Field>),
    TUPLE(Vec<Type>),
    ENUM(Vec<Type>),
    DEF(usize), // Points to typedef

    // Macros
    CONSTRAINT(Vec<Type>),

    // Special types
    INVALID,
    PARAM,
    FIELD,
    TYPE,
    VOID,
}

#[derive(Debug, PartialEq, Clone)]
pub struct Field {
    pub name: String,
    pub typ: Type,
}
