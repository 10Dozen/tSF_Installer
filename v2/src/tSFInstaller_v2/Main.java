package tSFInstaller_v2;

import javafx.application.Application;
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
import javafx.stage.Stage;

public class Main extends Application {

    @Override
    public void start(Stage primaryStage) {
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

        Label path = new Label("Installation folder");
        grid.add(path, 0,1);

        TextField pathField = new TextField();
        grid.add(pathField, 1,1);
        grid.add(new Separator(), 0,2,2,1);

        grid.add(new Label("Make backup?"),0,3);
        CheckBox backupChBox = new CheckBox();
        grid.add(backupChBox, 1,3);

        grid.add(new Label("Select components to install (set URL to use non-master branch)"),0,4,2,1);

        CheckBox comFunChBox = new CheckBox("dzn_Common Functions");
        grid.add(comFunChBox,0,5);
        grid.add(new TextField(), 1,5);

        CheckBox gearChBox = new CheckBox("dzn_Gear");
        grid.add(gearChBox,0,6);
        grid.add(new TextField(), 1,6);

        CheckBox dynaiChBox = new CheckBox("dzn_DynAI");
        grid.add(dynaiChBox,0,7);
        grid.add(new TextField(), 1,7);

        CheckBox civenChBox = new CheckBox("dzn_CivEn");
        grid.add(civenChBox,0,8);
        grid.add(new TextField(), 1,8);

        CheckBox tsfChBox = new CheckBox("tS Framework");
        grid.add(tsfChBox,0,9);
        grid.add(new TextField(), 1,9);

        grid.add(new Separator(), 0,10,2,1);
        grid.add(new Label("Set KitLink to install kits from collection"),0,11,2,1);
        grid.add(new Label("Kit #1"), 0,12);
        grid.add(new TextField(), 1,12);
        grid.add(new Label("Kit #2"), 0,13);
        grid.add(new TextField(), 1,13);
        grid.add(new Label("Kit #3"), 0,14);
        grid.add(new TextField(), 1,14);

        Button btn = new Button("Install");
        HBox hbBtn = new HBox(10);
        hbBtn.setAlignment(Pos.BOTTOM_RIGHT);
        hbBtn.getChildren().add(btn);
        grid.add(hbBtn, 1, 15);

        Scene scene = new Scene(grid, 640, 520);
        primaryStage.setScene(scene);
        primaryStage.centerOnScreen();
        primaryStage.show();
    }


    public static void main(String[] args) {
        launch(args);
    }
}
