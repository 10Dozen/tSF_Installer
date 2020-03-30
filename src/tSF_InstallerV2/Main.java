package tSF_InstallerV2;

import javafx.application.Application;
import javafx.application.Platform;
import javafx.concurrent.Task;
import javafx.event.ActionEvent;
import javafx.event.EventHandler;
import javafx.geometry.Insets;
import javafx.geometry.Pos;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.control.CheckBox;
import javafx.scene.control.Label;
import javafx.scene.control.Separator;
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

public class Main extends Application {

    @Override
    public void start(Stage primaryStage) throws IOException {
        Installer installer = new Installer();

        primaryStage.setTitle("tSF Installer (v2.2.2)");
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

        grid.add( installer.pathLabel, 1,1 );
        Button pathBtn = new Button("Installation Folder");
        pathBtn.setOnAction(new EventHandler<ActionEvent>() {
            @Override
            public void handle(ActionEvent event) {
                DirectoryChooser directoryChooser = new DirectoryChooser();

                System.out.println(installer.path);
                if (installer.path.isEmpty()) { System.out.println("No installation folder"); };

                File installDir;
                if (!installer.path.isEmpty()) {
                    installDir = new File(installer.path);
                    while (!installDir.exists()) {
                        installDir = installDir.getParentFile();
                    }

                    if (!installer.path.isEmpty() && installDir.exists()) {
                        directoryChooser.setInitialDirectory(installDir);
                    }
                }

                File selectedDirectory = directoryChooser.showDialog(primaryStage);

                if (selectedDirectory != null) {
                    installer.SetPathLabel( selectedDirectory.getAbsolutePath() );
                }
            }
        });
        grid.add(pathBtn, 0,1);
        grid.add(new Separator(), 0,2,2,1);

        grid.add(new Label("Make backup?"),0,3);
        grid.add(installer.backupChBx, 1, 3);

        grid.add(new Label("Select components to install (set URL to use non-master branch)"),0,4,2,1);

        grid.add(installer.commonFunctions.getCb(), 0,5);
        grid.add(installer.commonFunctions.getField(), 1, 5);

        grid.add(installer.gear.getCb(), 0, 6);
        grid.add(installer.gear.getField(), 1, 6);

        grid.add(installer.dynai.getCb(), 0, 7);
        grid.add(installer.dynai.getField(), 1, 7);

        grid.add(installer.civen.getCb(), 0, 8);
        grid.add(installer.civen.getField(), 1, 8);

        grid.add(installer.tsf.getCb(), 0, 9);
        grid.add(installer.tsf.getField(), 1, 9);

        grid.add(new Separator(), 0,10,2,1);
        grid.add(new Label("Set KitLink to install kits from collection"),0,11,2,1);

        for (int i = 0; i < installer.kits.count; i++) {
            grid.add(installer.kits.labels[i], 0, 12 + i);
            grid.add(installer.kits.fields[i], 1, 12 + i);
        }

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
                if (installer.path.isEmpty()) {
                    msg.setText("Installation folder is not selected!");
                    msg.setTextFill(Color.web("#9b0000"));
                } else {

                    Task<Boolean> task = new Task<Boolean>() {
                        @Override protected Boolean call() throws IOException, InterruptedException {
                            Platform.runLater(new Runnable() {
                                @Override
                                public void run() {
                                    msg.setText("Installation in progress!");
                                    msg.setTextFill(Color.web("#6dad00"));
                                    pathBtn.setDisable(true);
                                    btn.setDisable(true);
                                }
                            });

                            boolean result = installer.Install();

                            Platform.runLater(new Runnable() {
                                @Override
                                public void run() {
                                    msg.setText("All done!");
                                    msg.setTextFill(Color.web("#6dad00"));
                                    pathBtn.setDisable(false);
                                    btn.setDisable(false);
                                }
                            });

                            return result;
                        }
                    };
                    Thread installThread = new Thread(task);
                    installThread.setDaemon(true);
                    installThread.start();
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
}
