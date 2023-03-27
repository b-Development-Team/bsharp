use fset::FSet;
use interp::Interp;
use ir::IR;

fn main() {
    // Fset
    let mut fset = FSet::new();
    fset.import(&std::path::Path::new(&"test".to_string()))
        .unwrap();

    // IR
    let mut ir = IR::new(fset);
    ir.build().unwrap();
    if ir.errors.len() > 0 {
        for err in ir.errors {
            println!("{:?}\n", err);
        }
        return;
    }

    // Run
    let mut interp = Interp::new(ir, Box::new(std::io::stdout()));
    let mainind = interp.ir.funcs.iter().position(|f| f.name == "@main");
    if mainind.is_none() {
        panic!("main func not found");
    }
    let res = interp.run_fn(mainind.unwrap(), Vec::new()).unwrap();
    println!("OUTPUT: {:?}", res);
}
