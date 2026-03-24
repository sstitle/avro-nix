# flake-parts module: ports.repositories
#
# Defines typed options for every Repository port in the system. Each entity
# (item, order, …) gets an entry under ports.repositories. The module is the
# Nix-level dependency-injection mechanism: change an adapter here and nothing
# in the application code changes.
#
# The module also emits a packages.repository-config derivation containing the
# resolved configuration as JSON, which demo/server processes can read at
# runtime to know which adapter to instantiate.
#
# Adding a new entity:
#   1. Define its Avro Protocol in ports/<entity>/repository.avpr
#   2. Implement adapters in adapters/{go,python}/<entity>/
#   3. Add  ports.repositories.<entity> = { language = "go"; backend = "memory"; };
#      to the perSystem block in flake.nix — no module changes needed.
{ lib, flake-parts-lib, ... }:
let
  inherit (lib) mkOption types;
  inherit (flake-parts-lib) mkPerSystemOption;

  repositorySubmodule = types.submodule {
    options = {
      language = mkOption {
        type = types.enum [
          "go"
          "python"
        ];
        description = "Implementation language for this repository adapter.";
      };

      backend = mkOption {
        type = types.enum [
          "memory"
          "rpc"
        ];
        description = ''
          Storage backend for this repository adapter.
            memory  — in-process, no persistence (default for dev/test)
            rpc     — delegates to a remote server over HTTP+JSON
        '';
      };

      rpc = {
        host = mkOption {
          type = types.str;
          default = "localhost";
          description = "Hostname of the remote Repository server (backend = rpc only).";
        };
        port = mkOption {
          type = types.port;
          default = 8080;
          description = "Port of the remote Repository server (backend = rpc only).";
        };
      };
    };
  };
in
{
  options.perSystem = mkPerSystemOption (
    { config, pkgs, ... }:
    {
      options.ports = {
        repositories = mkOption {
          type = types.attrsOf repositorySubmodule;
          default = { };
          description = ''
            Repository port bindings. Each attribute is an entity name mapped to
            an adapter configuration. The Nix module system is the composition root:
            change adapter here, not in application code.

            Example (multiple entities, mixed languages and backends):
              ports.repositories = {
                item  = { language = "go";     backend = "memory"; };
                order = { language = "python"; backend = "rpc"; rpc.port = 9090; };
              };
          '';
        };
      };

      config.packages.repository-config = pkgs.writeText "repository-config.json" (
        builtins.toJSON config.ports.repositories
      );
    }
  );
}
