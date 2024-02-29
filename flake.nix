{
  description = "uTCC";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-23.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { ... } @ inputs: inputs.utils.lib.eachDefaultSystem (system:
    let
      pkgs = import inputs.nixpkgs { inherit system; };
    in
    {
      devShell = pkgs.mkShell {
        packages = with pkgs; [
          go

          dapr-cli

          minikube
          kubectl
          kubernetes-helm
        ];

        GOROOT = "${pkgs.go}/share/go";
      };
    });
}
