{
  pkgs,
  gitHooksLib,
}: let
  hooks = gitHooksLib.run {
    src = ../.;
    hooks = {
      gofmt.enable = true;
      govet.enable = true;
    };
  };
in
  pkgs.mkShell {
    packages = with pkgs;
      [
        go
        doctoc
        (python3.withPackages (ps: [ps.pyte]))
      ]
      ++ hooks.enabledPackages;

    shellHook = hooks.shellHook;
  }
