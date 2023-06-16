use super::*;
use include_dir::{include_dir, Dir};

static STDLIB: Dir = include_dir!("$CARGO_MANIFEST_DIR/stdlib");

impl FSet {
    pub fn include_stdlib(&mut self) {
        for f in STDLIB.files() {
            if let Err(e) = self.add_file_source(
                f.path().to_str().unwrap().to_string(),
                f.contents_utf8().unwrap().to_string(),
            ) {
                panic!("{}", e.fmt(self))
            };
        }
    }
}
