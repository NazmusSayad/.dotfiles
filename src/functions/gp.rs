use std::process::Command;

fn main() {
  let args: Vec<String> = std::env::args().skip(1).collect();
  if args.len() > 1 {
    eprintln!("Usage: gp [branch]");
    std::process::exit(1);
  }

  let current_output = Command::new("git")
    .args(["branch", "--show-current"])
    .output()
    .expect("Failed to read current branch");
  if !current_output.status.success() {
    std::process::exit(current_output.status.code().unwrap_or(1));
  }
  let current_branch = String::from_utf8_lossy(&current_output.stdout)
    .trim()
    .to_string();

  let target_branch = if args.is_empty() {
    current_branch.clone()
  } else {
    args[0].clone()
  };

  if target_branch.is_empty() {
    eprintln!("Unable to determine target branch");
    std::process::exit(1);
  }

  println!("Pulling {} into {} (default)", target_branch, current_branch);

  let prune_status = Command::new("git")
    .args(["prune", "--progress"])
    .status()
    .expect("Failed to run git prune");
  if !prune_status.success() {
    std::process::exit(prune_status.code().unwrap_or(1));
  }

  let pull_status = Command::new("git")
    .args(["pull", "origin", target_branch.as_str(), "--progress"])
    .status()
    .expect("Failed to run git pull");

  if !pull_status.success() {
    std::process::exit(pull_status.code().unwrap_or(1));
  }
}

