{ pkgs, inputs, lib, config, ... }:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{

  languages.python = {
    enable = true;
    version = "3.13.5";
  };

  packages = [
    pkgs-unstable.go
  ];

  git-hooks.hooks = {
    # Shell
    shellcheck.enable = true;

    # Golang
    govet.enable = true;
    gotest.enable = true;
    gofmt.enable = true;

    # Python
    ruff.enable = true;
    unit-tests = {
      enable = true;
      name = "py-unit-tests";
      types = [ "python" ];
      entry = "python -m unittest discover .";
      language = "python";
      pass_filenames = false;
    };
  };
}
