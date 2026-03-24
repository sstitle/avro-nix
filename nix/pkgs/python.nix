# Builds the Python adapter package (avro-adapters) from adapters/python/.
#
# Installs the avro_adapters namespace package and two entry-point scripts:
#   item-demo   — runs the hexagonal DI demo
#   item-server — starts the HTTP+JSON Repository server
#
# To add a new entity, add its avro_adapters/<entity>/ subpackage to
# [tool.setuptools].packages in adapters/python/pyproject.toml; no changes
# needed here.
{
  python3Packages,
  src,
}:
python3Packages.buildPythonPackage {
  pname = "avro-adapters-python";
  version = "0.1.0";
  inherit src;

  pyproject = true;

  build-system = [
    python3Packages.setuptools
  ];

  propagatedBuildInputs = [
    python3Packages.rich
  ];

  meta.description = "Avro hexagonal adapter package (Python)";
}
