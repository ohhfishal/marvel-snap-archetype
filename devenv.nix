{ pkgs, inputs, lib, config, ... }:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{

  languages.python = {
    enable = true;
    version = "3.13.5";
  };

  env.HOST = "localhost";
  env.PORT = "8080";

  dotenv.enable = true;
  dotenv.filename = ".env";

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
  processes = {
    stats-server = {
      exec = "go tool templ generate ./... && go run . serve";
      process-compose = {
        working_dir = "${config.env.DEVENV_ROOT}";
        log_location = "${config.env.DEVENV_ROOT}/logs/fastapi.log";
        availability = {
          restart = "on_failure";
          max_restarts = 3;
          backoff_seconds = 2;
        };
        readiness_probe = {
          http_get = {
            host = config.env.HOST;
            port = config.env.PORT;
            path = "/health";
          };
          initial_delay_seconds = 5;
          period_seconds = 10;
          timeout_seconds = 3;
          success_threshold = 1;
          failure_threshold = 3;
        };
      };
    };
  };
}
