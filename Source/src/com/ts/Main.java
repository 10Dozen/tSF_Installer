 package com.ts;

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

public class Main {
    private static final Map<String, String> Repositories;
    static
    {
        Repositories = new HashMap<String, String>();
        Repositories.put("REPO_DZN_DYNAI", "https://github.com/10Dozen/dzn_dynai/archive/master.zip");
        Repositories.put("REPO_DZN_GEAR", "https://github.com/10Dozen/dzn_gear/archive/master.zip");
        Repositories.put("REPO_DZN_CommonFunctions", "https://github.com/10Dozen/dzn_commonFunctions/archive/master.zip");
        Repositories.put("REPO_DZN_CIVEN", "https://github.com/10Dozen/dzn_civen/archive/master.zip");
        Repositories.put("REPO_DZN_TSF", "https://github.com/10Dozen/dzn_tSFramework/archive/master.zip");
    }

    public static void main(String[] args) throws IOException {

        // Read Settings.txt
        FileInputStream fstream = new FileInputStream("Settings.txt");
        Properties prop = new Properties();
        prop.load(fstream);

        boolean needDynai = GetBoolProperty(prop, "INSTALL_DZN_DYNAI");
        boolean needGear = GetBoolProperty(prop, "INSTALL_DZN_GEAR");
        boolean needCommonFunctions = GetBoolProperty(prop, "INSTALL_DZN_CommonFunctions");
        boolean needCiven = GetBoolProperty(prop, "INSTALL_DZN_CIVEN");
        boolean needTSF = GetBoolProperty(prop, "INSTALL_DZN_TSF");
        String outputFolderPath = prop.getProperty("INSTALLATION_FOLDER");

        System.out.println(" ------ tSF Installer ------ ");
        System.out.println(" Output folder: ".concat(outputFolderPath));
        File outputFolder = new File(outputFolderPath);
        if (!outputFolder.exists()) {
            System.out.println("Aborted! No INSTALLATION_FOLDER exists!");
            System.exit(0);
        }

        PrintInstallationList(needGear, needDynai, needCommonFunctions, needCiven, needTSF);

        ProcessRepository("dzn_DynAI", needDynai, GetStringProperty(prop, "REPO_DZN_DYNAI"), outputFolder);
        ProcessRepository("dzn_Gear", needGear, GetStringProperty(prop, "REPO_DZN_GEAR"), outputFolder);
        ProcessRepository("dzn_CommonFunctions", needCommonFunctions, GetStringProperty(prop, "REPO_DZN_CommonFunctions"), outputFolder);
        ProcessRepository("dzn_CivEn", needCiven, GetStringProperty(prop, "REPO_DZN_CIVEN"), outputFolder);
        ProcessRepository("dzn_tSF", needTSF, GetStringProperty(prop, "REPO_DZN_TSF"), outputFolder);

        System.out.println("Compiling init.sqf");
        GenerateInitSQF(outputFolderPath, needGear, needDynai, needCiven, needTSF);

        System.out.println("\n --------------------------- \n dzn_gear Kits:");
        ProcessKits(prop, needGear, outputFolderPath);

        System.out.println(" --------------------------- ");
        System.out.println(" All done! Have a nice day!");
    }

    public static boolean GetBoolProperty (Properties prop, String param) {
        return (Integer.parseInt( prop.getProperty(param) ) > 0);
    }

    public static String GetStringProperty (Properties prop, String param) {
        String val = prop.getProperty(param);
        if (val.length() == 0) {
            val = Repositories.get(param);
        }

        return val;
    }

    public static String GetRepoFileName(String src) {
        String[] segments = src.split("/");

        return (segments[segments.length-1]);
    }

    public static void PrintInstallationList(boolean d, boolean g, boolean cf, boolean c, boolean t) {
        System.out.println(" Install: ");
        if (d) { System.out.println("       - DynAI"); }
        if (g) { System.out.println("       - Gear"); }
        if (cf) { System.out.println("       - CommonFunctions"); }
        if (c) { System.out.println("       - CivEn"); }
        if (t) { System.out.println("       - tS Framework"); }
        System.out.println(" --------------------------- ");
    }

