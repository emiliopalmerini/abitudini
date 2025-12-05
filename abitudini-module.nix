{ config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.abitudini;
  abitudiniPkg = pkgs.callPackage ./. { };
in

{
  options.services.abitudini = {
    enable = mkEnableOption "Abitudini habit tracker service";

    port = mkOption {
      type = types.port;
      default = 8080;
      description = "Port for the Abitudini HTTP server to listen on";
    };

    dataDir = mkOption {
      type = types.path;
      default = "/var/lib/abitudini";
      description = "Directory where Abitudini will store its database and data";
    };

    package = mkOption {
      type = types.package;
      default = abitudiniPkg;
      description = "The abitudini package to use";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.abitudini = {
      description = "Abitudini Habit Tracker";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        # Execution
        ExecStart = "${cfg.package}/bin/abitudini";
        WorkingDirectory = "${cfg.package}/share/abitudini";

        # User/Security
        DynamicUser = true;
        StateDirectory = "abitudini";

        # File access
        ReadWritePaths = [ cfg.dataDir ];
        StateDirectoryMode = "0755";

        # Environment
        Environment = "PORT=${toString cfg.port}";

        # Process isolation
        Type = "simple";
        Restart = "on-failure";
        RestartSec = 10;

        # Sandboxing
        PrivateTmp = true;
        NoNewPrivileges = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        ProtectKernelTunables = true;
        ProtectKernelModules = true;
        ProtectControlGroups = true;
        RestrictNamespaces = true;
        RestrictRealtime = true;
      };
    };

    # Ensure data directory exists
    systemd.tmpfiles.rules = [
      "d ${cfg.dataDir} 0755 abitudini abitudini - -"
    ];
  };
}
