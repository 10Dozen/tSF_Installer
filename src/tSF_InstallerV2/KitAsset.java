package tSF_InstallerV2;

import javafx.scene.control.Label;
import javafx.scene.control.TextField;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by Zodius on 04.06.2017.
 */
public class KitAsset {
	public int count;
	public String[] urls;
	public Label[] labels;
	public TextField[] fields;


	public KitAsset(String[] urls) {
		this.count = urls.length;
		this.urls = urls;
		this.labels = new Label[this.count];
		this.fields = new TextField[this.count];

		for (int i = 0; i < this.count; i++) {
			TextField tf = new TextField(urls[i]);
			tf.setPromptText("URL to zip-archive");
			fields[i] = tf;
			labels[i] = new Label("Kit_".concat(Integer.toString(i)));
		}
	}

	public List<String> getProperties() {
		List<String> lines = new ArrayList<>();
		for (int i = 0; i < this.count; i++) {
			lines.add( "Kit_".concat(Integer.toString(i + 1)).concat("=").concat( this.urls[i] ) );
		}

		return lines;
	}

	public void Update() {
		for (int i = 0; i < this.count; i++) {
			urls[i] = fields[i].getText();
		}
	}

}
