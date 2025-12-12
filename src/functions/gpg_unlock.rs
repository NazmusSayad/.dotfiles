use std::{env, fs, path::PathBuf, process::Command};

fn main() {
  if let Ok(output) = Command::new("ps").args(["aux"]).output() {
    let stdout = String::from_utf8_lossy(&output.stdout);

    for keyword in ["gpg", "keyboxd"] {
      for line in stdout.lines() {
        if line.contains(keyword) && !line.contains("grep") {
          if let Some(pid) = line
            .split_whitespace()
            .find(|field| field.chars().all(|c| c.is_ascii_digit()))
          {
            println!("Found {} process with PID: {}", keyword, pid);
            let _ = Command::new("sudo").args(["kill", "-9", pid]).status();
          }
        }
      }
    }
  }

  let home = env::var("HOME")
    .or_else(|_| env::var("USERPROFILE"))
    .unwrap_or_else(|_| ".".to_string());

  let gnupg = PathBuf::from(home).join(".gnupg");

  if let Ok(entries) = fs::read_dir(gnupg) {
    for entry in entries.flatten() {
      let path = entry.path();

      if path.extension().and_then(|e| e.to_str()) == Some("lock") {
        let _ = fs::remove_file(path);
      }
    }
  }
}
