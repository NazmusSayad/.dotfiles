use std::{
  fs,
  io::{self, Write},
  process::{exit, Command},
};

fn main() {
  println!("This will reset the entire repository to the latest remote branch.");
  println!("Write 'yes' and press [Enter] to confirm.");
  print!("> ");
  io::stdout().flush().ok();

  let mut confirm = String::new();
  io::stdin().read_line(&mut confirm).ok();

  if confirm.trim() != "yes" {
    println!("Reset aborted");
    return;
  }

  let status = Command::new("git")
    .args(["fetch", "--all"])
    .status()
    .expect("Failed to run git fetch");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let remote_url = Command::new("git")
    .args(["remote", "get-url", "origin"])
    .output()
    .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
    .unwrap_or_default();

  let current_branch = Command::new("git")
    .args(["branch", "--show-current"])
    .output()
    .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
    .unwrap_or_default();

  println!("> Branch: {}", current_branch);
  println!("> Remote: {}", remote_url);

  let remote_output = Command::new("git")
    .args(["branch", "-r", "--format=%(refname:short)"])
    .output()
    .expect("Failed to list remote branches");

  let remote_branches: Vec<String> = String::from_utf8_lossy(&remote_output.stdout)
    .lines()
    .map(|l| l.trim().to_string())
    .filter(|l| !l.is_empty())
    .collect();

  for rb in remote_branches {
    if let Some((_, name)) = rb.split_once('/') {
      if name == current_branch {
        continue;
      }

      println!("> Deleting remote branch: {}", name);

      let _ = Command::new("git")
        .args(["push", "origin", "--delete", name])
        .status();
    }
  }

  println!("> Deleting git folder...");
  let _ = fs::remove_dir_all(".git");

  let init_arg = format!("--initial-branch={}", current_branch);

  let status = Command::new("git")
    .args(["init", &init_arg])
    .status()
    .expect("Failed to init repository");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args(["remote", "add", "origin", &remote_url])
    .status()
    .expect("Failed to add remote");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args(["add", "."])
    .status()
    .expect("Failed to stage files");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args(["commit", "-m", "Initial commit"])
    .status()
    .expect("Failed to commit");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args([
      "push",
      "--force",
      "--set-upstream",
      "origin",
      &current_branch,
    ])
    .status()
    .expect("Failed to push");

  exit(status.code().unwrap_or(1));
}

