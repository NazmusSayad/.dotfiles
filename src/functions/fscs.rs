use std::process::Command;

fn main() {
  let status = Command::new("fsutil.exe")
    .args(["file", "setCaseSensitiveInfo", ".", "enable", "recursive"])
    .status()
    .expect("Failed to execute fsutil command");

  if !status.success() {
    std::process::exit(status.code().unwrap_or(1));
  }
}
