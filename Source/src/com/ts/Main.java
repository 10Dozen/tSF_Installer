 package com.ts;

import java.io.*;
import java.net.URL;
import java.nio.channels.Channels;
import java.nio.channels.ReadableByteChannel;

import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;

public class Main {
    private static final Map<String, String> Repositories;
    static
    {
        Repositories = new HashMap<String, String>();
        Repositories.put("REPO_DZN_DYNAI", "https://github.com/10Dozen/dzn_dynai/archive/master.zip");
        Repositories.put("REPO_DZN_GEAR", "https://github.com/10Dozen/dzn_gear/archive/master.zip");
        Repositories.put("REPO_DZN_CIVEN", "https://github.com/10Dozen/dzn_civen/archive/master.zip");
        Repositories.put("REPO_DZN_TSF", "https://github.com/10Dozen/dzn_tSFramework/archive/master.zip");
    }

    public static void main(String[] args) throws IOException {

        // Read Settings.txt
        FileInputStream fstream = new FileInputStream("Settings.txt");
        Properties prop = new Properties();
        prop.load(fstream);

        boolean needGear = GetBoolProperty(prop, "INSTALL_DZN_GEAR");
        boolean needDynai = GetBoolProperty(prop, "INSTALL_DZN_DYNAI");
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

        System.out.println(" Install: "
                .concat("[ Gear - ".concat(Boolean.toString(needGear))).concat(" ]")
                .concat("[ DynAI - ".concat(Boolean.toString(needDynai))).concat(" ]")
                .concat("[ CivEn - ".concat(Boolean.toString(needCiven))).concat(" ]")
                .concat("[ tSF - ".concat(Boolean.toString(needTSF))).concat(" ]")
        );
        System.out.println(" --------------------------- ");

        ProcessRepository("dzn_DynAI", needDynai, GetStringProperty(prop, "REPO_DZN_DYNAI"), outputFolder);
        ProcessRepository("dzn_Gear", needGear, GetStringProperty(prop, "REPO_DZN_GEAR"), outputFolder);
        ProcessRepository("dzn_CivEn", needCiven, GetStringProperty(prop, "REPO_DZN_CIVEN"), outputFolder);
        ProcessRepository("dzn_tSF", needTSF, GetStringProperty(prop, "REPO_DZN_TSF"), outputFolder);

        System.out.println("Compiling init.sqf");
        GenerateInitSQF(outputFolderPath, needGear, needDynai, needCiven, needTSF);

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

    public static void ProcessRepository(String name, boolean isNeeded, String url, File outputFolder) throws IOException {
        if (!isNeeded) { return; }
        System.out.println("Installing ".concat(name));

        System.out.print("    Downloading: ");
        File folder = new File( DownloadRepository(url) );
        System.out.println(" >> Done");

        System.out.print("    Installing: ");
        CopyFolder(folder, outputFolder);
        DeleteFolder(folder);
        DeleteFolder(new File (GetRepoFileName(url)));
        System.out.println(" >> Done");
    }

    public static String DownloadRepository(String src) throws IOException {
        URL website = new URL(src);

        String filename = GetRepoFileName(src);

        ReadableByteChannel rbc = Channels.newChannel(website.openStream());
        FileOutputStream fos = new FileOutputStream(filename);
        fos.getChannel().transferFrom(rbc, 0, Long.MAX_VALUE);

        String folder = ExtractZipContents.UnzipFile(filename);

        rbc.close();
        fos.close();

        return folder;
    }

    public static void CopyFolder(File src, File dest) throws IOException {
        if (src.isDirectory()) {
            dest.mkdir();
            // System.out.println("Directory copied from " + src + " to " + dest);
            System.out.print("|");

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
            System.out.print("|");
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
        System.out.print("|");

        return dir.delete();
    }

    public static void GenerateInitSQF(
            String folder, boolean addGear, boolean addDynai, boolean addCiven, boolean addTSF
    ) throws IOException {
        List<String> lines = new ArrayList<String>();

        lines.add("// *****");

        if (addGear) {
            lines.add("");
            lines.add("// dzn Gear");
            lines.add("  // set true to engage Edit mode");
            lines.add("[false] execVM \"dzn_gear\\dzn_gear_init.sqf\";");
        }

        if (addDynai) {
            lines.add("");
            lines.add("// dzn DynAI");
            lines.add("[] execVM \"dzn_dynai\\dzn_dynai_init.sqf\";");
        }

        if (addCiven) {
            lines.add("");
            lines.add("// dzn CivEn");
            lines.add("[] execVM \"dzn_civen\\dzn_civen_init.sqf\";");
        }

        if (addTSF) {
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

}
