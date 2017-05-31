package tSFInstaller_v2;

import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.StandardOpenOption;
import java.util.HashMap;
import java.util.Map;

public class Installer {
    private static File log = new File("log.log");
    private static final Map<String, String> Repositories;
    private static final String[] FilesToBackup = {
            "Settings.sqf"
            , "tSF_briefing.sqf"
            , "Endings.hpp"
            , "Notes.sqf"
            , "NotesSettings.sqf"
            , "description.ext"
            , "initServer.sqf"
            , "init.sqf"
            , "overview.jpg"

            , "Zones.sqf"
            , "GearAssignementTable.sqf"
            , "Kits.sqf"
    };

    static {
        Repositories = new HashMap<String, String>();
        Repositories.put("REPO_DZN_DYNAI", "https://github.com/10Dozen/dzn_dynai/archive/master.zip");
        Repositories.put("REPO_DZN_GEAR", "https://github.com/10Dozen/dzn_gear/archive/master.zip");
        Repositories.put("REPO_DZN_CommonFunctions", "https://github.com/10Dozen/dzn_commonFunctions/archive/master.zip");
        Repositories.put("REPO_DZN_CIVEN", "https://github.com/10Dozen/dzn_civen/archive/master.zip");
        Repositories.put("REPO_DZN_TSF", "https://github.com/10Dozen/dzn_tSFramework/archive/master.zip");
    }

    public static void Install (
            String path
            , boolean doBackup
            , boolean doCF, boolean doG, boolean doDA, boolean doCN, boolean doTSF
            , String urlCF, String urlG, String urlDA, String urlCN, String urlTSF
            , String kit1, String kit2, String kit3
    ) throws IOException {
        PrintWriter writer = new PrintWriter(log);
        writer.print("");
        writer.close();

        LogToFile(" ------ tSF Installer ------ ");
        LogToFile(" Output folder: ".concat(path));

        File outputFolder = new File(path);
        if (!outputFolder.exists()) {
            LogToFile("Aborted! No INSTALLATION_FOLDER exists!");
            System.exit(0);
        }
        File tempFolder = new File("Temp");

        LogToFile(" Install: ");
        if (doCF) { LogToFile("       - CommonFunctions"); }
        if (doG) { LogToFile("       - Gear"); }
        if (doDA) { LogToFile("       - DynAI"); }
        if (doCN) { LogToFile("       - CivEn"); }
        if (doTSF) { LogToFile("       - tS Framework"); }
        LogToFile(" --------------------------- ");

        ProcessRepository("dzn_DynAI", doDA, urlDA);



    }

    public static void LogToFile(String line) throws IOException {
        System.out.println(line);
        Files.write(log.toPath(), "\n ".concat(line).getBytes(), StandardOpenOption.APPEND );
    }

    public static void ProcessRepository(String name, boolean isNeeded, String url) throws IOException {
        if (!isNeeded) { return; }
        LogToFile("Installing ".concat(name));

        LogToFile("    Downloading... ");
        // File folder = new File( DownloadRepository(url,false) );
        LogToFile("     Done!");

        LogToFile("    Unziping... ");
        // CopyFolder(folder, new File("Temp"));
        // DeleteFolder(folder);
        // DeleteFolder(new File (GetRepoFileName(url)));
        LogToFile("     Done!");
    }

}
