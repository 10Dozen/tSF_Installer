package tSF_InstallerV2;

import javafx.scene.control.CheckBox;
import javafx.scene.control.Label;

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
	static {
		Repositories = new HashMap<String, String>();
		Repositories.put("dzn_DynAI", "https://github.com/10Dozen/dzn_dynai/archive/master.zip");
		Repositories.put("dzn_Gear", "https://github.com/10Dozen/dzn_gear/archive/master.zip");
		Repositories.put("dzn_CommonFunctions", "https://github.com/10Dozen/dzn_commonFunctions/archive/master.zip");
		Repositories.put("dzn_CivEn", "https://github.com/10Dozen/dzn_civen/archive/master.zip");
		Repositories.put("dzn_tSFramework", "https://github.com/10Dozen/dzn_tSFramework/archive/master.zip");
	}
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

	public String path;
	public Label pathLabel;
	public boolean backup;
	public CheckBox backupChBx;

	public InstallAsset commonFunctions;
	public InstallAsset dynai;
	public InstallAsset gear;
	public InstallAsset civen;
	public InstallAsset tsf;
	public KitAsset kits;

	public Installer() {
		Properties prop = new Properties();

		try {
			FileInputStream fstream = new FileInputStream("Settings.cfg");
			prop.load(fstream);
			prop.setProperty("INSTALLATION_FOLDER", prop.getProperty("INSTALLATION_FOLDER", "").replace("@", "\\"));
		} catch (IOException e) {
			System.out.println("No file, using defaults!");
		}

		this.path = prop.getProperty("INSTALLATION_FOLDER", "");
		this.pathLabel = new Label();
		this.SetPathLabel( this.path );

		this.backup = Boolean.parseBoolean( prop.getProperty("MAKE_BACKUP", "true") );
		backupChBx = new CheckBox();
		backupChBx.setSelected(backup);

		this.commonFunctions = new InstallAsset(prop,"dzn_CommonFunctions");
		this.gear = new InstallAsset(prop,"dzn_Gear");
		this.dynai = new InstallAsset(prop,"dzn_DynAI"	);
		this.civen = new InstallAsset(prop,"dzn_CivEn");
		this.tsf = new InstallAsset(prop, "dzn_tSFramework");

		this.kits = new KitAsset(new String[] {
				prop.getProperty("Kit_1", "")
				, prop.getProperty("Kit_2", "")
				, prop.getProperty("Kit_3", "")
		});
	}

	public void SetPathLabel(String url) {
		String displayUrl = url.isEmpty() ? "Please select..." : url.replace("@", "\\");
		if ( url.length() > 68) {
			displayUrl = (url.substring(0,15)).concat(" ... ").concat(url.substring( url.length() - 50));
		}

		this.pathLabel.setText( displayUrl );
		this.path = url;
	}

	public boolean Install() throws IOException {
		PrintWriter writer = new PrintWriter(log);
		writer.print("");
		writer.close();

		UpdateInstallerSettings();

		LogToFile(" ------ tSF Installer ------ ");
		LogToFile(" Output folder: ".concat(path));

		if (!(new File(path)).exists()) {
			LogToFile("Aborted! No INSTALLATION_FOLDER exists!");
			System.exit(0);
		}

		LogToFile("Install:");
		LogInstallationList();

		LogToFile(" --------------------------- ");
		ProcessRepository(this.dynai);
		ProcessRepository(this.gear);
		ProcessRepository(this.commonFunctions);
		ProcessRepository(this.civen);
		ProcessRepository(this.tsf);

		LogToFile("Compiling init.sqf");
		GenerateInitSQF();

		LogToFile("\n --------------------------- \n dzn_gear Kits:");
		ProcessKits();
		LogToFile(" --------------------------- \n  Installation... ");
		UpdateFiles();
		LogToFile(" Done! \n  --------------------------- \n  All done! Have a nice day!");

		SaveSettings();

		return  true;
	}

	private void LogToFile(String line) throws IOException {
		System.out.println(line);
		Files.write(log.toPath(), "\n".concat(line).getBytes(), StandardOpenOption.APPEND );
	}

	private void LogInstallationList() throws IOException {
		if (this.commonFunctions.isNeeded()) { LogToFile("       - CommonFunctions"); }
		if (this.gear.isNeeded()) { LogToFile("       - Gear"); }
		if (this.dynai.isNeeded()) { LogToFile("       - DynAI"); }
		if (this.civen.isNeeded()) { LogToFile("       - CivEn"); }
		if (this.tsf.isNeeded()) { LogToFile("       - tS Framework"); }
	}

	private void ProcessRepository(InstallAsset asset) throws IOException {
		if (!asset.isNeeded()) { return; }
		LogToFile( "Installing ".concat(asset.getTitle()) );
		String url = asset.getUrl().isEmpty() ? Repositories.get(asset.getTitle()) : asset.getUrl();
		LogToFile("    Downloading... ".concat(url));
		File folder = new File( DownloadRepository(url,false) );
		CopyFolder(folder, new File("Temp"));
		DeleteFolder(folder);
		DeleteFolder(new File (GetRepoFileName(url)));
		LogToFile("     Done!");
	}

	private String DownloadRepository(String src, boolean isSingleFile) throws IOException {
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

	private String GetRepoFileName(String src) {
		String[] segments = src.split("/");
		return (segments[segments.length-1]);
	}

	private String UnzipFile(String filename) {
		File folder = new File ("");
		boolean isBasicFolderExist = false;
		try {
			// Open the zip file
			ZipFile zipFile = new ZipFile(filename);
			Enumeration<?> enu = zipFile.entries();
			while (enu.hasMoreElements()) {
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

	private void CopyFolder(File src, File dest) throws IOException {
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

	private boolean DeleteFolder(File dir) throws IOException {
		if (dir.isDirectory()) {
			File[] children = dir.listFiles();
			for (int i = 0; i < children.length; i++) {
				boolean success = DeleteFolder(children[i]);
				if (!success) { return false; }
			}
		}
		return dir.delete();
	}

	private void GenerateInitSQF() throws IOException {
		List<String> lines = new ArrayList<>();
		lines.add("//	Tacitcal Shift Framework initialization");
		lines.add("[] spawn {");
		lines.add("        waitUntil { !isNil \"MissionDate\" };");

		if (this.gear.isNeeded()) {
			lines.add("");
			lines.add("        // dzn Gear 	(set true to engage Edit mode)");
			lines.add("        [false] execVM \"dzn_gear\\dzn_gear_init.sqf\";");
		}

		if (this.dynai.isNeeded()) {
			lines.add("");
			lines.add("        // dzn DynAI");
			lines.add("        [] execVM \"dzn_dynai\\dzn_dynai_init.sqf\";");
		}

		if (this.civen.isNeeded()) {
			lines.add("");
			lines.add("        // dzn CivEn");
			lines.add("        [] execVM \"dzn_civen\\dzn_civen_init.sqf\";");
		}

		if (this.tsf.isNeeded()) {
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

	private void ProcessKits() throws IOException {
		if (!this.gear.isNeeded()) { return; }

		File outputFolder = new File("Temp".concat("\\dzn_gear"));
		File kitsSummary = new File("Temp".concat("\\dzn_gear\\Kits.sqf"));

		for (int i = 0; i < this.kits.count; i++) {

			if (!(this.kits.urls[i]).isEmpty()) {
				LogToFile("    Processing kit - ".concat(this.kits.urls[i]));
				String filename = DownloadRepository(this.kits.urls[i], true);

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
		}
	}

	private void UpdateFiles() throws IOException {
		File temp = new File("Temp");
		File output = new File(this.path);
		if (this.backup) {
			BackupFolder(temp, output);
		}

		CopyFolder(temp, output);
		DeleteFolder(temp);
	}

	private void BackupFolder(File src, File dest) throws IOException {
		if (src.isDirectory()) {
			String files[] = src.list();
			for (String file : files) {
				BackupFolder(new File(src, file), new File(dest, file));
			}
		} else {
			if ( !(Arrays.asList(FilesToBackup).contains(src.getName())) ) { return; }
			if ( !(dest.exists()) ) { return; }
			if (Arrays.equals(Files.readAllBytes(src.toPath()), Files.readAllBytes(dest.toPath()))) { return;  }

			String name = dest.toString();
			String[] exts = name.split(Pattern.quote("."));
			exts[exts.length - 1] = "old.".concat(exts[exts.length - 1]);
			dest.renameTo(new File( String.join(".", exts) ));
		}
	}

	private void UpdateInstallerSettings() {
		this.backup = this.backupChBx.isSelected();

		this.commonFunctions.Update();
		this.gear.Update();
		this.dynai.Update();
		this.civen.Update();
		this.tsf.Update();

		this.kits.Update();
	}

	private void SaveSettings() throws IOException {
		List<String> lines = new ArrayList<String>();

		lines.add("INSTALLATION_FOLDER=".concat( this.path.replace("\\","@") ));
		lines.add("MAKE_BACKUP=".concat(Boolean.toString( this.backup )));
		lines.addAll( commonFunctions.getProperties() );
		lines.addAll( gear.getProperties() );
		lines.addAll( dynai.getProperties() );
		lines.addAll( civen.getProperties() );
		lines.addAll( tsf.getProperties() );
		lines.addAll( kits.getProperties() );

		File file = new File("Settings.cfg");
		if (file.exists()) {
			PrintWriter writer = new PrintWriter(file);
			writer.print("");
			writer.close();
		}

		Files.write( (new File("Settings.cfg")).toPath(), lines, Charset.forName("UTF-8"));
	}

}
