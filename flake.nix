{
  description = "Check the runtime of program executions";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.best-of = pkgs.buildGoModule {
          pname = "best-of";
          version = self.shortRev or "dev";
          src = ./.;
          vendorHash = null;
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.best-of}/bin/best-of";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go ];
        };
      }
    );
}
