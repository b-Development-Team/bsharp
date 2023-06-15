use std::collections::HashSet;

use tokens::Tokenizer;
mod error;
mod stdlib;
pub use error::FSetError;
pub use parser::*;
pub use tokens::*;

pub struct File {
    pub name: String,
    pub ast: Option<parser::Parser>,
}

pub struct FSet {
    pub files: Vec<File>,
    imported: HashSet<String>,
}

impl FSet {
    pub fn new() -> FSet {
        let mut f = FSet {
            files: Vec::new(),
            imported: HashSet::new(),
        };
        f.include_stdlib();
        f
    }

    pub fn add_file_source(&mut self, name: String, source: String) -> Result<(), FSetError> {
        let ind = self.files.len();
        self.files.push(File {
            name: name.clone(),
            ast: None,
        });

        let mut tok = Tokenizer::new(source, ind);
        tok.tokenize()?;
        let mut p = parser::Parser::new(tok);
        p.parse()?;

        self.files[ind].ast = Some(p);

        Ok(())
    }

    pub fn import(&mut self, dir: &std::path::Path) -> Result<(), FSetError> {
        if self.imported.contains(&dir.to_str().unwrap().to_string()) {
            return Ok(());
        }
        self.imported.insert(dir.to_str().unwrap().to_string());
        // Go through all files in the directory ending in .bsp
        for entry in std::fs::read_dir(dir)? {
            let entry = entry?;
            let path = entry.path();
            if path.extension().is_some() && path.extension().unwrap() == "bsp" {
                let name = path.file_name().unwrap().to_str().unwrap().to_string();
                let source = std::fs::read_to_string(path)?;
                self.add_file_source(name, source)?;
            }
        }
        Ok(())
    }

    pub fn display_pos(&self, pos: &Pos) -> String {
        let file = &self.files[pos.file];
        format!("{}:{}:{}", file.name, pos.start_line + 1, pos.start_col + 1)
    }
}
