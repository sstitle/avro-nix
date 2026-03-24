{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShell {
  name = "avro-nix";
  
  buildInputs = with pkgs; [
    # Core tools
    git
    mask
    avro-tools
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
