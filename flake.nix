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

          devShells.default = import ./shell.nix { inherit pkgs; };
          treefmt.config = import ./treefmt.nix;
        };
    };
}
