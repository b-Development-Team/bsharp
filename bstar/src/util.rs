use super::*;

impl BStar {
    pub fn structfields(&mut self, t: &Type) -> Vec<String> {
        let mut names = Vec::new();
        if let TypeData::STRUCT(f) = t.data.concrete(&self.ir) {
            for field in f.iter() {
                names.push(field.name.clone());
            }
        }
        names.sort();
        return names;
    }

    pub fn hashtyp(&mut self, t: &Type) -> String {
        match &t.data {
            TypeData::DEF(v) => format!("d{:x}", v),
            TypeData::INT => "i".to_string(),
            TypeData::FLOAT => "f".to_string(),
            TypeData::CHAR => "c".to_string(),
            TypeData::BOOL => "b".to_string(),
            TypeData::BOX => "x".to_string(),
            TypeData::ARRAY(t) => format!("r{}", self.hashtyp(&t)),
            TypeData::STRUCT(f) => {
                let mut s = String::new();
                s.push_str("s");
                for field in f.iter() {
                    s.push_str(&self.hashtyp(&field.typ));
                }
                s
            }
            TypeData::TUPLE(f) => {
                let mut s = String::new();
                s.push_str("t");
                for field in f.iter() {
                    s.push_str(&self.hashtyp(field));
                }
                s
            }
            TypeData::ENUM(f) => {
                let mut s = String::new();
                s.push_str("e");
                for field in f.iter() {
                    s.push_str(&self.hashtyp(field));
                }
                s
            }
            _ => unreachable!(),
        }
    }
}
