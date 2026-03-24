{
  description = "Hexagonal architecture with Avro ports and Nix-wired adapters";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
        # keep-sorted start
        ./nix/modules/ports.nix
        # keep-sorted end
      ];

      systems = [
        # keep-sorted start
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
        # keep-sorted end
      ];

      perSystem =
        {
          config,
          self',
          inputs',
          pkgs,
          system,
          ...
        }:
        {
          # Composition root: declare which adapter backs each Repository port.
          # Change language or backend here — no application code changes.
          ports.repositories.item = {
            language = "go";
            backend = "memory";
          };

          # Nix-built adapter packages. Both are sandboxed (no network at build time):
          #   go-adapters     → item-demo, item-server  (from adapters/go/vendor/)
          #   python-adapters → item-demo, item-server  (from pyproject.toml + setuptools)
          packages = {
            go-adapters = pkgs.callPackage ./nix/pkgs/go.nix {
              # builtins.path renames the store entry; without this, the directory
              # name "go" conflicts with buildGoModule's GOPATH and go.mod is ignored.
              src = builtins.path {
                path = ./adapters/go;
                name = "avro-adapters-go-src";
              };
            };
            python-adapters = pkgs.callPackage ./nix/pkgs/python.nix { src = ./adapters/python; };
          };

          devShells.default = import ./shell.nix { inherit pkgs; };
          treefmt.config = import ./treefmt.nix;
        };
    };
}