    public static void ProcessRepository(String name, boolean isNeeded, String url, File outputFolder) throws IOException {
        if (!isNeeded) { return; }
        System.out.println("Installing ".concat(name));

        System.out.print("    Downloading... ");
        File folder = new File( DownloadRepository(url,false) );
        System.out.println(" Done!");

        System.out.print("    Installing... ");
        CopyFolder(folder, outputFolder);
        DeleteFolder(folder);
        DeleteFolder(new File (GetRepoFileName(url)));
        System.out.println("  Done!");
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
            out = ExtractZipContents.UnzipFile(filename);
        }

        rbc.close();
        fos.close();

        return out;
    }

    public static void CopyFolder(File src, File dest) throws IOException {
        if (src.isDirectory()) {
            dest.mkdir();
            // System.out.println("Directory copied from " + src + " to " + dest);
            // System.out.print("|");

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
            // System.out.print("|");
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

        //System.out.println("removing file or directory : " + dir.getName());
        // System.out.print("|");

        return dir.delete();
    }

    public static void GenerateInitSQF(
            String folder, boolean g, boolean d, boolean c, boolean t
    ) throws IOException {
        List<String> lines = new ArrayList<String>();

        lines.add("// *****");
        lines.add("tS_DEBUG = true;");

        if (g) {
            lines.add("");
            lines.add("// dzn Gear");
            lines.add("  // set true to engage Edit mode");
            lines.add("[false] execVM \"dzn_gear\\dzn_gear_init.sqf\";");
        }

        if (d) {
            lines.add("");
            lines.add("// dzn DynAI");
            lines.add("[] execVM \"dzn_dynai\\dzn_dynai_init.sqf\";");
        }

        if (c) {
            lines.add("");
            lines.add("// dzn CivEn");
            lines.add("[] execVM \"dzn_civen\\dzn_civen_init.sqf\";");
        }

        if (t) {
            lines.add("");
            lines.add("  // TS Framework");
            lines.add("[] execVM \"dzn_tSFramework\\dzn_tSFramework_Init.sqf\";");
            lines.add("  // dzn AAR");
            lines.add("[] execVM \"dzn_brv\\dzn_brv_init.sqf\";");
        }

        lines.add("");
        lines.add("// *****");

        Path file = Paths.get(folder.concat("//init.sqf"));
        Files.write(file, lines, Charset.forName("UTF-8"));
    }

    public static void ProcessKits(Properties prop, boolean g, String outputFolderPath) throws IOException {
        if (!g) { return; }
        File outputFolder = new File(outputFolderPath.concat("\\dzn_gear"));
        File kitsSummary = new File(outputFolderPath.concat("\\dzn_gear\\Kits.sqf"));

        for (int i = 1; i < 10; i++) {
            try {
                String kit = prop.getProperty("KIT_" + Integer.toString(i));

                if (kit.length() != 0) {
                    System.out.println("    Processing kit - ".concat(kit));
                    String filename = DownloadRepository(kit, true);

                    File downloadedKitFile = new File (filename);
                    File kitFile = new File(
                            outputFolder.getAbsolutePath()
                                    .concat("\\")
                                    .concat( URLDecoder.decode(filename, "UTF-8"))
                    );
                    downloadedKitFile.renameTo(kitFile);
                    System.out.println("    Kit file downloaded: ".concat(kitFile.getName()));

                    Files.write(
                            kitsSummary.toPath()
                            , "\n#include \"Kits US SF Ghost Recon 1-4-4.sqf\"".getBytes()
                            , StandardOpenOption.APPEND
                    );
                    System.out.println("    Kit file applied to Kits.sqf. Done!");
                }
            } catch (NullPointerException e) {}
        }
    }

}
