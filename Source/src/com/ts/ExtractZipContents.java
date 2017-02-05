package com.ts;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Enumeration;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

public class ExtractZipContents {

	public static String UnzipFile(String filename) {
		File folder = new File ("");
		boolean isBasicFolderExist = false;
		try {
			// Open the zip file
			ZipFile zipFile = new ZipFile(filename);
			Enumeration<?> enu = zipFile.entries();
			while (enu.hasMoreElements()) {
				System.out.print("|");

				ZipEntry zipEntry = (ZipEntry) enu.nextElement();

				String name = zipEntry.getName();
				long size = zipEntry.getSize();
				long compressedSize = zipEntry.getCompressedSize();
				/*
				System.out.printf("name: %-20s | size: %6d | compressed size: %6d\n",
						name, size, compressedSize);
				 */

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

}