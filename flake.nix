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
          abitudini = pkgs.buildGoModule {
            pname = "abitudini";
            version = "0.1.0";
            src = pkgs.lib.cleanSource ./.;

            vendorHash = "sha256-HSLR1NOKphqcqcnTyy35teJ3gcPeXUVGaAFlFwSHZEA=";

            # Enable CGO for SQLite support
            env.CGO_ENABLED = "1";

            # Build dependencies for SQLite
            buildInputs = with pkgs; [
              sqlite
              pkg-config
            ];

            # Copy static assets to output
            postInstall = ''
              mkdir -p $out/share/abitudini
              cp -r static $out/share/abitudini/
            '';

            meta = with pkgs.lib; {
              description = "A habit tracking application built with Go";
              homepage = "https://github.com/emiliopalmerini/abitudini";
              license = licenses.mit;
              maintainers = [ ];
            };
          };

          default = self.packages.${system}.abitudini;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            sqlite
            pkg-config
          ];
        };
      }
    );
}
