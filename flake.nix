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
    in
      pkgs.mkShell {
        packages = with pkgs; [
          go_latest
          go-tools

          gopls
          delve
          impl
          gotest

          graphviz # for pprof
          stockfish
          rlwrap # for terminal command history inside stockfish
        ];

        GO111MODULE = "on";

        # go 1.25 new json implementation
        GOEXPERIMENT = "jsonv2";

        # needed for running delve with cgo
        # https://wiki.nixos.org/wiki/Go#Using_cgo_on_NixOS
        hardeningDisable = ["fortify"];

        shellHook = ''
          ${pkgs.go_latest}/bin/go version
        '';
      };
  };
}
