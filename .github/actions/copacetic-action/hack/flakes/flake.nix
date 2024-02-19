{
  description = "Flake for copacetic";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs @ { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        lib = nixpkgs.lib;
        goPackageName = "github.com/project-copacetic/copacetic";
      in
      {
        packages = {
          copacetic = pkgs.buildGo121Module rec {
            pname = "copacetic";
            version = "0.6.0";

            src = pkgs.fetchFromGitHub {
              owner = "project-copacetic";
              repo = "copacetic";
              rev = "v${version}";
              sha256 = "sha256-0lMRpungk2wZ3o8VO3xZDmAn/FqGq4OXlRACnPzyL7c=";
            };

            subPackages = ["."];

            ldflags = [
              "-X ${goPackageName}/pkg/version.GitVersion=${src.rev}"
              "-X ${goPackageName}/pkg/version.GitVersion=${src.rev}"
              "-X main.version=v${version}"
              "-s -w -extldflags -static"
            ];

            vendorHash = "sha256-H4Hda5PjaqbIQiieX6R/YyzOAZEfw1bI3xFqaXQXNtY=";

            doCheck = false;

            CGO_ENABLED = 0;

            installPhase = ''
              runHook preInstall

              mkdir -p $out
              dir="$GOPATH/bin"
              [ -e "$dir" ] && cp -r $dir $out
              mv $out/bin/${pname} $out/bin/copa

              runHook postInstall
            '';

            meta = with lib; {
              description = "";
              homepage = "https://project-copacetic.github.io/copacetic/website/";
              license = licenses.asl20;
            };
          };
        };

        formatter = pkgs.alejandra;
      }
    );
}
