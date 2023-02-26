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

    pub fn matches(&self, typs: &Vec<TypeData>, ir: &IR) -> bool {
        match self.concrete(ir) {
            TypeData::INTERFACE { body, .. } => {
                for typ in typs {
                    for v in body.iter() {
                        if std::mem::discriminant(&v.data) == std::mem::discriminant(&typ) {
                            return true;
                        }
                    }
                }
                false
            }
            _ => {
                for t in typs {
                    if std::mem::discriminant(t) == std::mem::discriminant(self) {
                        return true;
                    }
                }
                false
            }
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
    ARRAY {
        param: Option<usize>,
        body: Box<Type>,
        scope: usize,
    }, // First param is generic
    STRUCT {
        params: Vec<usize>,
        fields: Vec<Field>,
        scope: usize,
    },
    TUPLE {
        params: Vec<usize>,
        body: Vec<Type>,
        scope: usize,
    },
    ENUM {
        params: Vec<usize>,
        body: Vec<Type>,
        scope: usize,
    },
    INTERFACE {
        params: Vec<usize>,
        body: Vec<Type>,
        scope: usize,
    },
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
