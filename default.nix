{
  pkgs ? import <nixpkgs> { },
}:

pkgs.buildGoModule {
  pname = "todo-go";
  version = "0.0.1";
  src = ./.;
  vendorHash = null;
}
