{
  pkgs ? import <nixpkgs> { },
}:
let
  avro-tools-wrapped = pkgs.writeShellScriptBin "avro-tools" ''
    exec ${pkgs.avro-tools}/bin/avro-tools "$@" 2> >(grep -v "NativeCodeLoader" >&2)
  '';
in
pkgs.mkShell {
  name = "avro-nix";

  buildInputs = with pkgs; [
    # Core tools
    git
    mask
    avro-tools-wrapped
    uv
  ];

  shellHook = ''
    echo "🚀 Development environment loaded!"
    echo "Available tools:"
    echo "  - mask: Task runner"
    echo ""
    echo "Run 'mask --help' to see available tasks."
    echo "Run 'nix fmt' to format all files."
  '';
}
