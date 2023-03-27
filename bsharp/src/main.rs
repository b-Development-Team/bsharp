use clap::{error::ErrorKind, CommandFactory, Parser};
use fset::FSet;
use interp::Interp;
use ir::IR;

#[derive(Parser)]
#[command(name = "bsharp", about = "The B# programming language!")]
struct Args {
    #[arg(help = "the directory to import")]
    dir: std::path::PathBuf,
}

fn main() {
    let args = Args::parse();

    // Fset
    let mut fset = FSet::new();
    if let Err(err) = fset.import(&args.dir) {
        Args::command().error(ErrorKind::Io, err.fmt(&fset)).exit();
    };

    // IR
    let mut ir = IR::new(fset);
    if let Err(err) = ir.build() {
        Args::command()
            .error(ErrorKind::InvalidValue, err.fmt(&ir))
            .exit();
    };
    if ir.errors.len() > 0 {
        for err in ir.errors.iter() {
            println!("{}", err.fmt(&ir));
        }
        println!("");
        Args::command()
            .error(
                ErrorKind::ValueValidation,
                format!("built with {} errors", ir.errors.len()),
            )
            .exit();
    }

    // Run
    let mut interp = Interp::new(ir, Box::new(std::io::stdout()));
    let mainind = interp.ir.funcs.iter().position(|f| f.name == "@main");
    if mainind.is_none() {
        panic!("main func not found");
    }
    if let Err(err) = interp.run_fn(mainind.unwrap(), Vec::new()) {
        Args::command()
            .error(ErrorKind::InvalidValue, format!("{:?}", err)) // TODO: Interp error formatting
            .exit();
    };
}
