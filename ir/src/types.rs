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

    pub fn fmt(&self, ir: &IR) -> String {
        match self {
            TypeData::DEF(v) => ir.types[*v].name.clone(),
            TypeData::ARRAY(v) => format!("Array<{}>", v.data.fmt(ir)),
            TypeData::STRUCT(v) => {
                let mut s = String::new();
                s.push_str("[STRUCT");
                for f in v {
                    s.push_str(&format!(" [FIELD {} {}]", f.name, f.typ.data.fmt(ir)));
                }
                s.push_str("]");
                s
            }
            TypeData::TUPLE(v) => {
                let mut s = String::new();
                s.push_str("[TUPLE");
                for f in v {
                    s.push_str(&format!(" {}", f.data.fmt(ir)));
                }
                s.push_str("]");
                s
            }
            TypeData::ENUM(v) => {
                let mut s = String::new();
                s.push_str("[ENUM");
                for f in v {
                    s.push_str(&format!(" {}", f.data.fmt(ir)));
                }
                s.push_str("]");
                s
            }
            TypeData::INT => "[INT]".to_string(),
            TypeData::FLOAT => "[FLOAT]".to_string(),
            TypeData::CHAR => "[CHAR]".to_string(),
            TypeData::BOOL => "[BOOL]".to_string(),
            TypeData::BOX => "[BX]".to_string(),
            TypeData::INVALID => "INVALID".to_string(),
            TypeData::PARAM => "PARAM".to_string(),
            TypeData::FIELD => "FIELD".to_string(),
            TypeData::TYPE => "TYPE".to_string(),
            TypeData::VOID => "VOID".to_string(),
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
