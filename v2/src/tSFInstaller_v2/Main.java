package tSFInstaller_v2;

import javafx.application.Application;
import javafx.event.ActionEvent;
import javafx.event.EventHandler;
import javafx.geometry.Insets;
import javafx.geometry.Pos;
import javafx.scene.Scene;
import javafx.scene.control.*;
import javafx.scene.layout.ColumnConstraints;
import javafx.scene.layout.GridPane;
import javafx.scene.layout.HBox;
import javafx.scene.text.Font;
import javafx.scene.text.FontWeight;
import javafx.scene.text.Text;
import javafx.scene.text.TextAlignment;
import javafx.stage.DirectoryChooser;
import javafx.stage.Stage;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.util.*;

public class Main extends Application {

    @Override
    public void start(Stage primaryStage) throws IOException {
        Properties prop = getSettings();

        primaryStage.setTitle("tSF Installer (v2.0)");

        GridPane grid = new GridPane();
        grid.setAlignment(Pos.CENTER);
        grid.setHgap(20);
        grid.setVgap(10);
        grid.setPadding(new Insets(25,25,25,25));

        ColumnConstraints column1 = new ColumnConstraints();
        column1.setPercentWidth(30);
        ColumnConstraints column2 = new ColumnConstraints();
        column2.setPercentWidth(70);
        grid.getColumnConstraints().addAll(column1, column2);

        Text title = new Text("tS Framework Installer");
        title.setFont(Font.font("Tahoma", FontWeight.BOLD, 20));
        title.setTextAlignment(TextAlignment.JUSTIFY);
        grid.add(title,0,0,2,1);

        Label pathField = new Label();
        if (prop.getProperty("INSTALLATION_FOLDER") == "") {
            pathField.setText("Please select...");
        } else {
            pathField.setText(prop.getProperty("INSTALLATION_FOLDER"));
        }
        grid.add(pathField, 1,1);

        Button pathBtn = new Button("Installation Folder");
        pathBtn.setOnAction(new EventHandler<ActionEvent>() {
            @Override
            public void handle(ActionEvent event) {
                DirectoryChooser directoryChooser = new DirectoryChooser();
                if ( prop.getProperty("INSTALLATION_FOLDER") != "") {
                    directoryChooser.setInitialDirectory(new File(prop.getProperty("INSTALLATION_FOLDER")));
                }
                File selectedDirectory =
                        directoryChooser.showDialog(primaryStage);

                if (selectedDirectory != null) {
                    pathField.setText(selectedDirectory.getAbsolutePath());
                }
            }
        });
        grid.add(pathBtn, 0,1);
        grid.add(new Separator(), 0,2,2,1);

        grid.add(new Label("Make backup?"),0,3);
        CheckBox backupChBox = new CheckBox();
        backupChBox.setSelected( Boolean.parseBoolean( prop.getProperty("MAKE_BACKUP") ) );
        grid.add(backupChBox, 1,3);

        grid.add(new Label("Select components to install (set URL to use non-master branch)"),0,4,2,1);

        SetInstallationLines(grid, prop,"dzn_Common Functions","INSTALL_DZN_CommonFunctions", "REPO_DZN_CommonFunctions",5);
        SetInstallationLines(grid, prop,"dzn_Gear","INSTALL_DZN_GEAR", "REPO_DZN_GEAR",6);
        SetInstallationLines(grid, prop,"dzn_DynAI","INSTALL_DZN_DYNAI", "REPO_DZN_DYNAI",7);
        SetInstallationLines(grid, prop,"dzn_CivEn","INSTALL_DZN_CIVEN", "REPO_DZN_CIVEN",8);
        SetInstallationLines(grid, prop,"dzn_tSFramework","INSTALL_DZN_TSF", "REPO_DZN_TSF",9);

        grid.add(new Separator(), 0,10,2,1);
        grid.add(new Label("Set KitLink to install kits from collection"),0,11,2,1);

        SetInstallationLines(grid, prop, "Kit #1", "", "KIT_1", 12 );
        SetInstallationLines(grid, prop, "Kit #2", "", "KIT_2", 13 );
        SetInstallationLines(grid, prop, "Kit #3", "", "KIT_3", 14 );

        Button btn = new Button("Install");
        HBox hbBtn = new HBox(10);
        hbBtn.setAlignment(Pos.BOTTOM_RIGHT);
        hbBtn.getChildren().add(btn);
        grid.add(hbBtn, 1, 15);

        btn.setOnAction(new EventHandler<ActionEvent>() {
            @Override
            public void handle(ActionEvent event) {
               /* String path
            , boolean doBackup
            , boolean doCF, boolean doG, boolean doDA, boolean doCN, boolean doTSF
            , String urlCF, String urlG, String urlDA, String urlCN, boolean urlTSF
            , String kit1, String kit2, String kit3*/

                try {


                    Installer.Install(
                        pathField.getText()
                        , backupChBox.isSelected()

                        , InstallNeeded.get("dzn_Common Functions").isSelected()
                            , InstallNeeded.get("dzn_Gear").isSelected()
                            , InstallNeeded.get("dzn_DynAI").isSelected()
                            , InstallNeeded.get("dzn_CivEn").isSelected()
                            , InstallNeeded.get("dzn_tSFramework").isSelected()

                            , InstallURLs.get("dzn_Common Functions").getText()
                            , InstallURLs.get("dzn_Gear").getText()
                            , InstallURLs.get("dzn_DynAI").getText()
                            , InstallURLs.get("dzn_CivEn").getText()
                            , InstallURLs.get("dzn_tSFramework").getText()
                        , "", "", ""
                    );
                } catch (IOException e) {
                    e.printStackTrace();
                }

                //  /*  pathField.getText()*/
            }
        });

        Scene scene = new Scene(grid, 640, 520);
        primaryStage.setScene(scene);
        primaryStage.centerOnScreen();
        primaryStage.show();
    }


    public static void main(String[] args) {
        launch(args);
    }

    public Properties getSettings() throws IOException {
        Properties prop = new Properties();

        try {
            FileInputStream fstream = new FileInputStream("Settings.txt");
            prop.load(fstream);
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

    static Map<String, CheckBox> InstallNeeded;
    static Map<String, TextField> InstallURLs;
    static {
        InstallNeeded = new HashMap<String, CheckBox>();
        InstallURLs = new HashMap<String, TextField>();
    }

    private void SetInstallationLines(
            GridPane grid, Properties prop, String title, String chbxPropName, String urlPropName, int line
    ) {
        if (chbxPropName != "") {
            CheckBox chbx = new CheckBox(title);
            chbx.setSelected(Boolean.parseBoolean(prop.getProperty(chbxPropName)));
            grid.add(chbx, 0, line);
            InstallNeeded.put(title, chbx);
        } else {
            grid.add(new Label(title),0,line);
        }

        TextField tf = new TextField( prop.getProperty(urlPropName));
        tf.setPromptText("URL to zip-archive");
        grid.add(tf,1,line);
        InstallURLs.put(title, tf);
    }


}
