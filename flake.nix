{
  description = "Abitudini - Habit Tracker Application";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          abitudini = pkgs.callPackage ./nix/package.nix {};
          default = self.packages.${system}.abitudini;
        };

        devShells.default = pkgs.callPackage ./nix/devShell.nix {};
      }
    );
}
