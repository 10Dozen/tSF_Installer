package tSF_InstallerV2;

import javafx.scene.control.CheckBox;
import javafx.scene.control.TextField;

import java.util.ArrayList;
import java.util.List;
import java.util.Properties;


public class InstallAsset {
	private boolean isNeeded;
	private String url;
	private CheckBox cb;
	private TextField field;
	private String title;

	public InstallAsset(Properties prop, String title) {
		this.title = title;
		this.isNeeded = Boolean.parseBoolean( prop.getProperty("INSTALL_".concat(title),"false") );
		this.url = prop.getProperty("REPO_".concat(title), "");

		this.cb = new CheckBox(title.replace("_", " "));
		this.cb.setSelected(this.isNeeded);

		this.field = new TextField();
		this.field.setPromptText("URL to zip-archive");
		this.field.setText( this.url.isEmpty() ? "" : this.url );
	}

	public boolean isNeeded() {
		return isNeeded;
	}

	public String getUrl() {
		return url;
	}

	public CheckBox getCb() {
		return this.cb;
	}

	public TextField getField() {
		return field;
	}

	public String getTitle() {
		return title;
	}

	public List<String> getProperties() {
		List<String> lines = new ArrayList<String>();
		lines.add("INSTALL_".concat(this.title).concat("=").concat( Boolean.toString(this.isNeeded)));
		lines.add("REPO_".concat(this.title).concat("=").concat(this.url));

		return lines;
	}

	public void Update() {
		this.isNeeded = this.cb.isSelected();
		this.url = this.field.getText();
	}

}
