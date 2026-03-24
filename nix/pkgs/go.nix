# Builds the Go adapter binaries (item-demo, item-server) from adapters/go/.
#
# Uses the vendored dependencies in adapters/go/vendor/ so no network access
# is required at build time (compatible with Nix sandbox).
#
# To add a new entity's commands, add its cmd/ subpackage to subPackages and
# extend postInstall with the rename.
{
  lib,
  buildGoModule,
  src,
}:
buildGoModule {
  pname = "avro-adapters-go";
  version = "0.1.0";
  inherit src;

  # Hash of `go mod vendor` output. Update this whenever go.mod changes:
  #   nix build .#go-adapters  # fails with "got: sha256-..."
  #   then paste that hash here.
  vendorHash = "sha256-DAgqc+wVjXZytSZcZENF3lxIBhtpamBpAPuKNrb58pk=";

  subPackages = [
    # keep-sorted start
    "cmd/demo"
    "cmd/server"
    # keep-sorted end
  ];

  # Rename generic cmd names to entity-scoped names.
  postInstall = ''
    mv $out/bin/demo   $out/bin/item-demo
    mv $out/bin/server $out/bin/item-server
  '';

  meta = {
    description = "Avro hexagonal adapter binaries (Go)";
    mainProgram = "item-server";
  };
}
