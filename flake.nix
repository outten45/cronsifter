{
  description = "go and stuff";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.default = pkgs.mkShell {
        nativeBuildInputs = [ 
          # pkgs.bashInteractive 
          pkgs.buildPackages.go
          pkgs.buildPackages.gnumake
          pkgs.buildPackages.gcc
          pkgs.buildPackages.sqlite-interactive
          pkgs.buildPackages.readline
        ];
        buildInputs = [ ];
        shellHook = ''
          echo "Starting nix-shell with fish..."
        '';
      };
    });
}
