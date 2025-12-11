use std::process::Command;

fn main() {
  let cwd = std::env::current_dir().expect("Failed to read cwd");
  println!("👉 Current directory: {}", cwd.display());

  let output = Command::new("fsutil.exe")
    .args(&["file", "setCaseSensitiveInfo", ".", "enable", "recursive"])
    .output()
    .expect("Failed to execute fsutil command");

  if output.status.success() {
    println!("✅ Case sensitivity enabled for current directory");
  } else {
    eprintln!("Failed to enable case sensitivity");
    eprintln!("Error: {}", String::from_utf8_lossy(&output.stderr));
    std::process::exit(1);
  }
}
