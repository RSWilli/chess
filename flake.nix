{
  description = "golang development flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    system = "x86_64-linux";
  in {
    devShells."${system}".default = let
      pkgs = import nixpkgs {
        inherit system;
      };
      local_go = pkgs.go_1_25;
    in
      pkgs.mkShell {
        packages = with pkgs; [
          local_go
          graphviz # for pprof
          stockfish
          rlwrap # for terminal command history inside stockfish
        ];

        GO111MODULE = "on";

        # needed for running delve with cgo
        # https://wiki.nixos.org/wiki/Go#Using_cgo_on_NixOS
        hardeningDisable = ["fortify"];

        shellHook = ''
          ${local_go}/bin/go version
        '';
      };
  };
}
