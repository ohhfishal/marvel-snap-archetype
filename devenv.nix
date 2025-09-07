{ pkgs, lib, config, ... }: {

  languages.python = {
    enable = true;
    version = "3.13.5";
  };

  # packages = with pkgs; [
  # ];

  git-hooks.hooks = {
    shellcheck.enable = true;
    ruff.enable = true;
    unit-tests = {
      enable = true;
      name = "Unit tests";
      entry = "python -m unittest discover .";
      language = "python";
      # pass_filenames = false;
    };
  };
}
