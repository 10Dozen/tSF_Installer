package tSFInstaller_v2;

import java.io.*;
import java.net.URL;
import java.net.URLDecoder;
import java.nio.channels.Channels;
import java.nio.channels.ReadableByteChannel;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;
import java.util.*;
import java.util.regex.Pattern;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

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

    public static boolean Install (
            String path
            , boolean doBackup
            , boolean doCF, boolean doG, boolean doDA, boolean doCN, boolean doTSF
            , String urlCF, String urlG, String urlDA, String urlCN, String urlTSF
            , String[] kits
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

        LogToFile("Install:");
        if (doCF) { LogToFile("       - CommonFunctions"); }
        if (doG) { LogToFile("       - Gear"); }
        if (doDA) { LogToFile("       - DynAI"); }
        if (doCN) { LogToFile("       - CivEn"); }
        if (doTSF) { LogToFile("       - tS Framework"); }
        LogToFile(" --------------------------- ");

        /*
        ProcessRepository("dzn_DynAI", doDA, urlDA);
        ProcessRepository("dzn_Gear", doG, urlG);
        ProcessRepository("dzn_CommonFunctions", doCF, urlCF);
        ProcessRepository("dzn_CivEn", doCN, urlCN);
        ProcessRepository("dzn_tSF", doTSF, urlTSF);
*/
        LogToFile("Compiling init.sqf");
        // GenerateInitSQF(doG,doDA,doCN,doTSF);
        LogToFile("\\n --------------------------- \\n dzn_gear Kits:");
        // ProcessKits(doG, path, kits);
        LogToFile(" --------------------------- \n  Installation... ");
        // UpdateFiles(tempFolder, outputFolder, doBackup);
        LogToFile(" Done! \n  --------------------------- \n  All done! Have a nice day!");

        saveSettings(path, doBackup , doCF, doG, doDA, doCN, doTSF, urlCF, urlG, urlDA, urlCN, urlTSF, kits);
        return true;
    }

    public static void LogToFile(String line) throws IOException {
        System.out.println(line);
        Files.write(log.toPath(), "\n".concat(line).getBytes(), StandardOpenOption.APPEND );
    }

    public static void ProcessRepository(String name, boolean isNeeded, String url) throws IOException {
        if (!isNeeded) { return; }
        LogToFile("Installing ".concat(name));

        LogToFile("    Downloading... ");
        File folder = new File( DownloadRepository(url,false) );
        LogToFile("     Done!");

        LogToFile("    Unziping... ");
        CopyFolder(folder, new File("Temp"));
        DeleteFolder(folder);
        DeleteFolder(new File (GetRepoFileName(url)));
        LogToFile("     Done!");
    }

    public static String DownloadRepository(String src, boolean isSingleFile) throws IOException {
        URL website = new URL(src);
        String filename = GetRepoFileName(src);

        ReadableByteChannel rbc = Channels.newChannel(website.openStream());
        FileOutputStream fos = new FileOutputStream(filename);
        fos.getChannel().transferFrom(rbc, 0, Long.MAX_VALUE);

        String out;
        if (isSingleFile) {
            out = filename;
        } else {
            out = UnzipFile(filename);
        }

        rbc.close();
        fos.close();

        return out;
    }

    public static String GetRepoFileName(String src) {
        String[] segments = src.split("/");

        return (segments[segments.length-1]);
    }

    public static String UnzipFile(String filename) {
        File folder = new File ("");
        boolean isBasicFolderExist = false;
        try {
            // Open the zip file
            ZipFile zipFile = new ZipFile(filename);
            Enumeration<?> enu = zipFile.entries();
            while (enu.hasMoreElements()) {
                // System.out.print("|");

                ZipEntry zipEntry = (ZipEntry) enu.nextElement();

                String name = zipEntry.getName();

                // Do we need to create a directory ?
                File file = new File(name);
                if (name.endsWith("/")) {
                    file.mkdirs();
                    if (!isBasicFolderExist) {
                        isBasicFolderExist = true;
                        folder = file;
                    }
                    continue;
                }

                File parent = file.getParentFile();
                if (parent != null) {
                    parent.mkdirs();
                }

                // Extract the file
                InputStream is = zipFile.getInputStream(zipEntry);
                FileOutputStream fos = new FileOutputStream(file);
                byte[] bytes = new byte[1024];
                int length;
                while ((length = is.read(bytes)) >= 0) {
                    fos.write(bytes, 0, length);
                }
                is.close();
                fos.close();

            }
            zipFile.close();
        } catch (IOException e) {
            e.printStackTrace();
        }

        return folder.toString();
    }

    public static void CopyFolder(File src, File dest) throws IOException {
        if (src.isDirectory()) {
            dest.mkdir();
            String files[] = src.list();

            for (String file : files) {
                File srcFile = new File(src, file);
                File destFile = new File(dest, file);
                CopyFolder(srcFile, destFile);
            }
        } else {
            InputStream in = new FileInputStream(src);
            OutputStream out = new FileOutputStream(dest);

            byte[] buffer = new byte[1024];
            int lenght;
            while ( (lenght = in.read(buffer)) > 0 ) {
                out.write(buffer, 0, lenght);
            }

            in.close();
            out.close();
        }
    }

    public static boolean DeleteFolder(File dir) throws IOException {
        if (dir.isDirectory()) {
            File[] children = dir.listFiles();
            for (int i = 0; i < children.length; i++) {
                boolean success = DeleteFolder(children[i]);
                if (!success) { return false; }
            }
        }
        return dir.delete();
    }

    public static void GenerateInitSQF(
           boolean g, boolean d, boolean c, boolean t
    ) throws IOException {
        List<String> lines = new ArrayList<String>();

        lines.add("//	Tacitcal Shift Framework initialization");
        lines.add("[] spawn {");
        lines.add("        waitUntil { !isNil \"MissionDate\" };");

        if (g) {
            lines.add("");
            lines.add("        // dzn Gear 	(set true to engage Edit mode)");
            lines.add("        [false] execVM \"dzn_gear\\dzn_gear_init.sqf\";");
        }

        if (d) {
            lines.add("");
            lines.add("        // dzn DynAI");
            lines.add("        [] execVM \"dzn_dynai\\dzn_dynai_init.sqf\";");
        }

        if (c) {
            lines.add("");
            lines.add("        // dzn CivEn");
            lines.add("        [] execVM \"dzn_civen\\dzn_civen_init.sqf\";");
        }

        if (t) {
            lines.add("");
            lines.add("        // TS Framework");
            lines.add("        [] execVM \"dzn_tSFramework\\dzn_tSFramework_Init.sqf\";");
            lines.add("        // dzn AAR");
            lines.add("        [] execVM \"dzn_brv\\dzn_brv_init.sqf\";");
        }

        lines.add("};");
        lines.add("// *****");

        Path file = Paths.get("Temp".concat("//init.sqf"));
        Files.write(file, lines, Charset.forName("UTF-8"));
    }

    public static void ProcessKits(boolean g, String path, String[] urls) throws IOException {
        if (!g) { return; }
        File outputFolder = new File(path.concat("\\dzn_gear"));
        File kitsSummary = new File(path.concat("\\dzn_gear\\Kits.sqf"));

        for (int i = 1; i < urls.length; i++) {
            try {
                String kit = urls[i];

                if (kit.length() != 0) {
                    LogToFile("    Processing kit - ".concat(kit));
                    String filename = DownloadRepository(kit, true);

                    File downloadedKitFile = new File (filename);
                    File kitFile = new File(
                            outputFolder.getAbsolutePath()
                                    .concat("\\")
                                    .concat( URLDecoder.decode(filename, "UTF-8"))
                    );
                    downloadedKitFile.renameTo(kitFile);
                    LogToFile("    Kit file downloaded: ".concat(kitFile.getName()));

                    Files.write(
                            kitsSummary.toPath()
                            , "\n#include \"".concat(kitFile.getName()).concat("\"").getBytes()
                            , StandardOpenOption.APPEND
                    );
                    LogToFile("    Kit file applied to Kits.sqf. Done!");
                }
            } catch (NullPointerException e) {}
        }
    }

    public static void UpdateFiles(File tempFolder, File outputFolder, boolean backupNeeded) throws IOException {
        if ( backupNeeded ) {
            BackupFolder(tempFolder, outputFolder);
        }

        CopyFolder(tempFolder, outputFolder);
        DeleteFolder(tempFolder);
    }

    public static void BackupFolder(File src, File dest) throws IOException {
        if (src.isDirectory()) {
            String files[] = src.list();
            //System.out.println("Backup (folder): Checking  " + src);
            for (String file : files) {
                BackupFolder(new File(src, file), new File(dest, file));
            }
        } else {
            if ( !(Arrays.asList(FilesToBackup).contains(src.getName())) ) {
                //System.out.println("Backup (file): File is not for backup - " + src.getName());
                return;
            }

            if ( !(dest.exists()) ) {
                //System.out.println("Backup (file): Destination file not exists - " + dest);
                return;
            }

            if (Arrays.equals(Files.readAllBytes(src.toPath()), Files.readAllBytes(dest.toPath()))) {
                //System.out.println("Backup (file): Files are the same, No backup (" + dest + ")");
                return;
            }

            //System.out.println("Backup (file): Backuping diff files! (" + dest + ")");

            String name = dest.toString();
            String[] exts = name.split(Pattern.quote("."));
            exts[exts.length - 1] = "old.".concat(exts[exts.length - 1]);
            dest.renameTo(new File( String.join(".", exts) ));
        }
    }

    public static Properties getSettings() throws IOException {
        Properties prop = new Properties();

        try {
            FileInputStream fstream = new FileInputStream("Settings.txt");
            prop.load(fstream);
            prop.setProperty("INSTALLATION_FOLDER", prop.getProperty("INSTALLATION_FOLDER").replace("@","\\"));
        } catch (IOException e) {
            System.out.println("No file, using defaults!");
            prop.setProperty("INSTALLATION_FOLDER", "");
            prop.setProperty("MAKE_BACKUP", "true");
            prop.setProperty("INSTALL_DZN_CommonFunctions", "true");
            prop.setProperty("INSTALL_DZN_GEAR", "true");
            prop.setProperty("INSTALL_DZN_DYNAI", "true");
            prop.setProperty("INSTALL_DZN_CIVEN", "false");
            prop.setProperty("INSTALL_DZN_TSF", "true");
            prop.setProperty("KIT_1", "");
            prop.setProperty("KIT_2", "");
            prop.setProperty("KIT_3", "");
            prop.setProperty("REPO_DZN_CommonFunctions", "");
            prop.setProperty("REPO_DZN_GEAR", "");
            prop.setProperty("REPO_DZN_DYNAI", "");
            prop.setProperty("REPO_DZN_CIVEN", "");
            prop.setProperty("REPO_DZN_TSF", "");
        }

        return prop;
    }

    public static void saveSettings(String path
            , boolean doBackup
            , boolean doCF, boolean doG, boolean doDA, boolean doCN, boolean doTSF
            , String urlCF, String urlG, String urlDA, String urlCN, String urlTSF
            , String[] kits
    ) throws IOException {
        List<String> lines = new ArrayList<String>();
        lines.add("INSTALLATION_FOLDER=".concat(path.replace("\\","@")));
        lines.add("MAKE_BACKUP=".concat(Boolean.toString(doBackup)));

        lines.add("INSTALL_DZN_DYNAI=".concat(Boolean.toString(doDA)));
        lines.add("INSTALL_DZN_GEAR=".concat(Boolean.toString(doG)));
        lines.add("INSTALL_DZN_CommonFunctions=".concat(Boolean.toString(doCF)));
        lines.add("INSTALL_DZN_CIVEN=".concat(Boolean.toString(doCN)));
        lines.add("INSTALL_DZN_TSF=".concat(Boolean.toString(doTSF)));

        lines.add("REPO_DZN_DYNAI=".concat(urlDA));
        lines.add("REPO_DZN_CommonFunctions=".concat(urlCF));
        lines.add("REPO_DZN_GEAR=".concat(urlG));
        lines.add("REPO_DZN_CIVEN=".concat(urlCN));
        lines.add("REPO_DZN_TSF=".concat(urlTSF));

        lines.add("KIT_1=".concat(kits[0]));
        lines.add("KIT_2=".concat(kits[1]));
        lines.add("KIT_3=".concat(kits[2]));

        File file = new File("Settings.txt");

        if (file.exists()) {
            PrintWriter writer = new PrintWriter(file);
            writer.print("");
            writer.close();
        };

        Files.write( (new File("Settings.txt")).toPath(), lines, Charset.forName("UTF-8"));
    }
}
