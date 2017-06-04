package tSFInstaller_v2;

import javafx.application.Application;
import javafx.event.ActionEvent;
import javafx.event.EventHandler;
import javafx.geometry.Insets;
import javafx.geometry.Pos;
import javafx.scene.Scene;
import javafx.scene.control.*;
import javafx.scene.image.Image;
import javafx.scene.image.ImageView;
import javafx.scene.layout.ColumnConstraints;
import javafx.scene.layout.GridPane;
import javafx.scene.layout.HBox;
import javafx.scene.paint.Color;
import javafx.stage.DirectoryChooser;
import javafx.stage.Stage;

import java.io.File;
import java.io.IOException;
import java.util.*;

public class Main extends Application {

    @Override
    public void start(Stage primaryStage) throws IOException {
        Properties prop = Installer.getSettings();

        primaryStage.setTitle("tSF Installer (v2.0)");
        primaryStage.getIcons().add(new Image("icon.jpg"));

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

        ImageView title = new ImageView();
        Image img = new Image ("logo.png");
        title.setFitWidth(590);
        title.setImage(img);
        grid.add(title, 0,0,2,1);
        /*
        Text title = new Text("tS Framework Installer");
        title.setFont(Font.font("Tahoma", FontWeight.BOLD, 20));
        title.setTextAlignment(TextAlignment.JUSTIFY);
        grid.add(title,0,0,2,1);
*/
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
        btn.setMinWidth(200);
        HBox hbBtn = new HBox(10);
        hbBtn.setAlignment(Pos.BOTTOM_RIGHT);
        hbBtn.getChildren().add(btn);
        hbBtn.setMinWidth(200);
        grid.add(hbBtn, 1, 15);

        Label msg = new Label();
        grid.add(msg, 0, 15,2,1);

        btn.setOnAction(new EventHandler<ActionEvent>() {
            @Override
            public void handle(ActionEvent event) {

                if (pathField.getText() == "Please select...") {
                    msg.setText("Installation folder is not selected!");
                    msg.setTextFill(Color.web("#9b0000"));
                } else {
                    try {
                        msg.setText("Installation in progress!");
                        msg.setTextFill(Color.web("#6dad00"));

                        String[] kits = {
                                InstallURLs.get("Kit #1").getText()
                                , InstallURLs.get("Kit #2").getText()
                                , InstallURLs.get("Kit #3").getText()
                        };

                        boolean result = Installer.Install(
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
                                , kits
                        );

                        if (result) {
                            msg.setText("All done!");
                            msg.setTextFill(Color.web("#6dad00"));
                        }
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                }
            }
        });

        Scene scene = new Scene(grid, 640, 530);
        primaryStage.setScene(scene);
        primaryStage.centerOnScreen();
        primaryStage.show();
    }


    public static void main(String[] args) {
        launch(args);
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